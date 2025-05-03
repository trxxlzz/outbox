package redis

import (
	"context"
	"log"
	"outbox/internal/infra/redis"
	"outbox/internal/mapper"
	"outbox/internal/service/transaction"
)

func StartRedisListener(ctx context.Context, consumer *redis.RedisConsumer, service *transaction.TransactionApi) error {
	if err := consumer.InitGroup(ctx); err != nil {
		log.Printf("failed to init Redis group: %v", err)
	}

	for {
		messages, err := consumer.ReadMessages(ctx)
		if err != nil {
			log.Printf("failed to read messages: %v", err)
			continue
		}

		for _, message := range messages {
			tx := mapper.MapToTransaction(message.Values)
			if err := service.ProcessTransaction(ctx, tx); err != nil {
				log.Printf("failed to process transaction: %v", err)
				continue
			}

			if err := consumer.Ack(ctx, message.ID); err != nil {
				log.Printf("failed to ack message: %v", err)
				continue
			}
		}
	}
}
