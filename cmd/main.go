package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"outbox/internal/config"
	"outbox/internal/infra"
	"outbox/internal/infra/kafka"
	"outbox/internal/infra/postgres"
	"outbox/internal/service/transaction"
	"outbox/internal/transport"
	txRepo "outbox/repository/transaction"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig("config/.env", "postgres")

	dbpool, err := postgres.NewDBConnection(ctx, cfg.DSN())
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer dbpool.Close()
	log.Println("Successfully connected to database")

	// Подключаем Kafka
	kafkaWriter := kafka.NewKafkaWriter(cfg.KafkaBrokersList(), cfg.KafkaTopic)
	defer kafkaWriter.Close()
	log.Println("Successfully connected to Kafka")

	// Подключаем Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr(),
		Password: cfg.RedisPassword,
	})
	defer redisClient.Close()
	log.Println("Successfully connected to Redis")

	// Создаём репозиторий
	repo := txRepo.NewTransactionRepository(dbpool)

	// Создаём Kafka продюсер
	producer := infra.NewKafkaProducer(kafkaWriter)

	// Создаём бизнес-слой
	service := transaction.NewTransactionService(repo, producer)

	// Создаём Redis consumer
	redisConsumer := infra.NewRedisConsumer(redisClient, "transactions", "tx_group", "consumer1")

	// Запускаем transport-слой
	err = transport.StartConsumer(ctx, redisConsumer, service)
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}
}
