package transaction

import (
	"outbox/internal/domain"
	"outbox/repository"
)

type TransactionService struct {
	transactionRepository repository.TransactionRepository
	Producer              domain.EventProducer
}

func NewTransactionService(repo repository.TransactionRepository, producer domain.EventProducer) *TransactionService {
	return &TransactionService{transactionRepository: repo, Producer: producer}
}
