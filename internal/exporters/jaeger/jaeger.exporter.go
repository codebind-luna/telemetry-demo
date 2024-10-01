package jaeger

import (
	"context"

	"github.com/codebind-luna/telemetry-demo/internal/interfaces"
	"github.com/codebind-luna/telemetry-demo/pkg/logger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
)

var _ interfaces.TracerExporter = (*jaegerExporter)(nil)

type jaegerExporter struct {
	exp sdkTrace.SpanExporter
}

// GetExporter - get the stdout exporter to retrieve the collected spans
func (a *jaegerExporter) GetExporter() sdkTrace.SpanExporter {
	return a.exp
}

// New - create a new instance of the jaeger exporter
func New(
	logger *logger.Logger,
) (*jaegerExporter, error) {
	headers := map[string]string{
		"content-type": "application/json",
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint("localhost:4318"),
			otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		return nil, err
	}
	return &jaegerExporter{
		exp: exporter,
	}, nil
}
