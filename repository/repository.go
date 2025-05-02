package repository

import (
	"context"
	"outbox/internal/model"
)

type TransactionRepository interface {
	Save(ctx context.Context, transaction model.Transaction) error
}
