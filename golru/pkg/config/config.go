package config

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type Tracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanOption) (context.Context, trace.Span)
}

type PostgresDBParams struct {
	DbName   string
	Host     string
	User     string
	Password string
}
