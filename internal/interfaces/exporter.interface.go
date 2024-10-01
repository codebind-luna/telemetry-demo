package interfaces

import (
	"errors"

	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
)

type TracerExporter interface {
	GetExporter() sdkTrace.SpanExporter
}

type TracerExporterType string

const (
	STDOUT  TracerExporterType = "stdout"
	JAEGER  TracerExporterType = "jaeger"
	ASPECTO TracerExporterType = "aspecto"
)

func (s TracerExporterType) String() string {
	return string(s)
}

var (
	// options used to validate valid TracerExporter types
	validOptions = map[string]TracerExporterType{
		STDOUT.String():  STDOUT,
		JAEGER.String():  JAEGER,
		ASPECTO.String(): ASPECTO,
	}

	// ErrInvalidTracerExporterType - error when invalid TracerExporter type is provided
	ErrInvalidTracerExporterType = errors.New("invalid TracerExporter type provided")
)

// ParseTracerExporterType - parse a string into a TracerExporter type
func ParseTracerExporterType(busType string) (TracerExporterType, error) {
	bus, ok := validOptions[busType]
	if !ok {
		return "", ErrInvalidTracerExporterType
	}
	return bus, nil
}
