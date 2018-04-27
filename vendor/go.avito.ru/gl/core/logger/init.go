package logger

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/fluent/fluent-logger-golang/fluent"
)

// InitLogger настраивает механизм логирования.
func InitLogger(logLevel string,
	useStderr bool,
	tag string,
	host string,
	port int,
	fields map[string]string) error {
	level, err := log.ParseLevel(logLevel)

	if err != nil {
		return fmt.Errorf("logger level error: %v", err)
	}

	log.SetLevel(level)

	if useStderr {
		log.SetOutput(os.Stderr)
		log.SetFormatter(&log.JSONFormatter{})

		return nil
	}

	var fluentConfig fluent.Config

	if strings.HasPrefix(host, "unix://") {
		fluentConfig = fluent.Config{
			FluentNetwork:    "unix",
			FluentSocketPath: host[7:],
			MaxRetry:         math.MaxInt32}
	} else {
		fluentConfig = fluent.Config{
			FluentNetwork: "tcp",
			FluentHost:    host,
			FluentPort:    port,
			MaxRetry:      math.MaxInt32}
	}

	fluentHook, err := NewFluentHook(fluentConfig, tag, fields)

	if err != nil {
		return err
	}

	log.AddHook(fluentHook)

	log.SetOutput(ioutil.Discard)

	return nil
}
