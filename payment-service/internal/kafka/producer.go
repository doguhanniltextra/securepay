package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"securepay/payment-service/models"
)


// Producer wrapper
type Producer struct {
	writer *kafka.Writer
}

// NewProducer creates a new producer
func NewProducer(brokers []string, topic string) *Producer {
	slog.Info("Initializing Kafka Producer", "brokers", brokers, "topic", topic)

	return &Producer{
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
func (kp *Producer) Close() error {
	return kp.writer.Close()
}

// ProducePaymentInitiatedEvent sends the event to Kafka
func (kp *Producer) ProducePaymentInitiatedEvent(ctx context.Context, event models.PaymentInitiatedEvent) error {
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
