package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"google.golang.org/grpc"

	// Import generated code packages
	accountv1 "securepay/proto/gen/go/account/v1"
	paymentv1 "securepay/proto/gen/go/payment/v1"
)

// PaymentServiceAddr is the address of the Payment Service.
const PaymentServiceAddr = "payment-service:8081"

// AccountServiceAddr is the address of the Account Service.
const AccountServiceAddr = "account-service:50051"

// NewPaymentServiceClient creates a new gRPC client for the Payment Service.
// It uses SPIFFE-based mTLS for secure communication.
func NewPaymentServiceClient(ctx context.Context, source *workloadapi.X509Source) (paymentv1.PaymentServiceClient, *grpc.ClientConn, error) {
	// Get mTLS credentials securely using SPIFFE.
	creds, err := PaymentServiceCredentials(source)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get payment service credentials: %w", err)
	}

	// Dial the Payment Service using gRPC with mTLS.
	// Blocks until connection is established or context times out.
	conn, err := grpc.DialContext(ctx, PaymentServiceAddr, creds, grpc.WithBlock())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial payment service: %w", err)
	}

	slog.Info("Connected to Payment Service", "address", PaymentServiceAddr)

	// Create and return the generated gRPC client.
	client := paymentv1.NewPaymentServiceClient(conn)
	return client, conn, nil
}

// NewAccountServiceClient creates a new gRPC client for the Account Service.
// It uses SPIFFE-based mTLS for secure communication.
func NewAccountServiceClient(ctx context.Context, source *workloadapi.X509Source) (accountv1.AccountServiceClient, *grpc.ClientConn, error) {
	// Get mTLS credentials securely using SPIFFE.
	creds, err := AccountServiceCredentials(source)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get account service credentials: %w", err)
	}

	// Dial the Account Service using gRPC with mTLS.
	// Blocks until connection is established or context times out.
	conn, err := grpc.DialContext(ctx, AccountServiceAddr, creds, grpc.WithBlock())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial account service: %w", err)
	}

	slog.Info("Connected to Account Service", "address", AccountServiceAddr)

	// Create and return the generated gRPC client.
	client := accountv1.NewAccountServiceClient(conn)
	return client, conn, nil
}
