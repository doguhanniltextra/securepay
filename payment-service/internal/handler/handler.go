package handler

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"securepay/payment-service/internal/kafka"
	"securepay/payment-service/internal/repository"
	"securepay/payment-service/internal/validator"
	"securepay/payment-service/models"
	pb "securepay/proto/gen/go/payment/v1"
)

// PaymentHandler implements pb.PaymentServiceServer
type PaymentHandler struct {
	pb.UnimplementedPaymentServiceServer
	repo      repository.Repository
	validator *validator.Validator
	producer  *kafka.Producer
}

// NewPaymentHandler creates a new PaymentHandler
func NewPaymentHandler(repo repository.Repository, val *validator.Validator, producer *kafka.Producer) *PaymentHandler {
	return &PaymentHandler{
		repo:      repo,
		validator: val,
		producer:  producer,
	}
}

func (h *PaymentHandler) InitiatePayment(ctx context.Context, req *pb.InitiatePaymentRequest) (*pb.InitiatePaymentResponse, error) {
	slog.Info("InitiatePayment called", "payment_id", req.PaymentId, "amount", req.Amount)

	// Validator
	if err := h.validator.ValidateInitiatePayment(req); err != nil {
		slog.Error("Validation failed", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// TODO: Idempotency Check

	// TODO: Balance Check (via Account Service gRPC)

	// Save to DB
	if err := h.repo.SavePayment(ctx, req); err != nil {
		slog.Error("Failed to save payment", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to save payment: %v", err)
	}

	// Create Kafka Event
	event := models.PaymentInitiatedEvent{
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

	// Fetch from DB
	payment, err := h.repo.GetPayment(ctx, req.PaymentId)
	if err != nil {
		slog.Error("Failed to get payment", "error", err)
		return nil, status.Errorf(codes.NotFound, "payment not found: %v", err)
	}

	// Map status string to enum using models constants
	var status pb.PaymentStatus
	switch models.PaymentStatus(payment.Status) {
	case models.StatusPending:
		status = pb.PaymentStatus_PENDING
	case models.StatusCompleted:
		status = pb.PaymentStatus_COMPLETED
	case models.StatusFailed:
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
