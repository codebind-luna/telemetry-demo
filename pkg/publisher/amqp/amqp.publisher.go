package amqppublisher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	amqp "github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/codebind-luna/telemetry-demo/internal/interfaces"
	"github.com/codebind-luna/telemetry-demo/pkg/logger"
	"github.com/codebind-luna/telemetry-demo/pkg/publisher/pubsubconstants"
	"github.com/codebind-luna/telemetry-demo/pkg/publisher/pubsubmodels"
	wotel "github.com/voi-oss/watermill-opentelemetry/pkg/opentelemetry"
)

var _ interfaces.Publisher = (*amqpPublisher)(nil)

type amqpPublisher struct {
	logger   *logger.Logger
	pulisher message.Publisher
}

func NewAMQPPublisher(logger *logger.Logger) (*amqpPublisher, error) {
	publisher, err := amqp.NewPublisher(
		NewConfigWithDeduplication("amqp://guest:guest@localhost:5672/"),
		watermill.NewStdLogger(false, false))
	if err != nil {
		return nil, err
	}

	return &amqpPublisher{pulisher: wotel.NewPublisherDecorator(publisher), logger: logger}, nil
}

func (a *amqpPublisher) PublishExpression(ctx context.Context, msg pubsubmodels.Event) error {
	metadata := map[string]string{
		pubsubconstants.DeduplicationMessageHeader: msg.EventId,
	}

	return a.publish(ctx, pubsubconstants.ExpressionTopicName, msg, metadata)
}

type contextKey string

const requestIDKey contextKey = "requestID"

func (a *amqpPublisher) publish(ctx context.Context, topic string, payload interface{}, metadata map[string]string) error {
	// Marshal the struct to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		a.logger.WithContext(ctx).WithFields("operation", "publish").WithError(fmt.Errorf("failed to marshal data: %v", err)).Error("failed")
	}

	derivedCtx := context.WithValue(ctx, requestIDKey, "67")
	msg := message.NewMessage(watermill.NewUUID(), jsonData)
	msg.SetContext(derivedCtx)

	for k, v := range metadata {
		msg.Metadata.Set(k, v)
	}

	if err := a.pulisher.Publish(topic, msg); err != nil {
		a.logger.WithContext(derivedCtx).WithFields("operation", "publish").WithError(fmt.Errorf("failed to publish message: %v", err)).Error("failed")
	}

	a.logger.WithContext(derivedCtx).WithFields("operation", "publish", "messageId", msg.UUID).Info("successful")
	return nil
}
