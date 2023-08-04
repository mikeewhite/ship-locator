package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mikeewhite/ship-locator/backend/internal/handlers/kafka"
	"github.com/mikeewhite/ship-locator/backend/pkg/clog"
	"github.com/mikeewhite/ship-locator/backend/pkg/config"
)

func main() {
	defer clog.Info("service stopped")
	defer clog.Flush()

	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("error on loading config: %s", err.Error()))
	}

	consumer, err := kafka.NewConsumer(*cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to initialise Kafka consumer: %s", err.Error()))
	}
	defer consumer.Shutdown()

	ctx, cancel := context.WithCancel(context.Background())
	gracefulShutdownOnSignal(cancel)
	if err = consumer.Read(ctx); err != nil {
		clog.Errorf("kafka consumer stopped due to error: %s", err.Error())
	}
}

func gracefulShutdownOnSignal(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-signals
		clog.Infow("shutting down",
			"signal", s.String())
		cancel()
	}()
}
