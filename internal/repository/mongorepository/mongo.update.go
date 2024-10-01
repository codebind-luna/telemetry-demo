package mongorepository

import (
	"context"
	"errors"

	"github.com/codebind-luna/telemetry-demo/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

// Update - implements interfaces.Repository.
func (m *mongorepo) Update(ctx context.Context, id string, status models.Status, value *int, errMsg *string) error {
	collection := m.getcalculationsCollection()

	filter := bson.D{{Key: "expID", Value: id}}

	updateFields := bson.D{
		{Key: "status", Value: status},
	}

	if value != nil {
		updateFields = append(updateFields, bson.E{Key: "value", Value: *value})
	}

	if errMsg != nil {
		updateFields = append(updateFields, bson.E{Key: "errorMsg", Value: *errMsg})
	}

	update := bson.D{
		{Key: "$set", Value: updateFields},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		m.logger.WithContext(ctx).WithFields("expID", id, "status", status, "value", value).WithError(errors.New(err.Error())).Error("update failed")
	}

	m.logger.WithContext(ctx).WithFields("expID", id, "status", status, "value", value).Infof("updated entry in mongo")

	return nil
}

func (m *mongorepo) MarkPending(ctx context.Context, id string) error {
	return m.Update(ctx, id, models.Pending, nil, nil)
}

func (m *mongorepo) MarkCompleted(ctx context.Context, id string, val *int) error {
	return m.Update(ctx, id, models.Completed, val, nil)
}

func (m *mongorepo) MarkFailed(ctx context.Context, id string, failureMessage *string) error {
	return m.Update(ctx, id, models.Failed, nil, failureMessage)
}

func (m *mongorepo) MarkInProgress(ctx context.Context, id string) error {
	return m.Update(ctx, id, models.InProgress, nil, nil)
}
