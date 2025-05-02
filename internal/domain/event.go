package domain

import (
	"context"
	"outbox/internal/model"
)

type EventProducer interface {
	Publish(ctx context.Context, tx model.Transaction) error
}
