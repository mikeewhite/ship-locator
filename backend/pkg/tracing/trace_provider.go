package tracing

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/mikeewhite/ship-locator/backend/pkg/clog"
	"github.com/mikeewhite/ship-locator/backend/pkg/config"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type TraceProvider struct {
	traceProvider *tracesdk.TracerProvider
}

func NewTraceProvider(ctx context.Context, cfg config.Config, appName string) (*TraceProvider, error) {
	exporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint(cfg.TracingCollectorAddress),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start trace exporter: %w", err)
	}
	traceProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(appName),
		)),
	)

	// register the TraceProvider as the global instance
	otel.SetTracerProvider(traceProvider)

	clog.Infof("Exporting traces for '%s' to collector at %s", appName, cfg.TracingCollectorAddress)

	return &TraceProvider{traceProvider: traceProvider}, nil
}

func (tp *TraceProvider) Shutdown(ctx context.Context) {
	// allow 5 seconds for the trace provider to flush
	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	if err := tp.traceProvider.Shutdown(cancelCtx); err != nil {
		clog.Errorf("failed to cleanly shutdown trace provider: %s", err.Error())
	}
}
