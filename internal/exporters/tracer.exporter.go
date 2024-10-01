package exporters

import (
	"errors"

	"github.com/codebind-luna/telemetry-demo/internal/exporters/jaeger"
	"github.com/codebind-luna/telemetry-demo/internal/exporters/stdout"
	"github.com/codebind-luna/telemetry-demo/internal/interfaces"
	"github.com/codebind-luna/telemetry-demo/pkg/logger"
)

// New - retrieve a tracer exporter
func New(
	tracerExporterType interfaces.TracerExporterType,
	logger *logger.Logger,
) (interfaces.TracerExporter, error) {
	switch tracerExporterType {
	case interfaces.STDOUT:
		return stdout.New(logger)
	case interfaces.JAEGER:
		return jaeger.New(logger)
	default:
		return nil, errors.New("invalid tracer type provided")
	}
}
