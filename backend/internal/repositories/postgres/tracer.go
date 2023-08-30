package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "github.com/mikeewhite/ship-locator/postgres"

type Tracer struct {
	tracer trace.Tracer
}

func newTracer() *Tracer {
	return &Tracer{
		tracer: otel.Tracer(tracerName),
	}
}

func (t *Tracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx
	}
	var opts []trace.SpanStartOption
	opts = append(opts, trace.WithAttributes(semconv.DBStatementKey.String(data.SQL)))
	opts = append(opts, trace.WithAttributes(makeParamsAttribute(data.Args)))

	dbOp := "UNKNOWN"
	parts := strings.Fields(data.SQL)
	if len(parts) > 0 {
		dbOp = strings.ToUpper(parts[0])
	}

	dbName := conn.Config().Config.Database

	// see https://opentelemetry.io/docs/specs/otel/trace/semantic_conventions/database/
	spanName := fmt.Sprintf("%s %s", dbOp, dbName)

	ctx, _ = t.tracer.Start(ctx, spanName, opts...)

	return ctx
}

func (t *Tracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	span := trace.SpanFromContext(ctx)
	if data.Err != nil {
		span.RecordError(data.Err)
		span.SetStatus(codes.Error, data.Err.Error())
	}
	span.End()
}

func makeParamsAttribute(args []any) attribute.KeyValue {
	ss := make([]string, len(args))
	for i := range args {
		ss[i] = fmt.Sprintf("%+v", args[i])
	}
	return attribute.Key("pgx.query.parameters").StringSlice(ss)
}
