package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	kgo "github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"securepay/account-service/config"
	"securepay/account-service/internal/cache"
	"securepay/account-service/internal/repository"
	"securepay/account-service/models"
)

type Consumer struct {
	reader *kgo.Reader
}

func NewConsumer(cfg *config.Config) *Consumer {
	reader := kgo.NewReader(kgo.ReaderConfig{
		Brokers:  cfg.KafkaBrokers,
		Topic:    cfg.KafkaTopic,
		GroupID:  "account-service-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{reader: reader}
}

// Start begins consuming payment.initiated events.
// After processing, it produces a result event via the resultProducer.
func (c *Consumer) Start(ctx context.Context, repo repository.Repository, balanceCache cache.Cache, resultProducer *Producer) {
	slog.InfoContext(ctx, "Starting Kafka Consumer", "topic", c.reader.Config().Topic)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// FetchMessage blocks until message received
				m, err := c.reader.FetchMessage(ctx)
				if err != nil {
					// Prepare for shutdown or transient error
					if ctx.Err() != nil {
						return // Context closed
					}
					slog.ErrorContext(ctx, "Failed to fetch message", "error", err)
					time.Sleep(time.Second) // Backoff
					continue
				}

				// Extract Trace Context from Kafka headers
				carrier := propagation.MapCarrier{}
				for _, h := range m.Headers {
					carrier[h.Key] = string(h.Value)
				}
				msgCtx := otel.GetTextMapPropagator().Extract(ctx, carrier)
				msgCtx, span := otel.Tracer("account-service").Start(msgCtx, "kafka.ConsumePaymentInitiatedEvent")

				slog.InfoContext(msgCtx, "Message received",
					"key", string(m.Key),
					"offset", m.Offset,
					"trace_id", span.SpanContext().TraceID().String(),
				)

				var event models.PaymentInitiatedEvent
				if err := json.Unmarshal(m.Value, &event); err != nil {
					slog.ErrorContext(msgCtx, "Failed to unmarshal event", "error", err)
					span.RecordError(err)
					span.End()
					c.reader.CommitMessages(msgCtx, m)
					continue
				}

				// Process Payment (Deduct & Credit Balance)
				resultEvent := models.PaymentResultEvent{
					PaymentID: event.PaymentID,
					Timestamp: time.Now().Format(time.RFC3339),
				}

				err = repo.ProcessPayment(msgCtx, event.FromAccount, event.ToAccount, event.Amount)
				if err != nil {
					slog.ErrorContext(msgCtx, "Failed to process payment",
						"error", err,
						"payment_id", event.PaymentID,
					)
					resultEvent.Status = "FAILED"
					resultEvent.Reason = err.Error()
				} else {
					slog.InfoContext(msgCtx, "Payment processed successfully", "payment_id", event.PaymentID)
					resultEvent.Status = "COMPLETED"

					// Invalidate Redis cache for both accounts
					if delErr := balanceCache.DeleteBalance(msgCtx, event.FromAccount); delErr != nil {
						slog.WarnContext(msgCtx, "Failed to invalidate cache for from_account",
							"account_id", event.FromAccount, "error", delErr,
						)
					} else {
						slog.InfoContext(msgCtx, "Cache invalidated", "account_id", event.FromAccount)
					}

					if delErr := balanceCache.DeleteBalance(msgCtx, event.ToAccount); delErr != nil {
						slog.WarnContext(msgCtx, "Failed to invalidate cache for to_account",
							"account_id", event.ToAccount, "error", delErr,
						)
					} else {
						slog.InfoContext(msgCtx, "Cache invalidated", "account_id", event.ToAccount)
					}
				}

				// Produce result event so payment-service can update status
				if resultProducer != nil {
					if pubErr := resultProducer.ProducePaymentResultEvent(msgCtx, resultEvent); pubErr != nil {
						slog.ErrorContext(msgCtx, "Failed to produce result event",
							"error", pubErr, "payment_id", event.PaymentID,
						)
						span.RecordError(pubErr)
					}
				}

				// Commit message after processing
				if err := c.reader.CommitMessages(msgCtx, m); err != nil {
					slog.ErrorContext(msgCtx, "Failed to commit message", "error", err)
					span.RecordError(err)
				}
				span.End()
			}
		}
	}()
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
