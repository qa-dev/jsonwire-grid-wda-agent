package middlewares

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

// LoggedResponseWriter оборачивает http.ResponseWriter
// Поддерживает возможность чтения ранее записанного статуса ответа
type LoggedResponseWriter struct {
	responseWriter http.ResponseWriter
	status         int
}

// Header возвращает результат вызова метода Header() оборачиваемого http.ResponseWriter
func (lrw *LoggedResponseWriter) Header() http.Header {
	return lrw.responseWriter.Header()
}

// Write возвращает результат вызова метода Write() оборачиваемого http.ResponseWriter
func (lrw *LoggedResponseWriter) Write(data []byte) (int, error) {
	return lrw.responseWriter.Write(data)
}

// WriteHeader возвращает результат вызова метода WriteHeader() оборачиваемого http.ResponseWriter
func (lrw *LoggedResponseWriter) WriteHeader(status int) {
	lrw.status = status
	lrw.responseWriter.WriteHeader(status)
}

// Status возвращает ранее записанный статус ответа
func (lrw *LoggedResponseWriter) Status() int {
	return lrw.status
}

// ResponseWriter возвращает оборочиваемый ResponseWriter
func (lrw *LoggedResponseWriter) ResponseWriter() http.ResponseWriter {
	return lrw.responseWriter
}

// Hijack имплементирует интерфейс Hijacker для http.ResponseWriter
func (lrw *LoggedResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := lrw.responseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("ResponseWriter doesn't support Hijacker interface")
	}
	return hijacker.Hijack()
}
