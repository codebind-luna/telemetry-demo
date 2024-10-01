package interfaces

import (
	"context"
	"errors"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/codebind-luna/telemetry-demo/pkg/publisher/pubsubmodels"
)

type EventProcessor interface {
	GetExpressionSubscriber() message.Subscriber
}

type Publisher interface {
	PublishExpression(ctx context.Context, message pubsubmodels.Event) error
}

type MessageBusType string

func (s MessageBusType) String() string {
	return string(s)
}

var (
	// options used to validate valid bus types
	validBusOptions = map[string]MessageBusType{
		AMQP.String(): AMQP,
	}
	// ErrInvalidMessageBusType - error when invalid bus type is provided
	ErrInvalidMessageBusType = errors.New("invalid message bus type provided")
)

// ParseMessageBusType - parse a string into a bus type
func ParseMessageBusType(busType string) (MessageBusType, error) {
	bus, ok := validBusOptions[busType]
	if !ok {
		return "", ErrInvalidMessageBusType
	}
	return bus, nil
}

const (
	AMQP MessageBusType = "amqp" // AMQP type (e.g RabbitMQ)
)
