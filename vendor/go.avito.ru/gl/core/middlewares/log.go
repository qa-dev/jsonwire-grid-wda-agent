package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/pborman/uuid"

	"go.avito.ru/gl/core/context"
	log "go.avito.ru/gl/core/logger"
)

const (
	// RequestIDHeader хранит название HTTP-заголовка для передачи ID запроса.
	RequestIDHeader = "X-Request-Id"
)

// LogMiddleware оборачивает добавляет логгирование к обработчикам запросов.
type LogMiddleware struct {
	ctx *context.BaseContext
}

// NewLogMiddleware возвращает новый LogMiddleware.
func NewLogMiddleware(ctx *context.BaseContext) *LogMiddleware {
	return &LogMiddleware{
		ctx: ctx,
	}
}

// Log оборачивает http.Handler для логирования времени выполнения.
func (m *LogMiddleware) Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		requestID := req.Header.Get(RequestIDHeader)

		if requestID == "" {
			requestID = generateRequestID()
			req.Header.Set(RequestIDHeader, requestID)
		}

		defer func() {
			if err := recover(); err != nil {
				log.WithFields(log.Fields{
					"request_id": requestID,
				}).Fatalf("Panic: %+v\n%s", err, debug.Stack())
			}
		}()

		path := req.URL.Path
		handlerName := strings.Replace(path, "/", "", -1)

		requestTimer := m.ctx.Statsd.NewTiming()

		lrw := &LoggedResponseWriter{responseWriter: resp}

		handler.ServeHTTP(lrw, req)

		requestTimer.Send(fmt.Sprintf(
			"service.api.%v_%v.request_time.%v",
			handlerName,
			req.Method,
			lrw.Status()))

		// Example: 200 POST /rec/ (127.0.0.1) 1.46ms
		log.WithFields(log.Fields{
			"request_id": requestID,
		}).Infof("%v %v %v (%v) %.2fms",
			lrw.Status(),
			req.Method,
			path,
			req.RemoteAddr,
			requestTimer.Duration().Seconds()*1000)
	})
}

func generateRequestID() string {
	return uuid.NewUUID().String()
}
