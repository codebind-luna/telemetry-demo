package middleware

import (
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// CustomResponseWriter wraps the standard ResponseWriter
type CustomResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

// WriteHeader captures the status code
func (rw *CustomResponseWriter) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// NewCustomResponseWriter creates a new CustomResponseWriter
func NewCustomResponseWriter(w http.ResponseWriter) *CustomResponseWriter {
	return &CustomResponseWriter{
		ResponseWriter: w,
		StatusCode:     http.StatusOK, // Default to 200
	}
}

// OTelMiddleware is a middleware for tracing HTTP requests.
func OTelMiddleware(service string, next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tracer := otel.Tracer(service)
		spanCtx, span := tracer.Start(ctx, "http_request")
		defer span.End()

		// Add span attributes
		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.url", r.URL.String()),
		)
		// Create a custom response writer to capture the status code
		customWriter := NewCustomResponseWriter(w)

		// Call the next handler in the chain
		next(customWriter, r.WithContext(spanCtx))

		// Now you can access the status code after the handler finishes
		statusCode := customWriter.StatusCode
		span.SetAttributes(
			attribute.String("http.status_code", strconv.Itoa(statusCode)),
		)
	}
}
