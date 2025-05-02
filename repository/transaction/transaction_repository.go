package transaction

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"outbox/internal/model"
	"outbox/repository"
)

var Psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type PostgresTransactionRepository struct {
	db *pgxpool.Pool
}

//func NewTransactionRepository(db *pgxpool.Pool) repository.TransactionRepository {
//	return &PostgresTransactionRepository{db: db}
//}

func (r *PostgresTransactionRepository) Save(ctx context.Context, transactions model.Transaction) error {
	query, args, err := Psql.
		Insert("transactions").
		Columns("id", "user_id", "amount", "currency", "status", "timestamp").
		Values(transactions.ID, transactions.UserID, transactions.Amount, transactions.Currency, transactions.Status, transactions.Timestamp).
		ToSql()
	if err != nil {
		return err
	}

	_, execErr := r.db.Exec(ctx, query, args...)
	return execErr
}
