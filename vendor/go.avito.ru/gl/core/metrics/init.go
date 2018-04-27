package metrics

import (
	"fmt"
	"strings"

	"gopkg.in/alexcesaro/statsd.v2"

	log "go.avito.ru/gl/core/logger"
)

// NewStatsd возвращает новый настроенный клиент statsd
func NewStatsd(host string, port int, protocol string, prefix string, enable bool) (*statsd.Client, error) {
	protocol = strings.ToLower(protocol)
	muted := !enable

	log.Infof(
		`Create statsd client to %v:%v via %v with prefix "%v", muted is %v.`,
		host, port, protocol, prefix, muted)

	client, err := statsd.New(
		statsd.Address(fmt.Sprintf("%v:%v", host, port)),
		statsd.Prefix(prefix),
		statsd.Network(protocol),
		statsd.Mute(muted))

	if err != nil {
		return nil, err
	}

	log.Info("Statsd client was created.")

	return client, nil
}
