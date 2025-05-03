package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"outbox/internal/model"
	modelRepo "outbox/internal/repository/model"
)

type TransactionRepository interface {
	SaveTransaction(ctx context.Context, tx pgx.Tx, transaction model.Transaction) error
	SaveOutboxEvent(ctx context.Context, tx pgx.Tx, eventType string, payload interface{}) error
	GetOutboxEvent(ctx context.Context) (modelRepo.Outbox, error)
	MarkOutboxProcessed(ctx context.Context, id int) error
}
