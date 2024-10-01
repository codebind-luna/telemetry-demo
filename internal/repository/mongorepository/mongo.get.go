package mongorepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/codebind-luna/telemetry-demo/internal/models"
	"github.com/codebind-luna/telemetry-demo/internal/repository/mongorepository/mongomodels"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Get implements interfaces.Repository.
func (m *mongorepo) Get(ctx context.Context, id string) (*models.Expression, error) {
	collection := m.getcalculationsCollection()
	//  Creates a filter to match a document that has the specified
	//  "expID" value
	filter := bson.D{{Key: "expID", Value: id}}

	var expression mongomodels.Expression

	getErr := collection.FindOne(context.TODO(), filter).Decode(&expression)
	if getErr != nil {
		if getErr == mongo.ErrNoDocuments {
			m.logger.WithContext(ctx).WithFields("expression", id).WithError(errors.New("no matching record found")).Error("fetch failed")
			return nil, fmt.Errorf("no matching record found for expression %s", id)
		} else {
			m.logger.WithContext(ctx).WithFields("expression", id).WithError(errors.New("failed to retrieve record")).Error("fetch failed")
			return nil, fmt.Errorf("failed to retrieve record for expression %s: %+v", id, getErr)
		}
	}

	m.logger.WithContext(ctx).WithFields("expression", id).Info("retrieved the record for expression")

	return expression.To(), nil
}
