package publisher

import (
	"github.com/codebind-luna/telemetry-demo/internal/interfaces"
	"github.com/codebind-luna/telemetry-demo/pkg/logger"
	amqppublisher "github.com/codebind-luna/telemetry-demo/pkg/publisher/amqp"
)

// New - create a new instance of a publisher
func New(
	option interfaces.MessageBusType,
	logger *logger.Logger,
) (interfaces.Publisher, error) {
	var (
		pub interfaces.Publisher
		err error
	)
	switch option {
	case interfaces.AMQP:
		pub, err = amqppublisher.NewAMQPPublisher(logger)
	default:
		return nil, interfaces.ErrInvalidMessageBusType
	}
	return pub, err
}
