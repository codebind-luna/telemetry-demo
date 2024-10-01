package eventprocessor

import (
	"github.com/codebind-luna/telemetry-demo/internal/eventprocessor/amqpeventprocessor"
	"github.com/codebind-luna/telemetry-demo/internal/interfaces"
)

// New - create a new instance of an event processor
func New(
	option interfaces.MessageBusType,
) (interfaces.EventProcessor, error) {

	switch option {
	case interfaces.AMQP:
		return amqpeventprocessor.New()
	default:
		return nil, interfaces.ErrInvalidMessageBusType
	}
}
