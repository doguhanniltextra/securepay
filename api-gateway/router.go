package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"securepay/api-gateway/endpoints"
	"securepay/api-gateway/middleware"
	accountv1 "securepay/proto/gen/go/account/v1"
	paymentv1 "securepay/proto/gen/go/payment/v1"
)

// NewRouter sets up the routes and middleware for the API Gateway.
// It accepts gRPC clients as dependencies.
func NewRouter(paymentClient paymentv1.PaymentServiceClient, accountClient accountv1.AccountServiceClient) http.Handler {
	mux := http.NewServeMux()

	// Public Health Check
	mux.HandleFunc("GET "+endpoints.HealthCheckPath, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Private API Endpoints (Middleware Applied)

	// POST /api/v1/payments
	mux.Handle(endpoints.InitiatePaymentPathPattern, middlewareChain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleInitiatePayment(w, r, paymentClient)
	})))

	// GET /api/v1/payments/{id}
	mux.Handle(endpoints.GetPaymentPathPattern, middlewareChain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		handleGetPayment(w, r, paymentClient, id)
	})))

	// GET /api/v1/accounts/{id}/balance
	mux.Handle(endpoints.CheckBalancePathPattern, middlewareChain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		handleCheckBalance(w, r, accountClient, id)
	})))

	return mux
}

// middlewareChain applies rate limiting and JWT authentication middleware.
func middlewareChain(next http.Handler) http.Handler {
	// Order: RateLimit (outer) -> JWT (inner) -> Handler
	return middleware.RateLimitMiddleware(middleware.AuthMiddleware(next))
}

func handleInitiatePayment(w http.ResponseWriter, r *http.Request, client paymentv1.PaymentServiceClient) {
	var req paymentv1.InitiatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Context with timeout for gRPC call
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := client.InitiatePayment(ctx, &req)
	if err != nil {
		// In production, map gRPC codes to HTTP status codes
		http.Error(w, fmt.Sprintf("Payment service error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func handleGetPayment(w http.ResponseWriter, r *http.Request, client paymentv1.PaymentServiceClient, id string) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	req := &paymentv1.GetPaymentRequest{PaymentId: id}
	resp, err := client.GetPayment(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Payment service error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func handleCheckBalance(w http.ResponseWriter, r *http.Request, client accountv1.AccountServiceClient, id string) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	req := &accountv1.CheckBalanceRequest{AccountId: id}
	resp, err := client.CheckBalance(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Account service error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
