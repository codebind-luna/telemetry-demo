package mongorepository

import (
	"context"
	"errors"

	"github.com/codebind-luna/telemetry-demo/internal/models"
	"github.com/codebind-luna/telemetry-demo/internal/repository/mongorepository/mongomodels"
)

// Save implements interfaces.Repository.
func (m *mongorepo) Save(ctx context.Context, exp *models.Expression) (string, error) {
	coll := m.getcalculationsCollection()
	calculation := mongomodels.NewExpression(*exp)

	// Inserts a document describing an expression into the collection
	_, insertErr := coll.InsertOne(ctx, calculation)
	if insertErr != nil {
		m.logger.WithContext(ctx).WithFields("expression", exp).WithError(errors.New(insertErr.Error())).Error("insert failed")
	}

	m.logger.WithContext(ctx).WithFields("expression", exp).Infof("created an entry in mongo")

	return exp.ID, nil
}
