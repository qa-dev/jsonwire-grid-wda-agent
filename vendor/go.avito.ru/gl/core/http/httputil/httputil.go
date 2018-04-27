package httputil

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=- //

//GetClientIPFromRequest пытается получить IP клиента на основе заголовков
func GetClientIPFromRequest(r *http.Request) (ip string) {
	ip = r.Header.Get("X-Real-Ip")

	if ip != "" {
		return ip
	}

	ip = r.Header.Get("X-Forwarded-For")

	if ip != "" {
		ipList := strings.Split(ip, ",")

		return strings.TrimSpace(ipList[len(ipList)-1])
	}

	return r.RemoteAddr
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=- //

// FloatParam достает параметр float из урла
func FloatParam(req *http.Request, name string) (float64, error) {
	strParam, err := StrParam(req, name)
	if err != nil {
		return 0, err
	}

	floatParam, err := strconv.ParseFloat(strParam, 64)
	if err != nil {
		return 0, err
	}

	return floatParam, nil
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=- //

// IntParam достает параметр int из урла
func IntParam(req *http.Request, name string) (int, error) {
	strParam, err := StrParam(req, name)
	if err != nil {
		return 0, err
	}

	intParam, err := strconv.Atoi(strParam)
	if err != nil {
		return 0, err
	}

	return intParam, nil
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=- //

// StrParam достает параметр str из урла
func StrParam(req *http.Request, name string) (string, error) {
	param := req.URL.Query().Get(name)
	if param == "" {
		return "", errors.New("no such param " + name)
	}

	return param, nil
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=- //

// FormFloatParam достает параметр float из формы
func FormFloatParam(req *http.Request, name string) (float64, error) {
	param, err := FormStrParam(req, name)
	if err != nil {
		return 0, errors.New("no such param " + name)
	}

	floatParam, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return 0, err
	}

	return floatParam, nil
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=- //

// FormIntParam достает параметр int из формы
func FormIntParam(req *http.Request, name string) (int, error) {
	param, err := FormStrParam(req, name)
	if err != nil {
		return 0, errors.New("no such param " + name)
	}

	intParam, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}

	return intParam, nil
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=- //

// FormStrParam достает параметр str из формы
func FormStrParam(req *http.Request, name string) (string, error) {
	req.ParseForm()
	param := req.FormValue(name)
	if param == "" {
		return "", errors.New("no such param " + name)
	}
	return param, nil
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=- //
