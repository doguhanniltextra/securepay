# Stage 1: Build
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy workspace configuration
COPY go.work go.work.sum ./

# Copy module definitions
COPY api-gateway/go.mod api-gateway/go.sum ./api-gateway/
COPY proto/go.mod proto/go.sum ./proto/
COPY payment-service/go.mod payment-service/go.sum ./payment-service/
COPY account-service/go.mod account-service/go.sum ./account-service/

# Copy sources for required modules
# We need proto for generated code and payment-service for main code.
COPY proto/ ./proto/
COPY payment-service/ ./payment-service/
# No need to copy api-gateway source, only go.mod to satisfy go.work (though go.work might complain if dir is empty? No, go.mod is enough usually)
# Actually, let's just make the directory if copying files might fail if they don't exist.
# But since we copied go.mod above, the directory exists.

# Build the application
WORKDIR /app/payment-service
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/payment-service .

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app
# Install CA certificates for secure communication (e.g. AWS or external APIs)
RUN apk --no-cache add ca-certificates
# Copy the binary
COPY --from=builder /app/bin/payment-service .
# Expose the application port
EXPOSE 8081
# Command to run
CMD ["./payment-service"]
