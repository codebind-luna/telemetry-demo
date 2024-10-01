package logger

import (
	"context"
	"os"

	"github.com/codebind-luna/telemetry-demo/pkg/constants"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.SugaredLogger
}

func NewLogger() *Logger {
	// Create a pretty JSON encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.LineEnding = zapcore.DefaultLineEnding
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// Create a core that uses the pretty encoder
	core := zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)

	// Create the logger
	logger := zap.New(core)
	defer logger.Sync() // Flushes buffer, if any

	// zapconf := zapdriver.NewProductionConfig()

	// l, err := zapconf.Build(
	// 	zapdriver.WrapCore(
	// 		zapdriver.ReportAllErrors(true),
	// 		zapdriver.ServiceName(ServiceName),
	// 	),
	// 	zap.AddCallerSkip(1),
	// )
	// if err != nil {
	// 	log.Fatalf("failed to init logger: %v", err)
	// }

	return &Logger{logger: logger.Sugar()}
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	span := trace.SpanFromContext(ctx).SpanContext()

	return &Logger{
		logger: l.logger.With(
			constants.TraceKey, span.TraceID().String(),
			constants.SpanKey, span.SpanID().String(),
		),
	}
}

func (l *Logger) WithFields(args ...interface{}) *Logger {
	return &Logger{
		logger: l.logger.With(args...),
	}
}

func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *Logger) WithError(err error) *Logger {
	return &Logger{logger: l.logger.With(constants.ErrorKey, err.Error())}
}
