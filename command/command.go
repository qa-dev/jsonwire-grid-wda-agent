// Package command предназначен для выполнения системных команд по работе с симулятором и WDA.
package command

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/jsonwire-grid-wda-agent/device"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"gopkg.in/alexcesaro/statsd.v2"
)

// Время ожидания URL'a WDA
const wdaURLWaitingTime time.Duration = time.Second * 30

var (
	regexpFindServerURL = regexp.MustCompile("ServerURLHere->http://(.+)<-ServerURLHere")
	regexpFindFilePath  = regexp.MustCompile("Writing diagnostic log for test session to:\n([a-zA-Z0-9\\.\\-\\_\\/]+)\\/Session")
)

// KillWDA предназначен для остановки процесса WDA связанного с запущенным симмулятором
func KillWDA(deviceID string) error {
	args := []string{"ps",
		"aux",
		"|",
		"grep",
		"xcodebuild",
		"|",
		"awk",
		"'{print $2 \" \" $18}'"}
	argsStr := strings.Join(args, " ")

	out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf(`%s`, argsStr)).Output()
	if err != nil {
		log.WithError(err).Error("Error find WDA xcodebuild process")
		return err
	}

	output := strings.Split(string(out), "\n")
	var pid string
	for _, wdaProcess := range output {
		if strings.Contains(wdaProcess, deviceID) {
			split := strings.Split(wdaProcess, " ")
			pid = split[0]
			break
		}
	}
	if pid != "" {
		log.Printf("Killing WDA xcodebuild process with [pid:%s, device ID:%s]", pid, deviceID)
		err := exec.Command("kill", pid).Run()
		if err != nil {
			log.WithError(err).Errorf("Error kill WDA process for pid:%s", pid)
			return err
		}
	} else {
		log.Println("Not found no one running WDA xcodebuild process")
	}
	return nil
}

// KillAllSimulators предназначен для остановки запущенных симуляторов.
func KillAllSimulators() {
	exec.Command("killall", "-15", "Simulator").Run()
}

// BootSimulator предназначен для запуска симулятора через open.
func BootSimulator(deviceID string) error {
	command := exec.Command("xcrun", "simctl", "boot", deviceID)
	stderr, err := command.StderrPipe()
	if err != nil {
		log.WithError(err).Errorf("Error while opening device: %v", deviceID)
		return err
	}

	if err := command.Start(); err != nil {
		log.WithError(err).Errorf("Error while opening device: %v", deviceID)
		return err
	}

	slurp, _ := ioutil.ReadAll(stderr)
	if string(slurp) != "" {
		log.Println(string(slurp))
	}

	if err := command.Wait(); err != nil {
		log.WithError(err).Errorf("Error while opening device: %v", deviceID)
		return err
	}
	// Выключение открытия окна симулятора (запущен по умолчанию в фоновом режиме)
	//exec.Command("open", "/Applications/Xcode.app/Contents/Developer/Applications/Simulator.app", "--args", "-CurrentDeviceUDID", deviceID).Run()
	log.Println("Simulator booted:", deviceID)
	return nil
}

// UninstallApp предназначен для удаления приложения с устройства.
func UninstallApp(deviceID string, bundleID string) error {
	// Не возвращаем ошибку в методе, так как теперь при удалении несуществующего appa она возникает
	command := exec.Command("xcrun", "simctl", "uninstall", deviceID, bundleID)
	stderr, err := command.StderrPipe()
	if err != nil {
		log.WithError(err).Errorf("Error while uninstalling app on deviceID:[%s], bundleID:[%s]", deviceID, bundleID)
		return nil
	}
	if err := command.Start(); err != nil {
		log.WithError(err).Errorf("Error while uninstalling app on deviceID:[%s], bundleID:[%s]", deviceID, bundleID)
		return nil
	}

	slurp, _ := ioutil.ReadAll(stderr)
	if string(slurp) != "" {
		log.Println(string(slurp))
	}

	if err := command.Wait(); err != nil {
		log.WithError(err).Errorf("Error while uninstalling app on deviceID:[%s], bundleID:[%s]", deviceID, bundleID)
		return nil
	}
	return nil
}

