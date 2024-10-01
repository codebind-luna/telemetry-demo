package tracermetrics

import (
	"context"

	"github.com/codebind-luna/telemetry-demo/internal/exporters"
	"github.com/codebind-luna/telemetry-demo/internal/interfaces"
	"github.com/codebind-luna/telemetry-demo/pkg/logger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

type TracerCollector interface {
	Start()
	Stop(ctx context.Context)
}

var _ TracerCollector = (*impl)(nil)

type impl struct {
	serviceName    string
	appVersion     string
	log            *logger.Logger
	tracerProvider *sdktrace.TracerProvider
	collector      interfaces.TracerExporter

	enableTracing bool
}

func initTracer(resources *resource.Resource, expSpan sdktrace.SpanExporter) *sdktrace.TracerProvider {
	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(expSpan),
		sdktrace.WithResource(resources),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func New(
	serviceName string,
	appVersion string,
	logger *logger.Logger,
	enableTracing bool,
	tracerExporterType interfaces.TracerExporterType,
) TracerCollector {
	tracerExp, _ := exporters.New(tracerExporterType, logger)
	return &impl{
		serviceName:   serviceName,
		appVersion:    appVersion,
		log:           logger,
		enableTracing: enableTracing,
		collector:     tracerExp,
	}
}

func (i *impl) Start() {
	attr := resource.WithAttributes(
		semconv.DeploymentEnvironmentKey.String("development"),
		semconv.ServiceNameKey.String(i.serviceName),
		semconv.ServiceVersionKey.String(i.appVersion),
	)

	resources, _ := resource.New(context.Background(),
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithHost(),
		attr)

	if i.enableTracing {
		i.log.Info("initializing tracer provider")
		tp := initTracer(resources, i.collector.GetExporter())

		i.tracerProvider = tp
	}
}

func (i *impl) Stop(ctx context.Context) {

	if i.enableTracing {
		if err := i.tracerProvider.Shutdown(ctx); err != nil {
			i.log.Errorf("Error shutting down tracer provider: %v", err)
		} else {
			i.log.Info("tracing provider shutdown successfully")
		}
	}
}
