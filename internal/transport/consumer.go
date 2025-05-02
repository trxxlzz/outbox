package transport

import (
	"context"
	"log"
	"outbox/internal/model"
	"outbox/internal/service/transaction"
)

type RedisConsumer interface {
	Listen(ctx context.Context, handler func(tx model.Transaction)) error
}

// StartConsumer запускает Redis consumer и передаёт обработку в бизнес-слой
func StartConsumer(ctx context.Context, consumer RedisConsumer, service *transaction.TransactionService) error {
	return consumer.Listen(ctx, func(tx model.Transaction) {
		if err := service.ProcessTransaction(ctx, tx); err != nil {
			log.Printf("failed to process transaction: %v", err)
		}
	})
}
