package context

import (
	"gopkg.in/alexcesaro/statsd.v2"
)

// BaseContext хранит базовые глобальные структуры приложения.
type BaseContext struct {
	Hostname string
	Statsd   *statsd.Client
}

// NewBaseContext возвращает новый BaseContext, проинициализированный на основе переданной конфигурации
func NewBaseContext(hostname string, statsd *statsd.Client) *BaseContext {
	return &BaseContext{
		Hostname: hostname,
		Statsd:   statsd,
	}
}
