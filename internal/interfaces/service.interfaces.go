package interfaces

import (
	"context"

	"github.com/codebind-luna/telemetry-demo/internal/models"
)

type ExpressionsService interface {
	Calculate(context.Context, string) (string, error)
	Get(context.Context, string) (*models.Expression, error)
}

type WorkerService interface {
	ProcessExpression(ctx context.Context) (*int, error)
}
