package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	// "go.opentelemetry.io/otel/attribute"
	// "go.opentelemetry.io/otel"
	// "go.opentelemetry.io/otel/attribute"
)

// Get handles requests to the "/get" route
func (s *expressionsHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// ctx, span := otel.Tracer("expressions-service").Start(ctx, "get-expression-handler")
	// defer span.End()

	// Extract the user ID from the URL path
	id := strings.TrimPrefix(r.URL.Path, "/expressions/")

	s.logger.WithContext(ctx).WithFields("expressionID", id).Info("received request to retrieve expression")

	exp, getErr := s.svc.Get(ctx, id)

	if getErr != nil {
		http.Error(w, getErr.Error(), http.StatusBadRequest)
		return
	}

	// Set Content-Type to application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Set HTTP status to 200 OK

	// Add attributes to the external API call span
	// span.SetAttributes(
	// 	attribute.String("expression.id", exp.ID),
	// 	attribute.String("expression.exp", exp.Exp),
	// 	attribute.String("expression.status", exp.Status.String()),
	// )

	// Encode user to JSON and write to response
	json.NewEncoder(w).Encode(exp)
}
