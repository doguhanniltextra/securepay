# Stage 1: Build
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy workspace configuration
COPY go.work go.work.sum ./

# Copy module definitions
COPY api-gateway/go.mod api-gateway/go.sum ./api-gateway/
COPY proto/go.mod proto/go.sum ./proto/
COPY payment-service/go.mod payment-service/go.sum ./payment-service/

# Download dependencies (this might fail if go.work.sum is missing or inconsistent, but we try)
# Alternatively, we just copy everything and build, relying on go mod tidy or vendor if present.
# Since we are in dev, let's copy sources first to ensure go.work resolves correctly.

COPY api-gateway/ ./api-gateway/
COPY proto/ ./proto/

# Build the application
WORKDIR /app/api-gateway
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api-gateway .

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /app

# Install CA certificates for secure communication
RUN apk --no-cache add ca-certificates

# Copy the binary
COPY --from=builder /app/bin/api-gateway .

# Expose the application port
EXPOSE 8080

# Command to run
CMD ["./api-gateway"]
