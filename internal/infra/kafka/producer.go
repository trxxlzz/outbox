package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"outbox/internal/domain"
	"outbox/internal/model"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func (p *KafkaProducer) Publish(ctx context.Context, tx model.Transaction) error {
	payload, err := json.Marshal(tx)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(tx.ID),
		Value: payload,
	}

	return p.writer.WriteMessages(ctx, msg)
}

func NewKafkaProducer(writer *kafka.Writer) domain.EventProducer {
	return &KafkaProducer{writer: writer}
}
