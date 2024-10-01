package amqppublisher

import (
	amqp "github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/codebind-luna/telemetry-demo/pkg/publisher/pubsubconstants"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

func NewConfigWithDeduplication(amqpURI string) amqp.Config {
	arguments := make(amqp091.Table)
	arguments[pubsubconstants.DeduplicationQueuePropertyKey] = pubsubconstants.DeduplicationQueuePropertyValue
	config := NewPublisherConfig(amqpURI)
	config.Exchange.Arguments = arguments
	config.Queue.Arguments = arguments
	return config
}

func NewPublisherConfig(amqpURI string) amqp.Config {
	return amqp.Config{
		Connection: amqp.ConnectionConfig{
			AmqpURI: amqpURI,
		},

		Marshaler: amqp.DefaultMarshaler{},

		Exchange: amqp.ExchangeConfig{
			GenerateName: func(topic string) string {
				return topic
			},
			Type: pubsubconstants.FanOutExchangeType,
		},
		Queue: amqp.QueueConfig{
			GenerateName: func(topic string) string {
				return topic
			},
		},
		QueueBind: amqp.QueueBindConfig{
			GenerateRoutingKey: func(topic string) string {
				return topic
			},
		},
		Publish: amqp.PublishConfig{
			GenerateRoutingKey: func(topic string) string {
				return ""
			},
		},
		TopologyBuilder: &amqp.DefaultTopologyBuilder{},
	}
}
