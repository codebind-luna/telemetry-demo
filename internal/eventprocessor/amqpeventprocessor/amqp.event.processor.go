package amqpeventprocessor

import (
	"fmt"

	amqp "github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/codebind-luna/telemetry-demo/internal/interfaces"
	"github.com/sirupsen/logrus"

	logrusadapter "logur.dev/adapter/logrus"
	watermillintegration "logur.dev/integration/watermill"
)

const (
	FanOutExchangeType = "fanout"
)

var _ interfaces.EventProcessor = (*amqpEventProcessor)(nil)

type amqpEventProcessor struct {
	expressionSubscriber message.Subscriber
}

func (a *amqpEventProcessor) GetExpressionSubscriber() message.Subscriber {
	return a.expressionSubscriber
}

// New - create a new instance of the amqp event processor
func New() (*amqpEventProcessor, error) {
	logger := &logrus.Logger{}
	// Create a logger instance
	log := watermillintegration.New(logrusadapter.New(logger))

	// Setup expression subscriber
	logger.Info("initializing expression subscriber")
	config := newConfigWithDeduplication(
		"amqp://guest:guest@localhost:5672/",
	)
	expressionSubscriber, err := amqp.NewSubscriber(config, log)
	if err != nil {
		return nil, fmt.Errorf("could not initialize ping pong subscriber: %+v", err)
	}
	return &amqpEventProcessor{
		expressionSubscriber: expressionSubscriber,
	}, nil
}
