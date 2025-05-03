package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
)

func (p *KafkaProducer) Publish(ctx context.Context, key string, payload []byte) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: payload,
	}

	return p.writer.WriteMessages(ctx, msg)
}
