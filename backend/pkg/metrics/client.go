package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/mikeewhite/ship-locator/backend/pkg/clog"
	"github.com/mikeewhite/ship-locator/backend/pkg/config"
)

type Client struct {
	httpServer *http.Server

	dbQueryTimeHistogram      *prometheus.HistogramVec
	kafkaConsumeTimeHistogram *prometheus.HistogramVec
}

func New(cfg config.Config) *Client {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:              cfg.PrometheusServerAddress,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      5 * time.Minute,
	}

	client := &Client{
		httpServer: server,
	}

	client.dbQueryTimeHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "db_query_time_secs",
		Help: "Duration of DB queries",
	}, []string{"query_name"})

	client.kafkaConsumeTimeHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "kafka_consume_time_secs",
		Help: "Kafka consume time",
	}, []string{"topic"})

	return client
}

func (c *Client) Serve(ctx context.Context) error {
	clog.Infof("Starting Prometheus client at %s", c.httpServer.Addr)
	go func() {
		<-ctx.Done()
		clog.Info("Stopping Prometheus client")
		c.Shutdown()
	}()
	return c.httpServer.ListenAndServe()
}

func (c *Client) Shutdown() {
	// allow 5 seconds for the server to shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := c.httpServer.Shutdown(ctx); err != nil {
		clog.Errorf("failed to shutdown prometheus client cleanly: %s", err.Error())
	}
}
