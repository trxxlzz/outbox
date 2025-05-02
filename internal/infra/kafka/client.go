package kafka

import (
	"github.com/segmentio/kafka-go"
	"time"
)

func NewKafkaWriter(brokers []string, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{}, // распределяет нагрузку по партициям
		RequiredAcks: kafka.RequireAll,    // ждём подтверждения от всех брокеров
		Async:        false,               // включить true, если хочешь асинхронную отправку
		BatchTimeout: 10 * time.Millisecond,
	}
}
