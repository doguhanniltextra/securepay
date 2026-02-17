package handler

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"securepay/account-service/internal/repository"
	pb "securepay/proto/gen/go/account/v1"
)

// AccountHandler implements pb.AccountServiceServer
type AccountHandler struct {
	pb.UnimplementedAccountServiceServer
	repo repository.Repository
}

// NewAccountHandler creates a new AccountHandler
func NewAccountHandler(repo repository.Repository) *AccountHandler {
	return &AccountHandler{repo: repo}
}

// CheckBalance get account balance
func (h *AccountHandler) CheckBalance(ctx context.Context, req *pb.CheckBalanceRequest) (*pb.CheckBalanceResponse, error) {
	slog.Info("CheckBalance called", "account_id", req.AccountId)

	if req.AccountId == "" {
		return nil, status.Error(codes.InvalidArgument, "account_id is required")
	}

	acc, err := h.repo.GetAccount(ctx, req.AccountId)
	if err != nil {
		slog.Error("Failed to get account", "error", err)
		return nil, status.Errorf(codes.NotFound, "account not found: %v", err)
	}

	return &pb.CheckBalanceResponse{
		AccountId: acc.ID,
		Balance:   acc.Balance,
		Currency:  acc.Currency,
	}, nil
}
