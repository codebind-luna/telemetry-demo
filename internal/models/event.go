package models

import "github.com/google/uuid"

type EventStatus string

const (
	UnProcessed Status = "unprocessed"
	Processed   Status = "processed"
	InProcess   Status = "inprocess"
)

func (s EventStatus) String() string {
	return string(s)
}

type Event struct {
	ID      string      `json:"_id,omitempty"`
	EventId string      `json:"eventId,omitempty"`
	Status  EventStatus `json:"status,omitempty"`
}

func NewEvent() *Expression {
	return &Expression{
		Status: UnProcessed,
		ID:     uuid.NewString(),
	}
}
