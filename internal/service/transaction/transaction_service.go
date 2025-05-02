package transaction

import (
	"context"
	"outbox/internal/model"
)

func (s *TransactionService) ProcessTransaction(ctx context.Context, tx model.Transaction) error {
	if err := s.transactionRepository.Save(ctx, tx); err != nil {
		return err
	}

	if err := s.Producer.Publish(ctx, tx); err != nil {
		return err
	}

	return nil
}
