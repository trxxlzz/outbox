package transaction

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"outbox/internal/model"
	"outbox/internal/repository"
	modelRepo "outbox/internal/repository/model"
)

var Psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) repository.TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) SaveTransaction(ctx context.Context, tx pgx.Tx, transaction model.Transaction) error {
	query, args, err := Psql.
		Insert("transactions").
		Columns("id", "user_id", "amount", "currency", "status", "timestamp").
		Values(transaction.ID, transaction.UserID, transaction.Amount, transaction.Currency, transaction.Status, transaction.Timestamp).
		ToSql()
	if err != nil {
		return err
	}

	if _, err = tx.Exec(ctx, query, args...); err != nil {
		return err
	}

	if err := r.SaveOutboxEvent(ctx, tx, "TransactionCreated", transaction); err != nil {
		return err
	}

	return nil
}

func (r *TransactionRepository) SaveOutboxEvent(ctx context.Context, tx pgx.Tx, eventType string, payload interface{}) error {
	eventPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	query, args, err := Psql.
		Insert("outbox").
		Columns("event_type", "payload").
		Values(eventType, eventPayload).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	return err
}

func (r *TransactionRepository) GetOutboxEvent(ctx context.Context) (modelRepo.Outbox, error) {
	query, args, err := Psql.
		Select("id", "event_type", "payload").
		From("outbox").
		Where(squirrel.Eq{"status": "new"}).
		Limit(1).
		ToSql()
	if err != nil {
		return modelRepo.Outbox{}, fmt.Errorf("GetOutboxEvent build query: %w", err)
	}

	var outbox modelRepo.Outbox

	row := r.db.QueryRow(ctx, query, args...)
	err = row.Scan(&outbox.ID, &outbox.Type, &outbox.Payload)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return modelRepo.Outbox{}, nil // no new events found
		}
		return modelRepo.Outbox{}, fmt.Errorf("GetOutboxEvent scan: %w", err)
	}

	return outbox, nil
}

func (r *TransactionRepository) MarkOutboxProcessed(ctx context.Context, id int) error {
	query, args, err := Psql.
		Update("outbox").
		Set("status", "done").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("MarkOutboxEventProcessed build query: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("MarkOutboxEventProcessed exec: %w", err)
	}

	return nil
}
