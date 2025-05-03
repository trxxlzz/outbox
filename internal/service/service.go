package service

import (
	"context"
	"outbox/internal/model"
)

type TransactionService interface {
	ProcessTransaction(ctx context.Context, tx model.Transaction) error
}
