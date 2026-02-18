package config

import (
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	Port         string
	DatabaseURL  string
	KafkaBrokers []string
	KafkaTopic   string
	SpiffeSocket  string
	RedisAddr     string
	RedisPassword string
}

// Load loads the configuration from environment variables
func Load() *Config {
	// Load .env file if exists (optional)
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using system environment variables")
	}

	return &Config{
		Port:         getEnv("PORT", ":8082"), // Default 8082 for Account Service
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		KafkaBrokers: getEnvList("KAFKA_BROKERS", []string{"localhost:9092"}),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "payment.initiated"),
		SpiffeSocket:  getEnv("SPIFFE_ENDPOINT_SOCKET", "unix:///tmp/spire-agent/public/api.sock"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvList(key string, fallback []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.Split(value, ",")
	}
	return fallback
}
