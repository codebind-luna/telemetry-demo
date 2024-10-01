package mongomodels

import (
	"time"

	"github.com/codebind-luna/telemetry-demo/internal/models"
)

type Expression struct {
	ID        string        `bson:"_id,omitempty"`
	ExpID     string        `bson:"expID,omitempty"`
	Exp       string        `bson:"exp,omitempty"`
	Value     *int          `bson:"value,omitempty"`
	Status    models.Status `bson:"status,omitempty"`
	CreatedAt int64         `bson:"createdAt"`
	UpdatedAt int64         `bson:"updatedAt"`
	ErrorMsg  *string       `bson:"errorMsg,omitempty"`
}

func NewExpression(exp models.Expression) Expression {
	return Expression{
		Exp:       exp.Exp,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		Status:    exp.Status,
		ExpID:     exp.ID,
	}
}

func (e *Expression) To() *models.Expression {
	exp := &models.Expression{}
	exp.Exp = e.Exp
	exp.ID = e.ExpID
	exp.Status = e.Status
	exp.Value = e.Value
	exp.ErrMsg = e.ErrorMsg
	return exp
}
