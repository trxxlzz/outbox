package kafka

import (
	"github.com/segmentio/kafka-go"
	"outbox/internal/domain"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(writer *kafka.Writer) domain.EventProducer {
	return &KafkaProducer{writer: writer}
}
