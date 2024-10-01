package models

import (
	"github.com/google/uuid"
)

type Status string

const (
	InProgress Status = "inprogress"
	Pending    Status = "pending"
	Completed  Status = "completed"
	Failed     Status = "failed"
)

func (s Status) String() string {
	return string(s)
}

type Expression struct {
	ID     string  `json:"id"`
	Exp    string  `json:"exp"`
	Value  *int    `json:"value,omitempty"`
	Status Status  `json:"status"`
	ErrMsg *string `json:"errMsg,omitempty"`
}

func NewEpression(exp string) *Expression {
	return &Expression{
		Exp:    exp,
		Status: Pending,
		ID:     uuid.NewString(),
	}
}
