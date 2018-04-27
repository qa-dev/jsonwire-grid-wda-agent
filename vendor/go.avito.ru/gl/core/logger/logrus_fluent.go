package logger

import (
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/fluent/fluent-logger-golang/fluent"
)

const (
	// DefaultTag stores default fluent tag.
	DefaultTag = "app"
)

// FluentHook to send logs via fluentd.
type FluentHook struct {
	Fluent     *fluent.Fluent
	DefaultTag string
	Fields     map[string]string
}

// NewFluentHook creates a new hook to send to fluentd.
func NewFluentHook(config fluent.Config, defaultTag string, fields map[string]string) (*FluentHook, error) {
	logger, err := fluent.New(config)
	if err != nil {
		return nil, err
	}

	if defaultTag == "" {
		defaultTag = DefaultTag
	}

	return &FluentHook{
		Fluent:     logger,
		DefaultTag: defaultTag,
		Fields:     fields,
	}, nil
}

// Fire implements logrus.Hook interface Fire method.
func (f *FluentHook) Fire(entry *logrus.Entry) error {
	for name, value := range f.Fields {
		entry.Data[name] = value
	}

	msg := f.buildMessage(entry)
	tag := f.DefaultTag
	rawTag, ok := entry.Data["tag"]
	if ok {
		tag = fmt.Sprint(rawTag)
	}
	f.Fluent.Post(tag, msg)

	return nil
}

// Levels implements logrus.Hook interface Levels method.
func (f *FluentHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (f *FluentHook) buildMessage(entry *logrus.Entry) map[string]interface{} {
	data := make(map[string]interface{})

	for k, v := range entry.Data {
		if k == "tag" {
			continue
		}
		data[k] = v
	}
	data["message"] = entry.Message

	level := entry.Level.String()

	data["type"] = strings.ToUpper(string(level[0]))

	return data
}
