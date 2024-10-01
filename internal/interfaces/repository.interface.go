package interfaces

import (
	"context"
	"fmt"

	"github.com/codebind-luna/telemetry-demo/internal/models"
)

type RepositoryType string

func (r RepositoryType) String() string {
	return string(r)
}

const (
	MongoRepository      RepositoryType = "mongo"
	PostgreSQLRepository RepositoryType = "postgreSQL"
)

var (
	validRepositories = []RepositoryType{
		MongoRepository,
	}
	ErrInvalidRepositoryType = fmt.Errorf("invalid repository type")
	repositoryMap            = map[string]RepositoryType{
		MongoRepository.String(): MongoRepository,
	}
)

func isValidRepository(repo RepositoryType) bool {
	for _, valid := range validRepositories {
		if valid == repo {
			return true
		}
	}
	return false
}

func ParseRepository(repo string) (RepositoryType, error) {
	r, ok := repositoryMap[repo]
	if !ok {
		return "", ErrInvalidRepositoryType
	}
	if !isValidRepository(r) {
		return "", ErrInvalidRepositoryType
	}
	return r, nil
}

type Repository interface {
	Healthy(ctx context.Context) bool
	Save(ctx context.Context, exp *models.Expression) (string, error)
	Get(ctx context.Context, id string) (*models.Expression, error)
	Update(ctx context.Context, id string, status models.Status, value *int, errMsg *string) error
	MarkPending(ctx context.Context, id string) error
	MarkInProgress(ctx context.Context, id string) error
	MarkCompleted(ctx context.Context, id string, val *int) error
	MarkFailed(ctx context.Context, id string, failureMessage *string) error
}
