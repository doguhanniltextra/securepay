# Architecture — SecurePay

## Overview
Zero-Trust Fintech Payment Platform.
Authentication via SPIFFE/SPIRE,
Distributed tracing with OpenTelemetry,
Event-driven architecture with Kafka.

## Services and Ports
| Service              | Language             | Port  | Role                                       |
|----------------------|----------------------|-------|--------------------------------------------|
| API Gateway          | Go 1.24              | 8080  | Routing, rate limiting, JWT verification   |
| Payment Service      | Go 1.24              | 8081  | Payment initiation, validation, state machine |
| Account Service      | Go 1.24              | 50051 | Balance management, gRPC server            |
| Notification Service | Java 21/Spring Boot  | 8083  | Kafka consumer, notification logging       |
| SPIRE Server         | SPIRE v1.x           | 8081  | SVID distribution                          |
| Kafka                | Apache Kafka         | 9092  | Event bus                                  |
| PostgreSQL           | PostgreSQL 16        | 5432  | Persistent data                            |
| Redis                | Redis 7              | 6379  | Cache + idempotency                        |
| Jaeger               | Jaeger v2            | 16686 | Trace visualization                        |
| Prometheus           | Latest               | 9090  | Metrics                                    |
| Grafana              | Latest               | 3000  | Dashboard                                  |

## SPIFFE ID Scheme
- spiffe://securepay.dev/api-gateway
- spiffe://securepay.dev/payment-service
- spiffe://securepay.dev/account-service
- spiffe://securepay.dev/notification-service

## Layer Structure
Layer 1 — External Access:
  Client → API Gateway (HTTP, JWT)
  No microservice is exposed to the outside

Layer 2 — Inter-Service (Zero-Trust):
  Payment Service ↔ Account Service (gRPC + mTLS via SPIFFE)
  NO static credentials, SVID-based identity

Layer 3 — Async Event Bus:
  Payment Service → Kafka → Account Service
  Payment Service → Kafka → Notification Service
  Separate consumer groups: account-service-group, notification-service-group

## Payment Flow (8 steps)
1. Client → API Gateway (HTTP POST /payments, JWT)
2. API Gateway → Payment Service (gRPC, mTLS/SPIFFE)
3. Payment Service → Account Service (gRPC, balance check)
4. Account Service → Redis (cache lookup)
5. Account Service → PostgreSQL (if cache miss)
6. Payment Service → Kafka (payment.initiated event)
7. Kafka → Account Service (update balance)
8. Kafka → Notification Service (log notification)

## Database
- payments schema: transactions table
- accounts schema: balances table
- Same PostgreSQL instance, different schemas (operational simplicity)

## Out of Scope
- Kafka transport security (to be documented in README)
- SPIRE HA setup (single node sufficient)
- OpenTelemetry Collector (services send directly to Jaeger)
