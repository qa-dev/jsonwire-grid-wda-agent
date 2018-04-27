// Package handlers предназначен для перехвата http запросов и их обслуживания.
package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/jsonwire-grid-wda-agent/command"
	"github.com/qa-dev/jsonwire-grid-wda-agent/config"
	"github.com/qa-dev/jsonwire-grid-wda-agent/device"
	"github.com/qa-dev/jsonwire-grid-wda-agent/wda"
	"io"
	"io/ioutil"
	"net/http"
	"fmt"
	"gopkg.in/alexcesaro/statsd.v2"
)

// DesiredCapabilities описывает входные capabilities теста.
type DesiredCapabilities struct {
	Capabilities Capabilities `json:"desiredCapabilities"`
}

// Capabilities описывает параметры для выбора устройства, платформы и приложения.
type Capabilities struct {
	DeviceName   string `json:"deviceName"`
	IOSVersion   string `json:"iOSVersion"`
	AppPath      string `json:"appPath"`
	BundleID     string `json:"bundleId"`
	PlatformName string `json:"platformName"`
}

// CreateSessionHandler содержит в себе конфиг приложения.
type CreateSessionHandler struct {
	cfg   *config.Config
	proxy *wda.Proxy
	stats *statsd.Client
}

// NewCreateSessionHandler описывает конструктор со структурой CreateSessionHandler.
func NewCreateSessionHandler(cfg *config.Config, pr *wda.Proxy, stats *statsd.Client) *CreateSessionHandler {
	return &CreateSessionHandler{
		cfg:   cfg,
		proxy: pr,
		stats: stats,
	}
}

// ServeHTTP перехватывает http запрос /session, передает управление на подготовку устройства к запуску теста.
// В случае не POST запросов, просто проксирует их.
func (h *CreateSessionHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Body.Close()
	var desiredCapabilities DesiredCapabilities
	err = json.Unmarshal(body, &desiredCapabilities)
	if err != nil {
		log.WithError(err).Error("Error unmarshal body in desired capabilities")
		fErr := fmt.Errorf("Error unmarshal body in desired capabilities: %v : ", err)
		http.Error(rw, fErr.Error(), http.StatusInternalServerError)
		return
	}

	if desiredCapabilities.Capabilities.DeviceName == "" {
		err = errors.New("Не задан ключ DeviceName в desired capabilities")
		log.Println(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if desiredCapabilities.Capabilities.IOSVersion == "" {
		err = errors.New("Не задан ключ iOSVersion в desired capabilities")
		log.Println(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if desiredCapabilities.Capabilities.BundleID == "" {
		err = errors.New("Не задан ключ BundleID в desired capabilities")
		log.Println(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var url string
	if r.Method == http.MethodPost {
		url, err = h.PrepareForTest(desiredCapabilities.Capabilities, h.stats)
		h.proxy.SetUrl(url)
		if err != nil {
			log.Println(err.Error())
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("WDA URL:", url)
	}

	r.URL.Host = url
	r.URL.Scheme = "http"
	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	copyHeader(rw.Header(), resp.Header)
	rw.WriteHeader(resp.StatusCode)
	io.Copy(rw, resp.Body)
}

// PrepareForTest выполняет подготовку устройства перед запуском теста.
// 1. Остановка WDA.
// 2. Получить список всех устройств.
// 3. Выбрать симулятор для работы.
// 4. Удалить приложение.
// 5. Установить приложение.
// 6. Запуск WDA.
func (h *CreateSessionHandler) PrepareForTest(capabilities Capabilities, stats *statsd.Client) (string, error) {
	log.Println("Get devices list")
	deviceList, err := device.GetDeviceList()
	if err != nil {
		return "", err
	}

	log.Println("Choose simulator...")
	result, err := command.ChooseSimulator(*deviceList, capabilities.DeviceName, capabilities.IOSVersion, h.cfg.IOS.WDA.Path, h.cfg.IOS.WDA.DevicePrefix, stats)
	if err != nil {
		return "", err
	}
	deviceID := result.DeviceID

	log.Println("Uninstall App")
	err = command.UninstallApp(deviceID, capabilities.BundleID)
	if err != nil {
		return "", err
	}

	log.Println("Install App")
	var appPath string
	if capabilities.AppPath != "" {
		appPath = capabilities.AppPath
	} else {
		appPath = h.cfg.IOS.AppPath + "/" + capabilities.BundleID
	}
	err = command.InstallApp(deviceID, appPath)
	if err != nil {
		return "", err
	}

	if h.cfg.Video.Enable {
		err := command.StartVideo(deviceID)
		if err != nil {
			log.Printf("Error starting video: %v on deviceID:[%s]", err, deviceID)
		}
	}
	return result.WDAURL, nil
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
