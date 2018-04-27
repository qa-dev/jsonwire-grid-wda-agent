package main

import (
	contx "context"
	"fmt"
	"github.com/qa-dev/jsonwire-grid-wda-agent/config"
	"github.com/qa-dev/jsonwire-grid-wda-agent/grid"
	"github.com/qa-dev/jsonwire-grid-wda-agent/handlers"
	"github.com/qa-dev/jsonwire-grid-wda-agent/logger"
	"github.com/qa-dev/jsonwire-grid-wda-agent/wda"
	"go.avito.ru/gl/core/context"
	"go.avito.ru/gl/core/metrics"
	"go.avito.ru/gl/core/middlewares"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

func main() {
	cfg := config.New()
	err := cfg.LoadFromFile(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatalf("Problem in loading config from file, %s", err.Error())
	}
	logger.Init(cfg.Logger)

	statsd := cfg.Statsd
	stats, err := metrics.NewStatsd(statsd.Host, statsd.Port, statsd.Protocol, statsd.Prefix, statsd.Enable)
	if err != nil {
		log.Fatalf("Statsd create socked error: %s", err.Error())
	}
	ctx := context.NewBaseContext("go-iosd", stats)

	if cfg.Grid.Host != "" {
		gridC := grid.NewRegisterAndCheckData(cfg.Grid.Host)
		c := make(chan error)
		go func() {
			err := <-c
			if err != nil {
				log.Fatal(err.Error())
			}
		}()
		gridC.RegisterAndCheck(cfg.Server.Port, cfg.Capabilities, c)
	}

	target := &url.URL{Scheme: "http", Host: "127.0.0.1:" + cfg.IOS.WDA.DevicePrefix}
	proxy := wda.NewProxy(target)

	http.Handle("/session", handlers.NewCreateSessionHandler(cfg, proxy, stats))
	http.Handle("/session/", handlers.NewSessionHandler(cfg, proxy))
	http.Handle("/health/", handlers.NewHealthHandler())
	http.Handle("/", proxy)

	mux := middlewares.NewLogMiddleware(ctx).Log(http.DefaultServeMux)

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	server := &http.Server{Addr: fmt.Sprintf(":%v", cfg.Server.Port), Handler: mux}
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			log.Fatalf("Listen serve error, %s", err.Error())
		}
	}()

	<-stop

	log.Println("Shutting down the server...")

	con, _ := contx.WithTimeout(contx.Background(), 5*time.Minute)
	server.Shutdown(con)

	log.Println("Server gracefully stopped")
}
