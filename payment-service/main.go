package main

import (

	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "securepay/proto/gen/go/payment/v1"
)

type server struct {
	pb.UnimplementedPaymentServiceServer
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Port 8081 as requested
	port := ":8081"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		slog.Error("Failed to listen", "error", err)
		os.Exit(1)
	}

	// Create gRPC server
	s := grpc.NewServer()
	
	// Register PaymentService
	pb.RegisterPaymentServiceServer(s, &server{})
	
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
