package main

import (
	"fmt"
	"log"
	"net/http"
	"securepay/api-gateway/middleware"
)

func main() {
	// Task 4.1: Modified main.go to include middleware and test routes.
	// This ensures JWT middleware is correctly applied and verified.

	port := 8080
	fmt.Printf("API Gateway starting on port %d...\n", port)

	mux := http.NewServeMux()

	// Public Health Check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	// Secure Endpoint (under /api/v1/)
	mux.HandleFunc("/api/v1/secure", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "SECURE DATA OK")
	})

	// Another Endpoint to verify non-/api/v1/ bypass (should be public)
	mux.HandleFunc("/public", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "PUBLIC DATA OK")
	})

	// Apply Middleware globally
	// The middleware itself handles path filtering (/api/v1/ vs others)
	handler := middleware.AuthMiddleware(mux)

	// Start server
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
