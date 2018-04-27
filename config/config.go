package config

import (
	"encoding/json"
	"errors"
	"github.com/Sirupsen/logrus"
	"log"
	"os"
)

type Config struct {
	Server       Server            `json:"server"`
	Logger       Logger            `json:"logger"`
	IOS          IOS               `json:"ios"`
	Grid         Grid              `json:"grid"`
	Capabilities map[string]string `json:"capabilities"`
	Statsd       Statsd            `json:"statsd"`
	Video        Video             `json:"video"`
}

type Server struct {
	Port int `json:"port"`
}

type Logger struct {
	Level logrus.Level `json:"level"`
}

type WDA struct {
	Path         string `json:"path"`
	DevicePrefix string `json:"devicePrefix"`
}

type IOS struct {
	WDA     WDA    `json:"wda"`
	AppPath string `json:"appPath"`
}

type Grid struct {
	Host string `json:"host"`
}

type Statsd struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Prefix   string `json:"prefix"`
	Enable   bool   `json:"enable"`
}

type S3 struct {
	Endpoint  *string `json:"endpoint,omitempty"`
	AccessKey string  `json:"access_key"`
	SecretKey string  `json:"secret_key"`
	Bucket    string  `json:"bucket"`
	Region    string  `json:"region"`
}

type Video struct {
	Enable bool `json:"enable"`
	S3     S3   `json:"s3"`
}

func New() *Config {
	return &Config{}
}

func (c *Config) LoadFromFile(path string) error {
	log.Printf("Loaded config: %s", path)
	if path == "" {
		return errors.New("Empty configuration file path")
	}

	configFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)

	return jsonParser.Decode(&c)
}
