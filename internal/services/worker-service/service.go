package workerService

import (
	"context"

	"github.com/codebind-luna/telemetry-demo/internal/backend"
	"github.com/codebind-luna/telemetry-demo/internal/interfaces"
	"github.com/codebind-luna/telemetry-demo/pkg/logger"
)

var _ interfaces.WorkerService = (*impl)(nil)

type impl struct {
	logger             *logger.Logger
	expressionsbackend backend.Expressions
}

func NewService(logger *logger.Logger, expressionsbackend backend.Expressions) *impl {
	return &impl{
		logger:             logger,
		expressionsbackend: expressionsbackend,
	}
}

func (s *impl) ProcessExpression(ctx context.Context) (*int, error) {
	return s.expressionsbackend.Process(ctx, "", "")
}
