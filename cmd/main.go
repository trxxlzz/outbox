package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"outbox/internal/config"
	kafkaInfra "outbox/internal/infra/kafka"
	"outbox/internal/infra/postgres"
	redis2 "outbox/internal/infra/redis"
	txRepo "outbox/internal/repository/transaction"
	"outbox/internal/service/transaction"
	kafkaTransport "outbox/internal/transport/kafka"
	redisTransport "outbox/internal/transport/redis"
	"time"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig("config/.env", "postgres")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	dbpool, err := postgres.NewDBConnection(ctx, cfg.DSN())
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer dbpool.Close()
	log.Println("Successfully connected to database")

	// Подключаем Kafka
	kafkaWriter := kafkaInfra.NewKafkaWriter(cfg.KafkaBrokersList(), cfg.KafkaTopic)
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
	producer := kafkaTransport.NewKafkaProducer(kafkaWriter)

	// Создаём бизнес-слой
	service := transaction.NewTransactionService(repo, producer, dbpool)

	go service.StartOutboxWorker(ctx, 2*time.Second)

	// Создаём Redis consumer
	redisConsumer := redis2.NewRedisConsumer(redisClient, "transactions", "tx_group", "consumer1")

	// Запуск transport-слоя (слушаем Redis и вызываем бизнес-логику)
	err = redisTransport.StartRedisListener(ctx, redisConsumer, service)
	if err != nil {
		log.Fatalf("Failed to start Redis listener: %v", err)
	}
}
