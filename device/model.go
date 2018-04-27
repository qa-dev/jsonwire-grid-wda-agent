package device

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os/exec"
	"strings"
)

type DeviceInfo struct {
	State        string `json:"state"`
	Availability string `json:"availability"`
	Name         string `json:"name"`
	UDID         string `json:"udid"`
}

type DeviceList struct {
	Devices map[string][]DeviceInfo `json:"devices"`
}

// DeviceInfo структура возвращаемой информации об устройстве
type DeviceInfoResult struct {
	State      string
	Name       string
	UDID       string
	IOSVersion string
}

func GetBootedList() ([]*DeviceInfoResult, error) {
	devVer, err := GetDeviceList()
	if err != nil {
		return nil, err
	}

	var bootedDevices []*DeviceInfoResult
	for iOSVersion, devList := range devVer.Devices {
		for _, devInfo := range devList {
			if devInfo.State == "Booted" {
				bootedDevices = append(
					bootedDevices,
					&DeviceInfoResult{
						Name:       devInfo.Name,
						State:      devInfo.State,
						UDID:       devInfo.UDID,
						IOSVersion: iOSVersion,
					})
			}
		}
	}

	if bootedDevices == nil {
		log.Println("Not found booted devices")
	}
	return bootedDevices, nil
}

// GetDeviceList предназначен для получения списка устройств.
func GetDeviceList() (*DeviceList, error) {
	args := []string{"xcrun", "simctl", "list", "-j", "devices"}
	argsStr := strings.Join(args, " ")
	out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf(`%s`, argsStr)).Output()
	if err != nil {
		log.WithError(err).Error("Error getting booted simulator")
		return nil, err
	}

	var deviceList DeviceList
	err = json.Unmarshal(out, &deviceList)
	if err != nil {
		fErr := fmt.Errorf("Error unmarshal json to DeviceList struct: (GetDeviceList) : %v : ", err)
		log.WithError(err).Error("Error unmarshal json to DeviceList struct")
		return nil, fErr
	}

	return &deviceList, nil
}
