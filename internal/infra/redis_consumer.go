package infra

import (
	"context"
	"github.com/redis/go-redis/v9"
	"outbox/internal/mapper"
	"outbox/internal/model"
)

type RedisConsumer struct {
	client   *redis.Client
	stream   string
	group    string
	consumer string
}

func NewRedisConsumer(client *redis.Client, stream, group, consumer string) *RedisConsumer {
	return &RedisConsumer{client: client, stream: stream, group: group, consumer: consumer}
}

func (c *RedisConsumer) Listen(ctx context.Context, handler func(tx model.Transaction)) error {
	_ = c.client.XGroupCreateMkStream(ctx, c.stream, c.group, "$").Err()
	for {
		result, err := c.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    c.group,
			Consumer: c.consumer,
			Streams:  []string{c.stream, ">"},
			Count:    10,
			Block:    0,
		}).Result()
		if err != nil {
			return err
		}

		for _, stream := range result {
			for _, message := range stream.Messages {
				tx := mapper.MapToTransaction(message.Values)
				handler(tx)
				c.client.XAck(ctx, c.stream, c.group, message.ID)
			}
		}
	}
}
