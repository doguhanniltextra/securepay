package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"securepay/account-service/config"
	"securepay/account-service/internal/cache"
	"securepay/account-service/internal/handler"
	"securepay/account-service/internal/kafka"
	"securepay/account-service/internal/logger"
	"securepay/account-service/internal/repository"
	"securepay/account-service/internal/spiffe"
	"securepay/account-service/internal/telemetry"
	"securepay/account-service/models"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	pb "securepay/proto/gen/go/account/v1"
)

func main() {
	// Use the OTel-aware JSON logger so that trace_id and span_id are
	// automatically injected into every log record that carries a span.
	slog.SetDefault(logger.New())

	// Load Configuration
	cfg := config.Load()

	// Initialize Tracer
	shutdownTracer, err := telemetry.InitTracer(context.Background(), "account-service")
	if err != nil {
		slog.Error("Failed to initialize tracer", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := shutdownTracer(context.Background()); err != nil {
			slog.Error("Failed to shutdown tracer", "error", err)
		}
	}()

	// Connect to Database
	if cfg.DatabaseURL == "" {
		slog.Error("DATABASE_URL is not set")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		slog.Error("Failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		slog.Error("Failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("Connected to database")

	repo := repository.NewPostgresRepository(db)

	// Seed Data
	seedAccounts(context.Background(), repo)

	// Initialize Redis Cache
	balanceCache := cache.NewRedisCache(cfg.RedisAddr, cfg.RedisPassword)
	slog.Info("Redis cache initialized", "addr", cfg.RedisAddr)

	// Initialize Kafka Consumer
	consumer := kafka.NewConsumer(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	consumer.Start(ctx, repo, balanceCache)

	// Initialize SPIFFE Workload API Source
	source, err := spiffe.InitSPIFFESource(ctx, cfg.SpiffeSocket)
	if err != nil {
		slog.Error("Failed to initialize SPIFFE source", "error", err)
		os.Exit(1)
	}
	defer source.Close()
	slog.Info("SPIFFE Source initialized successfully")

	// Create gRPC Server with mTLS
	creds := spiffe.ServerCredentials(source)
	s := grpc.NewServer(creds, grpc.StatsHandler(otelgrpc.NewServerHandler()))

	// Register AccountService
	h := handler.NewAccountHandler(repo, balanceCache)
	pb.RegisterAccountServiceServer(s, h)

	// Enable reflection
	reflection.Register(s)

	// Graceful Shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		slog.Info("Shutting down gRPC server...")
		s.GracefulStop()
		consumer.Close() // Close Kafka reader
		cancel()         // Cancel context for consumer loop
		slog.Info("Server stopped")
	}()

	slog.Info("Starting Account Service gRPC server", "port", cfg.Port)
	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		slog.Error("Failed to listen", "error", err)
		os.Exit(1)
	}

	if err := s.Serve(lis); err != nil {
		slog.Error("Failed to serve", "error", err)
		os.Exit(1)
	}
}

func seedAccounts(ctx context.Context, repo repository.Repository) {
	accounts := []models.Account{
		{
			ID:       "11111111-1111-1111-1111-111111111111", // Ahmet
			Balance:  1000.00,
			Currency: "TRY",
		},
		{
			ID:       "22222222-2222-2222-2222-222222222222", // Mehmet
			Balance:  500.00,
			Currency: "TRY",
		},
	}

	for _, acc := range accounts {
		if err := repo.UpsertAccount(ctx, &acc); err != nil {
			slog.Error("Failed to seed account", "id", acc.ID, "error", err)
		} else {
			slog.Info("Seeded account", "id", acc.ID)
		}
	}
}
