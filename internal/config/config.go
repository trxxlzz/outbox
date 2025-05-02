package config

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string

	// PostgreSQL
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string

	// Kafka
	KafkaBrokers string
	KafkaTopic   string
}

type DBType string

const (
	PostgreSQL DBType = "postgres"
)

// LoadConfig загружает конфиг из .env файла и переменных окружения
func LoadConfig(envFile string, dbType DBType) (*Config, error) {
	if err := godotenv.Load(envFile); err != nil {
		return nil, fmt.Errorf("error loading env file %s: %v", envFile, err)
	}

	config := &Config{
		AppEnv:        os.Getenv("APP_ENV"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		KafkaBrokers:  os.Getenv("KAFKA_BROKERS"),
		KafkaTopic:    os.Getenv("KAFKA_TOPIC"),
	}

	if err := config.validate(dbType); err != nil {
		return nil, err
	}

	return config, nil
}

// validate проверяет, что все обязательные поля заполнены
func (c *Config) validate(dbType DBType) error {
	missingFields := []string{}

	if c.KafkaBrokers == "" {
		missingFields = append(missingFields, "KAFKA_BROKERS")
	}
	if c.KafkaTopic == "" {
		missingFields = append(missingFields, "KAFKA_TOPIC")
	}
	if c.RedisHost == "" {
		missingFields = append(missingFields, "REDIS_HOST")
	}
	if c.RedisPort == "" {
		missingFields = append(missingFields, "REDIS_PORT")
	}

	if dbType == PostgreSQL {
		if c.DBHost == "" {
			missingFields = append(missingFields, "DB_HOST")
		}
		if c.DBPort == "" {
			missingFields = append(missingFields, "DB_PORT")
		}
		if c.DBUser == "" {
			missingFields = append(missingFields, "DB_USER")
		}
		if c.DBPassword == "" {
			missingFields = append(missingFields, "DB_PASSWORD")
		}
		if c.DBName == "" {
			missingFields = append(missingFields, "DB_NAME")
		}
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required config fields: %v", missingFields)
	}

	return nil
}

// DSN возвращает строку подключения к PostgreSQL
func (c *Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

// RedisAddr возвращает строку подключения к Redis
func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

// KafkaBrokersList возвращает список брокеров Kafka как []string
func (c *Config) KafkaBrokersList() []string {
	return strings.Split(c.KafkaBrokers, ",")
}

// NewRedisClient возвращает готовый redis.Client
func (c *Config) NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     c.RedisAddr(),
		Password: c.RedisPassword,
		DB:       0,
	})
}
