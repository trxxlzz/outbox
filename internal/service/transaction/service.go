package transaction

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"outbox/internal/domain"
	"outbox/internal/repository"
	"time"
)

type TransactionApi struct {
	transactionRepository repository.TransactionRepository
	Producer              domain.EventProducer
	db                    *pgxpool.Pool
}

func NewTransactionService(repo repository.TransactionRepository, producer domain.EventProducer, db *pgxpool.Pool) *TransactionApi {
	return &TransactionApi{
		transactionRepository: repo,
		Producer:              producer,
		db:                    db,
	}
}

func (s *TransactionApi) StartOutboxWorker(ctx context.Context, handlePeriod time.Duration) {
	ticker := time.NewTicker(handlePeriod)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping outbox worker...")
				return
			case <-ticker.C:
				outbox, err := s.transactionRepository.GetOutboxEvent(ctx)
				if err != nil {
					log.Printf("Failed to get outbox event: %v", err)
					continue
				}

				// Если нет новых событий → пропускаем итерацию
				if outbox.ID == 0 {
					continue
				}

				// Отправляем payload в Kafka (или другую систему)
				err = s.Producer.Publish(ctx, fmt.Sprintf("%d", outbox.ID), []byte(outbox.Payload))
				if err != nil {
					log.Printf("Failed to publish event to Kafka (id=%d): %v", outbox.ID, err)
					continue // не помечаем как done, попробуем позже
				}

				// Помечаем событие как processed
				err = s.transactionRepository.MarkOutboxProcessed(ctx, outbox.ID)
				if err != nil {
					log.Printf("Failed to mark outbox event as processed (id=%d): %v", outbox.ID, err)
					continue
				}

				log.Printf("Successfully processed outbox event id=%d", outbox.ID)
			}
		}
	}()
}
