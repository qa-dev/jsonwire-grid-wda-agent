package config

const (
	// DefaultPort задает порт по-умолчанию, на котором запущен сервер.
	DefaultPort = 8890
)

// Server задает настройки web-сервиса.
type Server struct {
	Port int `json:"port"`
}

// Logger задает настройки логирования.
type Logger struct {
	Level     string `json:"level"`
	UseStderr bool   `json:"use_stderr"`
	Tag       string `json:"tag"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
}

// Common задает общие настройки приложения.
type Common struct {
	LocalConfigPath string `json:"local_config_path"`
}

// Statsd задает настройки для работы с сервером statsd.
type Statsd struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Prefix   string `json:"prefix"`
	Enable   bool   `json:"enable"`
}

// Client задает настройки для API клиентов.
type Client struct {
	URL          string  `json:"url"`
	Timeout      float32 `json:"timeout"`
	MaxIdleConns int     `json:"max_idle_conns"`
}

type Validator struct {
	Path string `json:"path"`
}

// BaseConfig задает общие для всех сервисов настройки.
type BaseConfig struct {
	Common    Common    `json:"common"`
	Server    Server    `json:"server"`
	Logger    Logger    `json:"logger"`
	Statsd    Statsd    `json:"statsd"`
	Validator Validator `json:"validator"`
}
