package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mikeewhite/ship-locator/backend/internal/core/services/shipsrcsrv"
	"github.com/mikeewhite/ship-locator/backend/internal/handlers/graphql/searchgraph"
	"github.com/mikeewhite/ship-locator/backend/internal/handlers/kafka/consumer"
	"github.com/mikeewhite/ship-locator/backend/internal/repositories/shipsrc/elasticsearch"
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

	// initialise the metrics client
	metricsClient := metrics.New(*cfg)
	go func() {
		if err := metricsClient.Serve(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			clog.Errorf("metrics client stopped due to error: %s", err.Error())
		}
	}()

	// initialise the ship search service
	searchRepo, err := elasticsearch.New(ctx, *cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to initialise Elasticsearch repository: %s", err.Error()))
	}
	searchService := shipsrcsrv.New(searchRepo)

	// start the ship event consumer
	shipEventConsumer, err := consumer.NewShipEventConsumer(*cfg, searchService, metricsClient)
	if err != nil {
		panic(fmt.Sprintf("failed to initialise ship event consumer: %s", err.Error()))
	}
	defer shipEventConsumer.Shutdown()
	go func() {
		if err := shipEventConsumer.Read(ctx); err != nil && !errors.Is(err, context.Canceled) {
			clog.Errorf("kafka consumer stopped due to error: %s", err.Error())
		}
	}()

	server, err := searchgraph.New(*cfg, searchService)
	if err != nil {
		panic(fmt.Sprintf("failed to initialise GraphQL server: %s", err.Error()))
	}
	if err := server.Serve(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
