package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mikeewhite/ship-locator/backend/internal/core/services/collectorsrv"
	"github.com/mikeewhite/ship-locator/backend/internal/handlers/aishdl"
	"github.com/mikeewhite/ship-locator/backend/internal/handlers/msghdl"
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

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		s := <-signals
		clog.Infow("shutting down",
			"signal", s.String())
		cancel()
	}()

	producer := msghdl.NewKafkaPublisher()
	defer producer.Shutdown()
	service := collectorsrv.New(ctx, producer)
	defer service.Shutdown()
	listener, err := aishdl.NewWebSocketListener(*cfg, service)
	if err != nil {
		panic(fmt.Sprintf("failed to initialise websocket listener: %s", err.Error()))
	}
	defer listener.Shutdown()
	if err := listener.Listen(ctx); err != nil {
		clog.Errorf("websocket listener stopped due to error: %s", err.Error())
	}
}
