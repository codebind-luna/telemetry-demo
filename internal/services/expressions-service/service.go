package expressionsService

import (
	"context"

	"github.com/codebind-luna/telemetry-demo/internal/interfaces"
	"github.com/codebind-luna/telemetry-demo/internal/models"
	"github.com/codebind-luna/telemetry-demo/pkg/logger"
	"github.com/codebind-luna/telemetry-demo/pkg/publisher/pubsubmodels"
	"go.opentelemetry.io/otel"
)

var _ interfaces.ExpressionsService = (*impl)(nil)

type impl struct {
	logger    *logger.Logger
	repo      interfaces.Repository
	publisher interfaces.Publisher
}

func NewService(logger *logger.Logger, repo interfaces.Repository, publisher interfaces.Publisher) *impl {
	return &impl{
		logger:    logger,
		repo:      repo,
		publisher: publisher,
	}
}

// Calculate - takes a string expression
// Saves entry to datastore
// and publishes the event for future calculation
func (s *impl) Calculate(ctx context.Context, exp string) (string, error) {
	s.logger.WithContext(ctx).WithFields("expression", exp).Info("received request to process expression")
	expression := models.NewEpression(exp)
	str, err := s.repo.Save(ctx, expression)
	if err != nil {
		return "", err
	}

	publishErr := s.publisher.PublishExpression(ctx, pubsubmodels.Event{EventId: expression.ID, Expression: expression.Exp})
	if publishErr != nil {
		return "", publishErr
	}

	return str, nil
}

// Get - takes an ID of an expression and returns the expression data
// from backend datastore
func (s *impl) Get(ctx context.Context, id string) (*models.Expression, error) {
	// Retrieve the tracer from the context
	tracer := otel.Tracer("expressions-service")

	// Create a new span
	derivedCtx, span := tracer.Start(ctx, "get-expressions-service")
	defer span.End()

	s.logger.WithContext(derivedCtx).Infof("processing request to retrieve expression for id %s", id)
	exp, getErr := s.repo.Get(ctx, id)
	if getErr != nil {
		return nil, getErr
	}

	return exp, nil
}
