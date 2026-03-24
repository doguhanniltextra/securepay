package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	kgo "github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"securepay/payment-service/internal/repository"
	"securepay/payment-service/models"
)

// ResultConsumer consumes payment result events from account-service
// and updates the payment status in the database.
type ResultConsumer struct {
	reader *kgo.Reader
}

// NewResultConsumer creates a consumer for the payment.result topic.
func NewResultConsumer(brokers []string, topic string) *ResultConsumer {
	reader := kgo.NewReader(kgo.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  "payment-service-result-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &ResultConsumer{reader: reader}
}

// Start begins consuming payment result events and updating payment status.
func (c *ResultConsumer) Start(ctx context.Context, repo repository.Repository) {
	slog.Info("Starting Payment Result Consumer", "topic", c.reader.Config().Topic)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				m, err := c.reader.FetchMessage(ctx)
				if err != nil {
					if ctx.Err() != nil {
						return
					}
					slog.Error("Failed to fetch result message", "error", err)
					time.Sleep(time.Second)
					continue
				}

				// Extract Trace Context
				carrier := propagation.MapCarrier{}
				for _, h := range m.Headers {
					carrier[h.Key] = string(h.Value)
				}
				msgCtx := otel.GetTextMapPropagator().Extract(ctx, carrier)
				msgCtx, span := otel.Tracer("payment-service").Start(msgCtx, "kafka.ConsumePaymentResultEvent")

				slog.InfoContext(msgCtx, "Payment result received",
					"key", string(m.Key),
					"offset", m.Offset,
					"trace_id", span.SpanContext().TraceID().String(),
				)

				var event models.PaymentResultEvent
				if err := json.Unmarshal(m.Value, &event); err != nil {
					slog.ErrorContext(msgCtx, "Failed to unmarshal result event", "error", err)
					span.RecordError(err)
					span.End()
					c.reader.CommitMessages(msgCtx, m)
					continue
				}

				// Update payment status in database
				newStatus := models.PaymentStatus(event.Status)
				if err := repo.UpdatePaymentStatus(msgCtx, event.PaymentID, newStatus); err != nil {
					slog.ErrorContext(msgCtx, "Failed to update payment status",
						"error", err,
						"payment_id", event.PaymentID,
						"status", event.Status,
					)
					span.RecordError(err)
				} else {
					slog.InfoContext(msgCtx, "Payment status updated",
						"payment_id", event.PaymentID,
						"status", event.Status,
					)
				}

				// Commit
				if err := c.reader.CommitMessages(msgCtx, m); err != nil {
					slog.ErrorContext(msgCtx, "Failed to commit result message", "error", err)
					span.RecordError(err)
				}
				span.End()
			}
		}
	}()
}

// Close closes the consumer.
func (c *ResultConsumer) Close() error {
	return c.reader.Close()
}
