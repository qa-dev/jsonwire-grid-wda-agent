package clients

import (
	"net/http"
	"time"

	log "go.avito.ru/gl/core/logger"
)

// HTTPClient задает интерфейс http-клиента.
// Совместим с http.Client.
type HTTPClient interface {
	Do(request *http.Request) (*http.Response, error)
}

// NewHTTPClient возвращает новый проиницилизированный http.Client.
// Настраиваются timeout на получение ответа и максимальный размер пула keep-alive соединений.
func NewHTTPClient(timeout float32, maxIdleConns int) *http.Client {
	timeoutDurationMS := time.Duration(timeout*1000) * time.Millisecond

	log.Infof(`Create http client with timeout %v.`, timeoutDurationMS)

	return &http.Client{
		Timeout: timeoutDurationMS,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxIdleConns,
		},
	}
}
