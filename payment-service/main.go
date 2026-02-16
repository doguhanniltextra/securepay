package main

import (

	"database/sql"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "securepay/proto/gen/go/payment/v1"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found")
	}

	// Port 8081 as requested
	port := ":8081"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		slog.Error("Failed to listen", "error", err)
		os.Exit(1)
	}


	// Connect to Database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		slog.Error("DATABASE_URL is not set")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", dbURL)
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
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		kafkaTopic = "payment.initiated"
	}
	producer := NewKafkaProducer(nil, kafkaTopic)
	defer producer.Close()

	// Create gRPC server
	s := grpc.NewServer()
	
	// Register PaymentService
	handler := NewPaymentHandler(db, producer)
	pb.RegisterPaymentServiceServer(s, handler)
	
	// Enable reflection for debugging (e.g.grpcurl)
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

	slog.Info("Starting Payment Service gRPC server", "port", port)
	if err := s.Serve(lis); err != nil {
		slog.Error("Failed to serve", "error", err)
		os.Exit(1)
	}
}