// InstallApp предназначен для установки приложения на устройство.
func InstallApp(deviceID string, appPath string) error {
	command := exec.Command("xcrun", "simctl", "install", deviceID, appPath)
	stderr, err := command.StderrPipe()
	if err != nil {
		log.WithError(err).Errorf("Error while installing app on deviceID:[%s], appPath:[%s]", deviceID, appPath)
		return err
	}

	if err := command.Start(); err != nil {
		log.WithError(err).Errorf("Error while installing app on deviceID:[%s], appPath:[%s]", deviceID, appPath)
		return err
	}

	slurp, _ := ioutil.ReadAll(stderr)
	if string(slurp) != "" {
		log.Println(string(slurp))
	}

	if err := command.Wait(); err != nil {
		log.WithError(err).Errorf("Error while installing app on deviceID:[%s], appPath:[%s]", deviceID, appPath)
		return err
	}
	return nil
}

type StdRead struct {
	path string
}

func (s *StdRead) Write(p []byte) (n int, err error) {
	go func() {
		path := getPathFromLog(p)
		if path != "" {
			s.path = path
		}
	}()
	return len(p), err
}

func startFindUrlProcess(urlInterception chan string, cmdReader *StdRead, defaultLogFilePath string) (*os.File, string, error) {
	var logFilePath string

	time.Sleep(time.Second * 6)
	if cmdReader.path != "" {
		logFilePath = filepath.Join(cmdReader.path, "StandardOutputAndStandardError.txt")
	} else {
		logFilePath = defaultLogFilePath
	}

	f, err := os.Open(logFilePath)
	if err != nil {
		log.WithError(err).Error("Error open device log file: ", logFilePath)
		return nil, "", err
	}

	go func() {
		url := getURLFromSimulatorLog(f)
		if url != "" {
			urlInterception <- url
		}
	}()

	return f, logFilePath, nil
}

func startWDAProcess(deviceID string, wdaPath string) (*exec.Cmd, *StdRead) {
	// Назначение, на каком устройстве стартуем WDA
	dest := fmt.Sprintf(`"platform=iOS Simulator,id=%s"`, deviceID)

	args := []string{"xcodebuild", "-project", wdaPath, "-scheme", "WebDriverAgentRunner", "-destination", dest, "test"}
	argsStr := strings.Join(args, " ")

	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf(`%s`, argsStr))
	cmdReader := &StdRead{}
	cmd.Stderr = cmdReader

	go cmd.Run()

	return cmd, cmdReader
}

// StartWDA предназначен для запуска WDA.
func StartWDA(deviceID string, wdaPath string, stats *statsd.Client) (string, error) {
	// Путь файла с логами WDA
	var logFilePath string

	// Команда запуска WDA
	var cmd *exec.Cmd

	// Канал для итераций поиска WDA URL
	var urlInterception = make(chan string)

	// Получаем стандартный путь к файлу лога и очищаем его
	defaultLogFilePath, err := getDefaultLogFilePath(deviceID)
	if err != nil {
		return "", err
	}

	// Запускаем процесс WDA
	cmd, cmdReader := startWDAProcess(deviceID, wdaPath)

	// Запускаем процесс поиска WDA URL
	file, logFilePath, err := startFindUrlProcess(urlInterception, cmdReader, defaultLogFilePath)
	if err != nil {
		log.WithError(err).Error("Error starting find WDA URL process (first iteration)")
		return "", err
	}

	select {
	case result := <-urlInterception:
		// Нашли WDA URL на первой итерации
		stats.Increment("wda_url.first_iteration_find")
		log.Println("Info: [WDA URL Getting at first iteration]")
		err := file.Close()
		if err != nil {
			log.WithError(err).Error("Error close file: " + logFilePath)
		}
		return result, nil
	case <-time.After(wdaURLWaitingTime):
		// Если после 30 секунд не нашли URL, значит вероятно подвис WDA. Пробуем перезапустить процесс.
		log.Printf("Do not wait URL from WDA log, rerun process [PID]:%d", cmd.Process.Pid)

		err := file.Close()
		if err != nil {
			log.WithError(err).Error("Error close file: " + logFilePath)
		}

		err = KillWDA(deviceID)
		if err != nil {
			log.WithError(err).Error("Error kill WDA after first iteration")
			return "", err
		}

		// Запускаем новый процесс WDA
		log.Println("Run new WDA process")
		cmd, cmdReader = startWDAProcess(deviceID, wdaPath)

		urlInterception = make(chan string)
		// Запускаем новый процесс поиска WDA URL
		file, logFilePath, err = startFindUrlProcess(urlInterception, cmdReader, defaultLogFilePath)
		if err != nil {
			log.WithError(err).Error("Error starting find WDA URL process (first iteration)")
			return "", err
		}
	}

	select {
	case result := <-urlInterception:
		// Нашли WDA URL на второй итерации
		stats.Increment("wda_url.second_iteration_find")
		log.Println("Info: [WDA URL Getting at second iteration]")
		err := file.Close()
		if err != nil {
			log.WithError(err).Error("Error close file: " + logFilePath)
		}
		return result, nil
	case <-time.After(wdaURLWaitingTime):
		err := file.Close()
		if err != nil {
			log.WithError(err).Error("Error close file: " + logFilePath)
		}

		// Последняя попытка получить WDA URL (при считывании файла целиком)
		bs, err := ioutil.ReadFile(logFilePath)
		if err != nil {
			cmd.Process.Signal(os.Interrupt)
			return "", err
		}

		m := regexpFindServerURL.FindSubmatch(bs)
		if len(m) > 0 {
			result := string(m[1])
			stats.Increment("wda_url.last_chance_find")
			log.Println("Info: [WDA URL Getting at last chance]")
			return result, nil
		}

		cmd.Process.Signal(os.Interrupt)
		stats.Increment("wda_url.not_find")
		return "", fmt.Errorf("Do not wait URL from WDA log:\n [File Path]: %s \n [Content Of Log File]:\n %s ", logFilePath, string(bs))
	}
}

