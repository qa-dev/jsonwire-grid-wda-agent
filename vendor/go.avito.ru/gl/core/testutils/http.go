package testutils

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// MockHTTPClient реализует интерфейс http-клиента для отдачи заданного ответа.
type MockHTTPClient struct {
	Response *http.Response
	Error    error
}

// NewMockHTTPClient создает новый тестовый клиент с заданными статусом и телом ответа.
func NewMockHTTPClient(statusCode int, body string) *MockHTTPClient {
	return NewMockHTTPClientWithError(statusCode, body, nil)
}

// NewMockHTTPClientWithError создает новый тестовый клиент который при вызове возвращает ошибку.
func NewMockHTTPClientWithError(statusCode int, body string, expectedError error) *MockHTTPClient {
	return &MockHTTPClient{
		Response: MakeResponse(statusCode, body),
		Error:    expectedError,
	}
}

// SetError устанавливает ошибку возвращаемую при вызове Do
func (c *MockHTTPClient) SetError(expectedError error) {
	c.Error = expectedError
}

// Do выполняет запрос и возвращает ответ.
func (c *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.Response, c.Error
}

// MakeResponse создает новый ответ с заданным статусом и телом.
func MakeResponse(statusСode int, body string) *http.Response {
	resp := new(http.Response)

	resp.StatusCode = statusСode
	resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(body)))

	return resp
}
