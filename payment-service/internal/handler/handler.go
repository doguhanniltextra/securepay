package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"securepay/payment-service/internal/cache"
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
	cache     cache.Cache
}

// NewPaymentHandler creates a new PaymentHandler
func NewPaymentHandler(repo repository.Repository, val *validator.Validator, producer *kafka.Producer, cache cache.Cache) *PaymentHandler {
	return &PaymentHandler{
		repo:      repo,
		validator: val,
		producer:  producer,
		cache:     cache,
	}
}

func (h *PaymentHandler) InitiatePayment(ctx context.Context, req *pb.InitiatePaymentRequest) (*pb.InitiatePaymentResponse, error) {
	ctx, span := otel.Tracer("payment-service").Start(ctx, "handler.InitiatePayment")
	defer span.End()

	slog.InfoContext(ctx, "InitiatePayment called", "payment_id", req.PaymentId, "amount", req.Amount)

	// Validator
	if err := h.validator.ValidateInitiatePayment(req); err != nil {
		slog.ErrorContext(ctx, "Validation failed", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Idempotency Check
	idempotencyKey := fmt.Sprintf("idempotency:%s", req.IdempotencyKey)
	cachedResp, err := h.cache.Get(ctx, idempotencyKey)
	if err == nil && cachedResp != "" {
		slog.InfoContext(ctx, "Returning cached response for idempotency", "key", req.IdempotencyKey)
		var resp pb.InitiatePaymentResponse
		if err := json.Unmarshal([]byte(cachedResp), &resp); err == nil {
			return &resp, nil
		}
		slog.WarnContext(ctx, "Failed to unmarshal cached response", "error", err)
	}

	// TODO: Balance Check (via Account Service gRPC)

	// Save to DB
	if err := h.repo.SavePayment(ctx, req); err != nil {
		slog.ErrorContext(ctx, "Failed to save payment", "error", err)
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
		slog.ErrorContext(ctx, "Failed to produce payment initiated event", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to produce payment event: %v", err)
	}

	// Prepare PENDING response
	resp := &pb.InitiatePaymentResponse{
		PaymentId: req.PaymentId,
		Status:    pb.PaymentStatus_PENDING,
		Message:   "Payment initiated",
	}

	// Save idempotency record to Redis
	respJSON, _ := json.Marshal(resp)
	if err := h.cache.Set(ctx, idempotencyKey, string(respJSON), 24*time.Hour); err != nil {
		slog.WarnContext(ctx, "Failed to set idempotency key in cache", "error", err)
	}

	slog.InfoContext(ctx, "Payment initiated successfully", "payment_id", req.PaymentId)
	return resp, nil
}

func (h *PaymentHandler) GetPayment(ctx context.Context, req *pb.GetPaymentRequest) (*pb.GetPaymentResponse, error) {
	ctx, span := otel.Tracer("payment-service").Start(ctx, "handler.GetPayment")
	defer span.End()

	slog.InfoContext(ctx, "GetPayment called", "payment_id", req.PaymentId)

	if req.PaymentId == "" {
		return nil, status.Error(codes.InvalidArgument, "payment_id is required")
	}

	// Fetch from DB
	payment, err := h.repo.GetPayment(ctx, req.PaymentId)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get payment", "error", err)
		return nil, status.Errorf(codes.NotFound, "payment not found: %v", err)
	}

	// Map status string to enum using models constants
	var paymentStatus pb.PaymentStatus
	switch models.PaymentStatus(payment.Status) {
	case models.StatusPending:
		paymentStatus = pb.PaymentStatus_PENDING
	case models.StatusCompleted:
		paymentStatus = pb.PaymentStatus_COMPLETED
	case models.StatusFailed:
		paymentStatus = pb.PaymentStatus_FAILED
	default:
		paymentStatus = pb.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
	}

	return &pb.GetPaymentResponse{
		PaymentId:   payment.ID,
		Status:      paymentStatus,
		Message:     "Payment details retrieved",
		Amount:      payment.Amount,
		Currency:    payment.Currency,
		FromAccount: payment.FromAccount,
		ToAccount:   payment.ToAccount,
	}, nil
}