func getDefaultLogFilePath(deviceID string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		log.WithError(err).Error("Error getting current user")
		return "", err
	}

	logFilePath := filepath.Join(usr.HomeDir, "Library/Logs/CoreSimulator", deviceID, "system.log")

	_, err = os.Stat(logFilePath)
	if os.IsNotExist(err) {
		// Если файла с логом не было ранее, значит симмулятор запущен впервые и нужно подождать пока он подготовится
		// Также в это время он создает папку с system.log файлом
		log.Println("Waiting prepare simulator for first run...")
		time.Sleep(15 * time.Second)
	}
	err = os.Truncate(logFilePath, 0)
	if err != nil {
		log.WithError(err).Warnf("Error truncate file %s", logFilePath)
	}
	return logFilePath, nil
}

type ChooseResult struct {
	WDAURL   string
	DeviceID string
}

// ChooseSimulator помогает выбрать нужный симулятор в зависимости от desired capabilities.
func ChooseSimulator(deviceList device.DeviceList, deviceNameFromCapabilities string, iOSVersionFromCapabilities string, wdaPath string, keyPrefix string, stats *statsd.Client) (*ChooseResult, error) {
	var deviceID string
	var url string

	// Получает список устройств
	devices, err := device.GetBootedList()
	if err != nil {
		return nil, err
	}

	if devices != nil {
		log.Println("Founded booted simulators")
		//Проверяем все девайсы на соответствие desired capabilities
		for _, device := range devices {
			customDeviceName := keyPrefix + deviceNameFromCapabilities
			if device.Name == customDeviceName && device.IOSVersion == iOSVersionFromCapabilities {
				log.Println("The required simulator is already running")
				// Запущен нужный симулятор
				deviceID = device.UDID

				//Убиваем процесс WDA
				log.Printf("Killing WDA xcodebuild process [destination:%s]", deviceID)
				KillWDA(deviceID)

				//Запускаем новый процесс
				log.Println("Starting WDA...")
				url, err = StartWDA(deviceID, wdaPath, stats)
				if err != nil {
					return nil, err
				}
			} else {
				// Если нашли девайс с тем же префиксом то сразу закрываем его и останавливаем WDA для него
				if strings.Contains(device.Name, keyPrefix) || keyPrefix == "" {
					//Запущен неверный симмулятор с тем же префиксом, остановим его
					log.Println("Founded simulator do not match desired capabilities")
					log.Printf("Shutting down wrong simulator [%s]", device.UDID)
					err := shutdownSimulator(device.UDID)
					if err != nil {
						return nil, err
					}

					//Убиваем процесс WDA
					log.Printf("Killing WDA xcodebuild process [destination:%s]", device.UDID)
					KillWDA(device.UDID)
				}
			}
		}

		if deviceID == "" {
			//Запускаем нужный симулятор в соответствии с desiredCapabilities
			result, err := bootFromDesiredCapabilities(deviceList, deviceNameFromCapabilities, keyPrefix, iOSVersionFromCapabilities)
			if err != nil {
				return nil, err
			}
			deviceID = result

			log.Println("Starting WDA...")
			url, err = StartWDA(deviceID, wdaPath, stats)
			if err != nil {
				return nil, err
			}
		}
	} else {
		log.Println("Not found booted devices, boot from desired capabilities...")
		result, err := bootFromDesiredCapabilities(deviceList, deviceNameFromCapabilities, keyPrefix, iOSVersionFromCapabilities)
		if err != nil {
			return nil, err
		}
		deviceID = result

		log.Println("Starting WDA...")
		url, err = StartWDA(deviceID, wdaPath, stats)
		if err != nil {
			return nil, err
		}
	}

	time.Sleep(time.Second * 1)
	log.Println("Selected device:", deviceID)

	return &ChooseResult{WDAURL: url, DeviceID: deviceID}, nil
}

