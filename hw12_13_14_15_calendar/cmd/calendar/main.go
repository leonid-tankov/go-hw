package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	conf := config.NewConfig(configFile)
	logg := logger.New(conf.Logger.Level, os.Stdout)
	store := storage.NewStorageByType(conf, logg)
	calendar := app.New(logg, store)

	server := internalhttp.NewServer(conf.HTTP.Host, conf.HTTP.Port, logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("App Calendar is running...")

	if err := server.Start(); !errors.Is(err, http.ErrServerClosed) {
		logg.Error("failed to start http server: %v", err)
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
