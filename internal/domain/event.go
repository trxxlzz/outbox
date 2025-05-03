package domain

import (
	"context"
)

type EventProducer interface {
	Publish(ctx context.Context, key string, payload []byte) error
}
