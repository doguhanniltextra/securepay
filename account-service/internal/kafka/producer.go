package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	kgo "github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"securepay/account-service/models"
)

// Producer wraps a kafka.Writer to produce payment result events.
type Producer struct {
	writer *kgo.Writer
}

// NewProducer creates a new Kafka producer for the given topic.
func NewProducer(brokers []string, topic string) *Producer {
	slog.Info("Initializing Kafka Result Producer", "brokers", brokers, "topic", topic)

	return &Producer{
		writer: &kgo.Writer{
			Addr:     kgo.TCP(brokers...),
			Topic:    topic,
			Balancer: &kgo.LeastBytes{},
		},
	}
}

// Close closes the producer.
func (p *Producer) Close() error {
	return p.writer.Close()
}

// ProducePaymentResultEvent sends a payment result event to Kafka.
func (p *Producer) ProducePaymentResultEvent(ctx context.Context, event models.PaymentResultEvent) error {
	ctx, span := otel.Tracer("account-service").Start(ctx, "kafka.ProducePaymentResultEvent")
	defer span.End()

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal result event: %w", err)
	}

	// Inject trace context into Kafka headers
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	headers := make([]kgo.Header, 0, len(carrier))
	for k, v := range carrier {
		headers = append(headers, kgo.Header{Key: k, Value: []byte(v)})
	}

	msg := kgo.Message{
		Key:     []byte(event.PaymentID),
		Value:   payload,
		Headers: headers,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write result message: %w", err)
	}

	slog.InfoContext(ctx, "Produced payment result event",
		"payment_id", event.PaymentID,
		"status", event.Status,
	)
	return nil
}
