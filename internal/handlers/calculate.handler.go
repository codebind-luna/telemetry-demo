package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type postCalculateParams struct {
	Expression string `json:"expression"`
}

type calculateResponse struct {
	ID string `json:"id"`
}

// Calculate handles requests to the "/calculate" route
func (s *expressionsHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Only allow POST requests
	if r.Method != http.MethodPost {
		s.logger.WithContext(ctx).WithFields("request.method", r.Method, "operation", "calculate expression").WithError(errors.New("method not allowed")).Error("operation failed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.WithContext(ctx).WithFields("request.method", r.Method, "operation", "calculate expression").WithError(errors.New("unable to read request body")).Error("operation failed")
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Unmarshal JSON data into User struct
	var expression postCalculateParams
	if err := json.Unmarshal(body, &expression); err != nil {
		s.logger.WithContext(ctx).WithFields("request.method", r.Method, "operation", "calculate expression").WithError(errors.New("invalid JSON format")).Error("operation failed")
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	expID, createErr := s.svc.Calculate(ctx, expression.Expression)
	if createErr != nil {
		s.logger.WithContext(ctx).WithFields("request.method", r.Method, "operation", "calculate expression").WithError(errors.New("invalid JSON format")).WithFields("expression", expression.Expression)
		http.Error(w, "Unable to create expression", http.StatusBadRequest)
		return
	}

	// Set Content-Type to application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Set HTTP status to 200 OK

	s.logger.WithContext(ctx).WithFields("expression", expression.Expression, "request.method", r.Method, "operation", "calculate expression").Info("successfully created expression for future processing")

	json.NewEncoder(w).Encode(calculateResponse{ID: expID})
}
