package handlers

import (
	"github.com/codebind-luna/telemetry-demo/internal/interfaces"
	"github.com/codebind-luna/telemetry-demo/pkg/logger"
	
)

type expressionsHandler struct {
	logger *logger.Logger
	svc    interfaces.ExpressionsService
}

// NewExpressionsHandler creates a new instance of expressionsHandler
func NewExpressionsHandler(
	logger *logger.Logger,
	svc interfaces.ExpressionsService) *expressionsHandler {
	return &expressionsHandler{
		svc:    svc,
		logger: logger,
	}
}
