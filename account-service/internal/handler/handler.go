package handler

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"securepay/account-service/internal/cache"
	"securepay/account-service/internal/repository"
	pb "securepay/proto/gen/go/account/v1"
)

// AccountHandler implements pb.AccountServiceServer
type AccountHandler struct {
	pb.UnimplementedAccountServiceServer
	repo  repository.Repository
	cache cache.Cache
}

// NewAccountHandler creates a new AccountHandler
func NewAccountHandler(repo repository.Repository, cache cache.Cache) *AccountHandler {
	return &AccountHandler{repo: repo, cache: cache}
}

// CheckBalance get account balance (read-aside cache pattern)
func (h *AccountHandler) CheckBalance(ctx context.Context, req *pb.CheckBalanceRequest) (*pb.CheckBalanceResponse, error) {
	slog.Info("CheckBalance called", "account_id", req.AccountId)

	if req.AccountId == "" {
		return nil, status.Error(codes.InvalidArgument, "account_id is required")
	}

	// 1. Check Redis cache
	entry, err := h.cache.GetBalance(ctx, req.AccountId)
	if err != nil {
		slog.Warn("Cache get failed, falling back to DB", "error", err)
	}
	if entry != nil {
		slog.Info("Cache hit", "account_id", req.AccountId)
		return &pb.CheckBalanceResponse{
			AccountId: req.AccountId,
			Balance:   entry.Balance,
			Currency:  entry.Currency,
		}, nil
	}

	// 2. Cache miss -- fetch from PostgreSQL
	slog.Info("Cache miss", "account_id", req.AccountId)
	acc, err := h.repo.GetAccount(ctx, req.AccountId)
	if err != nil {
		slog.Error("Failed to get account", "error", err)
		return nil, status.Errorf(codes.NotFound, "account not found: %v", err)
	}

	// 3. Write to Redis cache with TTL 60s
	cacheEntry := &cache.BalanceEntry{
		Balance:  acc.Balance,
		Currency: acc.Currency,
	}
	if err := h.cache.SetBalance(ctx, req.AccountId, cacheEntry); err != nil {
		slog.Warn("Failed to set cache", "error", err)
	}

	return &pb.CheckBalanceResponse{
		AccountId: acc.ID,
		Balance:   acc.Balance,
		Currency:  acc.Currency,
	}, nil
}
