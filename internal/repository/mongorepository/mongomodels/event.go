package mongomodels

import (
	"time"

	"github.com/codebind-luna/telemetry-demo/internal/models"
)

type Event struct {
	ID        string        `bson:"_id,omitempty"`
	EventId   string        `bson:"eventId,omitempty"`
	CreatedAt int64         `bson:"createdAt"`
	UpdatedAt int64         `bson:"updatedAt"`
	lock      *bool         `bson:"lock,omitempty"`
	Status    models.Status `bson:"status,omitempty"`
}

func NewEvent(eventId string) Event {
	return Event{
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		EventId:   eventId,
		Status:    models.UnProcessed,
	}
}
