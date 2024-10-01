package amqpeventprocessor

import (
	amqp "github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/codebind-luna/telemetry-demo/pkg/publisher/pubsubconstants"
	amqp091 "github.com/rabbitmq/amqp091-go"
)

func newConfigWithDeduplication(amqpURI string) amqp.Config {
	arguments := make(amqp091.Table)
	arguments[pubsubconstants.DeduplicationQueuePropertyKey] = pubsubconstants.DeduplicationQueuePropertyValue
	config := setupSubscriber(amqpURI, pubsubconstants.FanOutExchangeType)
	config.Exchange.Arguments = arguments
	config.Queue.Arguments = arguments
	return config
}

// setupSubscriber - configuration for any subscriber
func setupSubscriber(
	amqpURI string,
	exchangeType string,
) amqp.Config {
	// Generate the config
	amqpConfig := amqp.Config{
		Connection: amqp.ConnectionConfig{
			AmqpURI: amqpURI,
		},
		Marshaler: amqp.DefaultMarshaler{},
		Exchange: amqp.ExchangeConfig{
			GenerateName: func(topic string) string {
				return topic
			},
			Type: exchangeType,
		},
		Queue: amqp.QueueConfig{
			GenerateName: func(topic string) string {
				return topic
			},
		},
		QueueBind: amqp.QueueBindConfig{
			GenerateRoutingKey: func(topic string) string {
				return ""
			},
		},
		TopologyBuilder: &amqp.DefaultTopologyBuilder{},
	}
	return amqpConfig
}
