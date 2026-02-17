package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"

	"securepay/account-service/config"
	"securepay/account-service/internal/repository"
	"securepay/account-service/models"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(cfg *config.Config) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.KafkaBrokers,
		Topic:    cfg.KafkaTopic,
		GroupID:  "account-service-group", // Convention
		MinBytes: 10e3,                    // 10KB
		MaxBytes: 10e6,                    // 10MB
		// CommitInterval: time.Second, // Auto-commit interval? Or manual?
		// StartOffset: kafka.LastOffset, // Only new messages? Or explicit?
	})

	return &Consumer{reader: reader}
}

func (c *Consumer) Start(ctx context.Context, repo repository.Repository) {
	slog.Info("Starting Kafka Consumer", "topic", c.reader.Config().Topic)
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
					slog.Error("Failed to fetch message", "error", err)
					time.Sleep(time.Second) // Backoff
					continue
				}

				slog.Info("Message received", "key", string(m.Key), "offset", m.Offset)

				var event models.PaymentInitiatedEvent
				if err := json.Unmarshal(m.Value, &event); err != nil {
					slog.Error("Failed to unmarshal event", "error", err)
					// Commit to move on? or Dead Letter Queue?
					// For now, commit to avoid stuck loop on bad message
					c.reader.CommitMessages(ctx, m)
					continue
				}

				// Process Payment (Deduct Balance)
				// We need "From" and "To" and "Amount"
				err = repo.ProcessPayment(ctx, event.FromAccount, event.ToAccount, event.Amount)
				if err != nil {
					slog.Error("Failed to process payment", "error", err, "payment_id", event.PaymentID)
					// If insufficient funds, we should probably record failure somewhere
					// But we have no producer. Logging is best effort.
				} else {
					slog.Info("Payment processed successfully", "payment_id", event.PaymentID)
				}
				
				// Commit message after processing
				if err := c.reader.CommitMessages(ctx, m); err != nil {
					slog.Error("Failed to commit message", "error", err)
				}
			}
		}
	}()
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
