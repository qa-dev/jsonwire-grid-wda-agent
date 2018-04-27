// Package grid предназначен для работы с jsonwire-grid, регистрации в нём и проверки соединения с ним.
package grid

import (
	"bytes"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"go.avito.ru/gl/core/clients"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"fmt"
)

// ResponseGrid описывает структуру ответа от jsonwire-grid.
// Success = true, если хост зарегистрирован, false в противном случае.
type ResponseGrid struct {
	Success bool `json:"success"`
}

// RegisterRequestData описывает структуру запроса в jsonwire-grid.
// Конфигурация Configuration состоит из параметров host и port данного сервиса.
// Платформа всегда по умолчанию - WDA и задается в Capabilities.
type RegisterRequestData struct {
	Configuration Configuration  `json:"configuration"`
	Capabilities  []Capabilities `json:"capabilities"`
}

// Configuration состоит из параметров host и port данного сервиса.
type Configuration struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// Capabilities описывает платформу. (по умолчанию WDA)
type Capabilities map[string]string

// RegisterAndCheckData описывает структуру, состоящую из адреса хоста jsonwire-grid и http клиента.
type RegisterAndCheckData struct {
	gridHost string
	client   *http.Client
}

// NewRegisterAndCheckData описывает конструктор со структурой RegisterAndCheckData
func NewRegisterAndCheckData(gridHost string) *RegisterAndCheckData {
	return &RegisterAndCheckData{
		gridHost: gridHost,
		client:   clients.NewHTTPClient(60, 10),
	}
}

// register описывает регистрацию в jsonwire-grid.
func (h *RegisterAndCheckData) register(host string, port int, capabilities map[string]string) error {
	log.Println("Register in grid...")
	var reqData = &RegisterRequestData{
		Configuration: Configuration{Host: host, Port: port},
		Capabilities:  []Capabilities{Capabilities(capabilities)},
	}

	b, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	var jsonStr = []byte(b)
	reqURL := url.URL{
		Scheme: "http",
		Path:   "/grid/register",
		Host:   h.gridHost,
	}
	req, err := http.NewRequest("POST", reqURL.String(), bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Println("Response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("Response Body:", string(body))
	return nil
}

// checkConnection метод предназначен для проверки соединения с grid.
func (h *RegisterAndCheckData) checkConnection(host string, port int) (bool, error) {
	reqURL := url.URL{
		Scheme:   "http",
		Path:     "/grid/api/proxy",
		RawQuery: "id=http://" + host + ":" + strconv.Itoa(port),
		Host:     h.gridHost,
	}

	resp, err := http.Get(reqURL.String())
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var ResponseGrid ResponseGrid
	err = json.Unmarshal(body, &ResponseGrid)
	if err != nil {
		fErr := fmt.Errorf("Error unmarshal response from grid: (checkConnection): %v: ", err)
		log.WithError(err).Error("Error unmarshal response from grid: URL:" + reqURL.String())
		return false, fErr
	}

	return ResponseGrid.Success, nil
}

// RegisterAndCheck метод предназначен для регистрации и отслеживания соединения с jsonwire-grid каждые 5 секунд.
// Если соединение с jsonwire-grid потеряно, то регистрация происходит снова.
func (h *RegisterAndCheckData) RegisterAndCheck(port int, capabilities map[string]string, c chan error) {
	host, err := getIpv4()
	if err != nil {
		c <- err
		return
	}

	err = h.register(host, port, capabilities)
	if err != nil {
		c <- err
		return
	}

	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for range ticker.C {
			result, err := h.checkConnection(host, port)
			if err != nil {
				c <- err
			}
			if result == false {
				log.Println("Cannot find proxy in registry")
				log.Println("Register in grid again...")
				err := h.register(host, port, capabilities)
				if err != nil {
					c <- err
				}
			}
		}
	}()

	return
}

// getIpv4 метод получает ip адрес, на котором стартует данный сервис.
func getIpv4() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", nil
}
