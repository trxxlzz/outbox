package transaction

import (
	"context"
	"log"
	"outbox/internal/model"
)

func (s *TransactionApi) ProcessTransaction(ctx context.Context, transaction model.Transaction) (err error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	log.Printf("Received transaction: %+v", transaction)

	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	err = s.transactionRepository.SaveTransaction(ctx, tx, transaction)
	if err != nil {
		return err // rollback сработает в defer
	}

	return nil
}
