package repository

import (
	"errors"

	"github.com/codebind-luna/telemetry-demo/internal/interfaces"
	"github.com/codebind-luna/telemetry-demo/internal/repository/mongorepository"
	"github.com/codebind-luna/telemetry-demo/pkg/logger"
)

// New - retrieve a repository
func New(
	repoType interfaces.RepositoryType,
	logger *logger.Logger,
) (interfaces.Repository, error) {
	switch repoType {
	case interfaces.MongoRepository:
		return mongorepository.New(logger)
	default:
		return nil, errors.New("invalid repository provided")
	}
}
