package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/danielfurman/ports-microservices/internal/portssvc"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go cancelOnShutdownSignal(cancel, newShutdownSignalCh())

	var cfg portssvc.Config
	if err := env.Parse(&cfg); err != nil {
		logrus.WithError(err).Fatal("Failed to read config from environment")
	}

	server := portssvc.NewServer(cfg)
	if err := server.Serve(ctx); err != nil {
		logrus.WithError(err).Fatal("Ports server stopped")
	}

	logrus.Info("Ports service finished successfully")
}

func newShutdownSignalCh() chan os.Signal {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	return signalCh
}

func cancelOnShutdownSignal(cancel context.CancelFunc, signalCh chan os.Signal) {
	s := <-signalCh
	logrus.WithField("signal", s).Info("Received shutdown signal - stopping application")
	cancel()
}
