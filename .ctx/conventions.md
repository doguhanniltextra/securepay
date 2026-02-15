# Conventions — SecurePay

## Go Services (api-gateway, payment-service, account-service)

### File Naming
- main.go          → entry point
- router.go        → HTTP/gRPC routing
- handler.go       → request handlers
- validator.go     → validation logic
- repository.go    → DB operations
- cache.go         → Redis operations
- kafka.go         → Kafka producer/consumer
- spiffe.go        → SPIFFE/SPIRE integration
- telemetry.go     → OpenTelemetry setup
- state.go         → state machine

### Error Handling
- Every function returns (result, error)
- Error wrap: fmt.Errorf("context: %w", err)
- NO panic, every error must be handled
- gRPC errors: status.Errorf(codes.X, "message")

### Logging
- Use log/slog package (Go 1.21+)
- Not console.log: slog.Info(), slog.Error()
- Every log must include trace_id field

### OpenTelemetry
- Every service contains telemetry.go file
- OTLP exporter → Jaeger (localhost:4317)
- W3C TraceContext propagation
- Span names: "service.operation" format

## Java Service (notification-service)

### Package Structure
com.securepay.notification/
  consumer/     → Kafka consumer
  model/        → Data models
  config/       → Spring configuration
  telemetry/    → OTel configuration

### Dependencies (pom.xml)
- spring-kafka
- opentelemetry-spring-boot-starter
- opentelemetry-exporter-otlp

## Kafka
- Topic: payment.initiated
- Consumer groups: account-service-group, notification-service-group
- Event schema:
  {
    "payment_id": "uuid",
    "from_account": "uuid",
    "to_account": "uuid",
    "amount": 150.00,
    "currency": "TRY",
    "timestamp": "ISO8601"
  }

## PostgreSQL
- Schema: payments (Payment Service)
- Schema: accounts (Account Service)
- Each table: created_at, updated_at timestamp
- Each table: version int (for optimistic locking)
- UUID type: uuid (PostgreSQL native)

## Redis Key Formats
- Balance cache: balance:{account_id}
- Idempotency: idempotency:{key}
- TTL balance: 60 seconds
- TTL idempotency: 24 hours (86400 seconds)

## gRPC
- Proto files: in proto/ folder
- Go generated: proto/gen/go/
- Java generated: proto/gen/java/
- Service naming: XxxService
- Method naming: VerbNoun (CheckBalance, InitiatePayment)

## SPIFFE
- Trust domain: securepay.dev
- Socket path: /tmp/spire-agent/public/api.sock
- SVID TTL: 1 hour
- Go SDK: github.com/spiffe/go-spiffe/v2
- Java SDK: io.spiffe:java-spiffe-core