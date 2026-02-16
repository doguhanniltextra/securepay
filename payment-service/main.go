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

	"securepay/payment-service/config"
	"securepay/payment-service/internal/handler"
	"securepay/payment-service/internal/kafka"
	"securepay/payment-service/internal/repository"
	"securepay/payment-service/internal/spiffe"
	"securepay/payment-service/internal/validator"
	pb "securepay/proto/gen/go/payment/v1"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load configuration
	cfg := config.Load()

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

	// Initialize Kafka Producer
	producer := kafka.NewProducer(cfg.KafkaBrokers, cfg.KafkaTopic)
	defer producer.Close()

	// Initialize SPIFFE Workload API Source
	source, err := spiffe.InitSPIFFESource(context.Background(), cfg.SpiffeSocket)
	if err != nil {
		slog.Error("Failed to initialize SPIFFE source", "error", err)
		os.Exit(1)
	}
	defer source.Close()
	slog.Info("SPIFFE Source initialized successfully")

	// Initialize Components
	repo := repository.NewPostgresRepository(db)
	val := validator.New()
	h := handler.NewPaymentHandler(repo, val, producer)

	// Create gRPC server with mTLS credentials
	creds := spiffe.PaymentServiceServerCredentials(source)
	s := grpc.NewServer(creds)
	
	// Register PaymentService
	pb.RegisterPaymentServiceServer(s, h)
	
	// Enable reflection for debugging (e.g. grpcurl)
	reflection.Register(s)

	// Graceful shutdown handling
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		slog.Info("Shutting down gRPC server...")
		s.GracefulStop()
		slog.Info("Server stopped")
	}()

	slog.Info("Starting Payment Service gRPC server", "port", cfg.Port)
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
