# Mimari — SecurePay

## Genel Bakış
Zero-Trust Fintech Payment Platform.
SPIFFE/SPIRE ile servis kimlik doğrulama,
OpenTelemetry ile distributed tracing,
Kafka ile event-driven mimari.

## Servisler ve Portlar
| Servis               | Dil                  | Port  | Görev                                      |
|----------------------|----------------------|-------|--------------------------------------------|
| API Gateway          | Go 1.24              | 8080  | Routing, rate limiting, JWT doğrulama      |
| Payment Service      | Go 1.24              | 8081  | Ödeme başlatma, validasyon, state machine  |
| Account Service      | Go 1.24              | 50051 | Bakiye yönetimi, gRPC server               |
| Notification Service | Java 21/Spring Boot  | 8083  | Kafka consumer, bildirim kaydı             |
| SPIRE Server         | SPIRE v1.x           | 8081  | SVID dağıtımı                              |
| Kafka                | Apache Kafka         | 9092  | Event bus                                  |
| PostgreSQL           | PostgreSQL 16        | 5432  | Kalıcı veri                                |
| Redis                | Redis 7              | 6379  | Cache + idempotency                        |
| Jaeger               | Jaeger v2            | 16686 | Trace görselleştirme                       |
| Prometheus           | Latest               | 9090  | Metrics                                    |
| Grafana              | Latest               | 3000  | Dashboard                                  |

## SPIFFE ID Şeması
- spiffe://securepay.dev/api-gateway
- spiffe://securepay.dev/payment-service
- spiffe://securepay.dev/account-service
- spiffe://securepay.dev/notification-service

## Katman Yapısı
Katman 1 — Dış Erişim:
  Client → API Gateway (HTTP, JWT)
  Hiçbir mikroservis dışarıya açık değil

Katman 2 — Servisler Arası (Zero-Trust):
  Payment Service ↔ Account Service (gRPC + mTLS via SPIFFE)
  Statik credential YOK, SVID tabanlı kimlik

Katman 3 — Async Event Bus:
  Payment Service → Kafka → Account Service
  Payment Service → Kafka → Notification Service
  Ayrı consumer group'lar: account-service-group, notification-service-group

## Ödeme Akışı (8 adım)
1. Client → API Gateway (HTTP POST /payments, JWT)
2. API Gateway → Payment Service (gRPC, mTLS/SPIFFE)
3. Payment Service → Account Service (gRPC, bakiye kontrolü)
4. Account Service → Redis (cache lookup)
5. Account Service → PostgreSQL (cache miss ise)
6. Payment Service → Kafka (payment.initiated event)
7. Kafka → Account Service (bakiye güncelle)
8. Kafka → Notification Service (bildirim kaydı)

## Veritabanı
- payments schema: transactions tablosu
- accounts schema: balances tablosu
- Aynı PostgreSQL instance, farklı schema (operasyonel basitlik)

## Kapsam Dışı
- Kafka transport security (README'de belgelenecek)
- SPIRE HA kurulumu (single node yeterli)
- OpenTelemetry Collector (servisler doğrudan Jaeger'a gönderir)
