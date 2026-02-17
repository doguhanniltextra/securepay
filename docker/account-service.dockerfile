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
# We need proto and account-service
COPY proto/ ./proto/
COPY account-service/ ./account-service/

# Build the application
WORKDIR /app/account-service
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/account-service .

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /app

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Copy the binary
COPY --from=builder /app/bin/account-service .

# Expose port
EXPOSE 8082

# Command to run
CMD ["./account-service"]
