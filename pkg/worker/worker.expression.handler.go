package worker

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"

	"go.opentelemetry.io/otel/trace"
)

type contextKey string

const requestIDKey contextKey = "requestID"

func (w *workerImpl) handleExpression(msg *message.Message) error {
	ctx := msg.Context()

	// Retrieve the value using the same key type
	requestID := ctx.Value(requestIDKey)

	if requestID != nil {
		fmt.Printf("Processing request with ID: %s\n", requestID)
	} else {
		fmt.Println("No request ID found in context")
	}

	span := trace.SpanFromContext(ctx).SpanContext()
	// w.svc.ProcessExpression(msg.Context())
	w.log.WithContext(ctx).WithFields("messageId", msg.UUID, "trace_id", span.TraceID().String()).Info("successfully processed")
	return nil
}
