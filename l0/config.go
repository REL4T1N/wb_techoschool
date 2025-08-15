package main

import (
	"log"
	"os"
)

type Config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	KafkaHost string
	KafkaPort string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func LoadConfig() Config {
	cfg := Config{
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5433"),
		PostgresUser:     getEnv("POSTGRES_USER", "order_user"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "order_pass"),
		PostgresDB:       getEnv("POSTGRES_DB", "orders_db"),

		KafkaHost: getEnv("KAFKA_HOST", "localhost"),
		KafkaPort: getEnv("KAFKA_PORT", "9092"),
	}
	return cfg
}

func MustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable %s is required", key)
	}
	return val
}
