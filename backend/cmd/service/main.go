package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mikeewhite/ship-locator/backend/internal/core/services/shipsrv"
	"github.com/mikeewhite/ship-locator/backend/internal/handlers/graphql"
	"github.com/mikeewhite/ship-locator/backend/internal/handlers/kafka"
	"github.com/mikeewhite/ship-locator/backend/internal/repositories/postgres"
	"github.com/mikeewhite/ship-locator/backend/pkg/clog"
	"github.com/mikeewhite/ship-locator/backend/pkg/config"
	"github.com/mikeewhite/ship-locator/backend/pkg/metrics"
)

func main() {
	defer clog.Info("service stopped")
	defer clog.Flush()

	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("error on loading config: %s", err.Error()))
	}

	ctx, cancel := context.WithCancel(context.Background())
	gracefulShutdownOnSignal(cancel)

	metricsClient := metrics.New(*cfg)
	go func() {
		if err := metricsClient.Serve(ctx); err != nil && err != http.ErrServerClosed {
			clog.Errorf("metrics client stopped due to error: %s", err.Error())
		}
	}()

	repo, err := postgres.NewPostgres(ctx, *cfg, metricsClient)
	if err != nil {
		panic(fmt.Sprintf("failed to initialise Postgres repository: %s", err.Error()))
	}
	service := shipsrv.New(repo)
	consumer, err := kafka.NewConsumer(*cfg, shipsrv.New(repo), metricsClient)
	if err != nil {
		panic(fmt.Sprintf("failed to initialise Kafka consumer: %s", err.Error()))
	}
	defer consumer.Shutdown()
	go func() {
		if err := consumer.Read(ctx); err != nil && err != context.Canceled {
			clog.Errorf("kafka consumer stopped due to error: %s", err.Error())
		}
	}()

	server, err := graphql.New(*cfg, service)
	if err != nil {
		panic(fmt.Sprintf("failed to initialise GraphQL server: %s", err.Error()))
	}
	if err := server.Serve(ctx); err != nil && err != http.ErrServerClosed {
		clog.Errorf("graphql server stopped due to error: %s\n", err.Error())
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
