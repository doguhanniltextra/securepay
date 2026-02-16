package main

import (
	"context"
	"database/sql"

	"log/slog"
	"time"


	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "securepay/proto/gen/go/payment/v1"
)

// PaymentHandler implements pb.PaymentServiceServer
type PaymentHandler struct {
	pb.UnimplementedPaymentServiceServer
	repo      *PostgresRepository
	validator *Validator
	producer  *KafkaProducer
}

// NewPaymentHandler creates a new PaymentHandler
func NewPaymentHandler(db *sql.DB, producer *KafkaProducer) *PaymentHandler {
	return &PaymentHandler{
		repo:      NewPostgresRepository(db),
		validator: NewValidator(),
		producer:  producer,
	}
}

func (h *PaymentHandler) InitiatePayment(ctx context.Context, req *pb.InitiatePaymentRequest) (*pb.InitiatePaymentResponse, error) {
	slog.Info("InitiatePayment called", "payment_id", req.PaymentId, "amount", req.Amount)

	// Task 23 - Validator
	if err := h.validator.ValidateInitiatePayment(req); err != nil {
		slog.Error("Validation failed", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// TODO: Idempotency Check

	// TODO: Balance Check (via Account Service gRPC)

	// TODO: Task 22 - Repository (Save to DB)
	if err := h.repo.SavePayment(ctx, req); err != nil {
		slog.Error("Failed to save payment", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to save payment: %v", err)
	}

	// Task 25 - Kafka Producer
	event := PaymentInitiatedEvent{
		PaymentID:   req.PaymentId,
		FromAccount: req.FromAccount,
		ToAccount:   req.ToAccount,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	if err := h.producer.ProducePaymentInitiatedEvent(ctx, event); err != nil {
		slog.Error("Failed to produce payment initiated event", "error", err)
		// We return error here to signal failure to the client, 
		// although DB record is already created (in PENDING state).
		// In a real system, we might want to transactions or outbox pattern.
		return nil, status.Errorf(codes.Internal, "failed to produce payment event: %v", err)
	}

	// Return PENDING response
	return &pb.InitiatePaymentResponse{
		PaymentId: req.PaymentId,
		Status:    pb.PaymentStatus_PENDING,
		Message:   "Payment initiated",
	}, nil
}

func (h *PaymentHandler) GetPayment(ctx context.Context, req *pb.GetPaymentRequest) (*pb.GetPaymentResponse, error) {
	slog.Info("GetPayment called", "payment_id", req.PaymentId)

	if req.PaymentId == "" {
		return nil, status.Error(codes.InvalidArgument, "payment_id is required")
	}

	// TODO: Task 22 - Repository (Fetch from DB)
	payment, err := h.repo.GetPayment(ctx, req.PaymentId)
	if err != nil {
		slog.Error("Failed to get payment", "error", err)
		return nil, status.Errorf(codes.NotFound, "payment not found: %v", err)
	}
	
	// Map status string to enum
	var status pb.PaymentStatus
	switch payment.Status {
	case string(StatusPending):
		status = pb.PaymentStatus_PENDING
	case string(StatusCompleted):
		status = pb.PaymentStatus_COMPLETED
	case string(StatusFailed):
		status = pb.PaymentStatus_FAILED
	default:
		status = pb.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
	}

	return &pb.GetPaymentResponse{
		PaymentId:   payment.ID,
		Status:      status,
		Message:     "Payment details retrieved",
		Amount:      payment.Amount,
		Currency:    payment.Currency,
		FromAccount: payment.FromAccount,
		ToAccount:   payment.ToAccount,
	}, nil
}
