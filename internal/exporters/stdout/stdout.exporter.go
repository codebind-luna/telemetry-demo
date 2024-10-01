package stdout

import (
	"github.com/codebind-luna/telemetry-demo/internal/interfaces"
	"github.com/codebind-luna/telemetry-demo/pkg/logger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
)

var _ interfaces.TracerExporter = (*stdoutExporter)(nil)

type stdoutExporter struct {
	exp sdkTrace.SpanExporter
}

// GetExporter - get the stdout exporter to retrieve the collected spans
func (a *stdoutExporter) GetExporter() sdkTrace.SpanExporter {
	return a.exp
}

// New - create a new instance of the amqp event processor
func New(
	logger *logger.Logger,
) (*stdoutExporter, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}
	return &stdoutExporter{
		exp: exporter,
	}, nil
}
