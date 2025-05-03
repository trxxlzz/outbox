package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisConsumer struct {
	client   *redis.Client
	stream   string
	group    string
	consumer string
}

func NewRedisConsumer(client *redis.Client, stream, group, consumer string) *RedisConsumer {
	return &RedisConsumer{
		client:   client,
		stream:   stream,
		group:    group,
		consumer: consumer,
	}
}

func (c *RedisConsumer) InitGroup(ctx context.Context) error {
	return c.client.XGroupCreateMkStream(ctx, c.stream, c.group, "$").Err()
}

func (c *RedisConsumer) ReadMessages(ctx context.Context) ([]redis.XMessage, error) {
	result, err := c.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    c.group,
		Consumer: c.consumer,
		Streams:  []string{c.stream, ">"},
		Count:    10,
		Block:    0,
	}).Result()
	if err != nil {
		return nil, err
	}

	var messages []redis.XMessage
	for _, stream := range result {
		messages = append(messages, stream.Messages...)
	}
	return messages, nil
}

func (c *RedisConsumer) Ack(ctx context.Context, messageID string) error {
	return c.client.XAck(ctx, c.stream, c.group, messageID).Err()
}
