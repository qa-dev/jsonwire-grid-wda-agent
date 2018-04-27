package middlewares

import (
	"net/http"

	log "go.avito.ru/gl/core/logger"
)

// PostChecker проверяет, что запрос - HTTP POST.
func PostChecker(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		requestID := req.Header.Get("X-Request-Id")
		if req.Method != "POST" {
			log.WithFields(log.Fields{
				"request_id": requestID,
			}).Warnf("HTTP method is not allowed: %v.", req.Method)
			http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

			return
		}
		handler.ServeHTTP(rw, req)
	})
}