// getURLFromSimulatorLog получает URL, на котором стартует WDA.
func getURLFromSimulatorLog(f *os.File) (string) {
	var url string
	buf := make([]byte, 1024)
	for {
		n, _ := f.Read(buf)
		if n != 0 {
			m := regexpFindServerURL.FindSubmatch(buf[:n])
			if len(m) > 0 {
				url = string(m[1])
				break
			}
		} else {
			time.Sleep(time.Millisecond * 20)
		}
	}
	return url
}

func getPathFromLog(s []byte) string {
	var path string
	for i := 0; i < 15; i++ {
		m := regexpFindFilePath.FindSubmatch(s)
		if len(m) > 0 {
			path = string(m[1])
			break
		} else {
			time.Sleep(time.Second)
		}
	}
	return path
}

// bootFromDesiredCapabilities запускает устройство в зависимости от desired capabilities.
func bootFromDesiredCapabilities(devicesList device.DeviceList, deviceName string, keyPrefix string, iOSVersion string) (string, error) {
	var deviceID string
	customDeviceName := keyPrefix + deviceName
	for _, devInfo := range devicesList.Devices[iOSVersion] {
		if devInfo.Name == customDeviceName {
			deviceID = devInfo.UDID
			break
		}
	}
	if deviceID == "" {
		log.Println("Не найдено устройство:" + deviceName + " " + iOSVersion + " в списке устройств")
		log.Printf("Создание нового устройства: [%s / %s]", deviceName, iOSVersion)
		result, err := createSimulator(deviceName, keyPrefix, iOSVersion)
		if err != nil {
			log.WithError(err).Printf("Не удалось создать устройство с параметрами: [%s / %s]", deviceName, iOSVersion)
			return "", err
		}
		deviceID = result
	}

	err := BootSimulator(deviceID)
	if err != nil {
		return "", err
	}
	return deviceID, nil
}

// shutdownSimulator останавливает работу симулятора.
func shutdownSimulator(deviceID string) error {
	command := exec.Command("xcrun", "simctl", "shutdown", deviceID)
	stderr, err := command.StderrPipe()
	if err != nil {
		log.WithError(err).Errorf("Error while shutting down device with id %s", deviceID)
		return err
	}

	if err := command.Start(); err != nil {
		log.WithError(err).Errorf("Error while shutting down device with id %s", deviceID)
		return err
	}

	slurp, _ := ioutil.ReadAll(stderr)
	if string(slurp) != "" {
		log.Println(string(slurp))
	}

	if err := command.Wait(); err != nil {
		log.WithError(err).Errorf("Error while shutting down device with id %s", deviceID)
		return err
	}
	return nil
}

func createSimulator(deviceName string, keyPrefix string, iOSVersion string) (string, error) {
	deviceTypeString := strings.Replace(deviceName, " ", "-", -1)
	deviceType := "com.apple.CoreSimulator.SimDeviceType." + deviceTypeString

	iOSVersionString := strings.Replace(iOSVersion, " ", "-", -1)
	iOSVersionString = strings.Replace(iOSVersionString, ".", "-", -1)
	runtime := "com.apple.CoreSimulator.SimRuntime." + iOSVersionString

	customDeviceName := keyPrefix + deviceName
	out, err := exec.Command("xcrun", "simctl", "create", customDeviceName, deviceType, runtime).Output()
	if err != nil {
		log.WithError(err).Errorf("Error while creating simulator [%s / %s]", deviceName, iOSVersion)
		return "", err
	}

	deviceID := strings.TrimSpace(string(out))
	return deviceID, nil
}
