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
	"github.com/mikeewhite/ship-locator/backend/internal/handlers/kafka/consumer"
	"github.com/mikeewhite/ship-locator/backend/internal/handlers/kafka/producer"
	"github.com/mikeewhite/ship-locator/backend/internal/repositories/shipsrc/elasticsearch"

	"github.com/mikeewhite/ship-locator/backend/internal/core/services/shipsrv"
	"github.com/mikeewhite/ship-locator/backend/internal/handlers/graphql"
	"github.com/mikeewhite/ship-locator/backend/internal/repositories/postgres"
	"github.com/mikeewhite/ship-locator/backend/pkg/clog"
	"github.com/mikeewhite/ship-locator/backend/pkg/config"
	"github.com/mikeewhite/ship-locator/backend/pkg/metrics"
	"github.com/mikeewhite/ship-locator/backend/pkg/tracing"
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

	traceProvider, err := tracing.NewTraceProvider(ctx, *cfg, "ship-data-service")
	if err != nil {
		panic(fmt.Sprintf("error on initialising trace provider: %s", err.Error()))
	}
	defer traceProvider.Shutdown(ctx)

	metricsClient := metrics.New(*cfg)
	go func() {
		if err := metricsClient.Serve(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			clog.Errorf("metrics client stopped due to error: %s", err.Error())
		}
	}()

	// initialise the ship search service
	// TODO - this should be moved to the dedicated search microservice (it's currently here so that
	// the GraphQL server can access it)
	searchRepo, err := elasticsearch.New(ctx, *cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to initialise Elasticsearch repository: %s", err.Error()))
	}
	searchService := shipsrcsrv.New(searchRepo)

	// initialise the ship data service
	repo, err := postgres.NewPostgres(ctx, *cfg, metricsClient)
	defer repo.Shutdown(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to initialise Postgres repository: %s", err.Error()))
	}
	shipEventProducer, err := producer.NewShipEventProducer(*cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to initialise ship event producer: %s", err.Error()))
	}
	service := shipsrv.New(repo, shipEventProducer)

	consumer, err := consumer.NewShipDataConsumer(*cfg, service, searchService, metricsClient)
	if err != nil {
		panic(fmt.Sprintf("failed to initialise Kafka consumer: %s", err.Error()))
	}
	defer consumer.Shutdown()
	go func() {
		if err := consumer.Read(ctx); err != nil && !errors.Is(err, context.Canceled) {
			clog.Errorf("kafka consumer stopped due to error: %s", err.Error())
		}
	}()

	server, err := graphql.New(*cfg, service, searchService)
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
