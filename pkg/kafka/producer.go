package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

var writer *kafka.Writer

// InitProducer initializes a global Kafka writer
func InitProducer(broker string, topic string) {
	writer = &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

// CloseProducer closes the global Kafka writer
func CloseProducer() error {
	if writer != nil {
		return writer.Close()
	}
	return nil
}

// ProduceMessage sends a message using the global writer
func ProduceMessage(key string, value []byte) error {
	if writer == nil {
		return nil // or return error "producer not initialized"
	}
	return writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: value,
		},
	)
}
