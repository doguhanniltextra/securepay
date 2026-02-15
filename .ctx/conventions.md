# Kod Kuralları — SecurePay

## Go Servisleri (api-gateway, payment-service, account-service)

### Dosya İsimlendirme
- main.go          → entry point
- router.go        → HTTP/gRPC routing
- handler.go       → request handler'lar
- validator.go     → validasyon logic
- repository.go    → DB işlemleri
- cache.go         → Redis işlemleri
- kafka.go         → Kafka producer/consumer
- spiffe.go        → SPIFFE/SPIRE entegrasyonu
- telemetry.go     → OpenTelemetry setup
- state.go         → state machine

### Hata Yönetimi
- Her fonksiyon (result, error) döner
- Hata wrap: fmt.Errorf("context: %w", err)
- panic YOK, her hata handle edilmeli
- gRPC hataları: status.Errorf(codes.X, "mesaj")

### Logging
- log/slog paketi kullan (Go 1.21+)
- console.log değil: slog.Info(), slog.Error()
- Her log'da trace_id field'ı olmalı

### OpenTelemetry
- Her servis telemetry.go dosyası içerir
- OTLP exporter → Jaeger (localhost:4317)
- W3C TraceContext propagation
- Span isimleri: "service.operation" formatında

## Java Servisi (notification-service)

### Paket Yapısı
com.securepay.notification/
  consumer/     → Kafka consumer
  model/        → Data model'lar
  config/       → Spring konfigürasyonu
  telemetry/    → OTel konfigürasyonu

### Bağımlılıklar (pom.xml)
- spring-kafka
- opentelemetry-spring-boot-starter
- opentelemetry-exporter-otlp

## Kafka
- Topic: payment.initiated
- Consumer groups: account-service-group, notification-service-group
- Event şeması:
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
- Her tablo: created_at, updated_at timestamp
- Her tablo: version int (optimistic locking için)
- UUID tip: uuid (PostgreSQL native)

## Redis Key Formatları
- Bakiye cache: balance:{account_id}
- Idempotency: idempotency:{key}
- TTL bakiye: 60 saniye
- TTL idempotency: 24 saat (86400 saniye)

## gRPC
- Proto dosyaları: proto/ klasöründe
- Go generated: proto/gen/go/
- Java generated: proto/gen/java/
- Service naming: XxxService
- Method naming: VerbNoun (CheckBalance, InitiatePayment)

## SPIFFE
- Trust domain: securepay.dev
- Socket path: /tmp/spire-agent/public/api.sock
- SVID TTL: 1 saat
- Go SDK: github.com/spiffe/go-spiffe/v2
- Java SDK: io.spiffe:java-spiffe-core
