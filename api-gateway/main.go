package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Logger Init (JSON format for production)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Context with Graceful Shutdown Signals (SIGINT, SIGTERM)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info("Starting API Gateway...", "port", 8080)

	// 1. Initialize SPIFFE X.509 Source
	// This connects to the SPIRE Agent via Unix socket defined in spiffe.go.
	// Ensure SPIRE Agent is running and socket is accessible.
	source, err := InitSPIFFESource(ctx)
	if err != nil {
		slog.Error("Failed to initialize SPIFFE source", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := source.Close(); err != nil {
			slog.Error("Failed to close SPIFFE source", "error", err)
		}
	}()
	slog.Info("SPIFFE Source initialized successfully")

	// 2. Initialize gRPC Clients
	// Use a short timeout for initial connection establishment
	dialCtx, dialCancel := context.WithTimeout(ctx, 5*time.Second)
	defer dialCancel()

	// Payment Service Client
	paymentClient, paymentConn, err := NewPaymentServiceClient(dialCtx, source)
	if err != nil {
		slog.Warn("Failed to connect to Payment Service (continuing without it)", "error", err)
	} else {
		defer paymentConn.Close()
		slog.Info("Payment Service Client initialized")
	}

	// Account Service Client
	accountClient, accountConn, err := NewAccountServiceClient(dialCtx, source)
	if err != nil {
		slog.Warn("Failed to connect to Account Service (continuing without it)", "error", err)
	} else {
		defer accountConn.Close()
		slog.Info("Account Service Client initialized")
	}

	// 3. Setup Router (Inject dependencies)
	router := NewRouter(paymentClient, accountClient)

	// 4. Start HTTP Server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Run server in a specific goroutine
	go func() {
		slog.Info("HTTP Server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP Server failed", "error", err)
			os.Exit(1)
		}
	}()

	// 5. Wait for Shutdown Signal
	<-ctx.Done()
	slog.Info("Shutdown signal received, initiating graceful shutdown...")

	// Create a timeout context for shutdown operations
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("API Gateway shutdown complete")
}
