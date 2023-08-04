package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mikeewhite/ship-locator/backend/internal/core/services/collectorsrv"
	"github.com/mikeewhite/ship-locator/backend/internal/handlers/kafka"
	"github.com/mikeewhite/ship-locator/backend/internal/handlers/websocket"
	"github.com/mikeewhite/ship-locator/backend/pkg/clog"
	"github.com/mikeewhite/ship-locator/backend/pkg/config"
)

func main() {
	defer clog.Info("collector stopped")
	defer clog.Flush()

	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("error on loading config: %s", err.Error()))
	}

	producer, err := kafka.NewProducer(*cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to initialise Kafka producer: %s", err.Error()))
	}
	defer producer.Shutdown()

	ctx, cancel := context.WithCancel(context.Background())
	gracefulShutdownOnSignal(cancel)
	service := collectorsrv.New(ctx, producer)
	defer service.Shutdown()
	listener, err := websocket.NewWebSocketListener(*cfg, service)
	if err != nil {
		panic(fmt.Sprintf("failed to initialise websocket listener: %s", err.Error()))
	}
	defer listener.Shutdown()
	if err := listener.Listen(ctx); err != nil {
		clog.Errorf("websocket listener stopped due to error: %s", err.Error())
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
