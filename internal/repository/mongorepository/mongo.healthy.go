package mongorepository

import (
	"context"

	"github.com/codebind-luna/telemetry-demo/internal/repository/mongorepository/mongoutils"
)

// Healthy - check if mongodb is healthy
func (m *mongorepo) Healthy(ctx context.Context) bool {
	if err := mongoutils.Ping(m.client, ctx, m.logger); err != nil {
		return false
	}
	return true
}
