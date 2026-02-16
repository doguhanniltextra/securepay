package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/segmentio/kafka-go"
)

// PaymentInitiatedEvent structure based on conventions
type PaymentInitiatedEvent struct {
	PaymentID   string  `json:"payment_id"`
	FromAccount string  `json:"from_account"`
	ToAccount   string  `json:"to_account"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Timestamp   string  `json:"timestamp"`
}

// KafkaProducer wrapper
type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer creates a new producer
func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	if len(brokers) == 0 {
		// Default to localhost:9092 if not provided, or read from env
		envBroker := os.Getenv("KAFKA_BROKERS")
		if envBroker != "" {
			brokers = strings.Split(envBroker, ",")
		} else {
			brokers = []string{"localhost:9092"}
		}
	}
	
	slog.Info("Initializing Kafka Producer", "brokers", brokers, "topic", topic)

	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
			// Async by default usually, but we might want to check errors.
			// WriteMessages is blocking/sync by default in kafka-go which is good for reliability here.
		},
	}
}

// Close closes the producer
func (kp *KafkaProducer) Close() error {
	return kp.writer.Close()
}

// ProducePaymentInitiatedEvent sends the event to Kafka
func (kp *KafkaProducer) ProducePaymentInitiatedEvent(ctx context.Context, event PaymentInitiatedEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(event.PaymentID), // Use PaymentID as key for ordering guarantees if partitioned
		Value: payload,
	}

	// WriteMessages blocks until the message is sent
	if err := kp.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	slog.Info("Produced event to Kafka", "topic", kp.writer.Topic, "payment_id", event.PaymentID)
	return nil
}
