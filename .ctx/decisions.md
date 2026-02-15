# Architectural Decisions — SecurePay

## [2026-02-15] Polyglot architecture: Go + Java

**Reason:** Language selection based on service responsibility.
**Choice:**
- Go: API Gateway, Payment Service, Account Service (critical services)
- Java/Spring Boot: Notification Service (simple IO-bound service)
**Outcome:**
- Go-heavy profile highlighted in CV
- Java is not completely excluded, strengthening the polyglot narrative
- In interviews, can explain: "Go for critical services, Spring Boot for simple IO services"

## [2026-02-15] Service communication: gRPC + mTLS

**Reason:** gRPC integrates cleaner with SPIFFE mTLS.
**Choice:** gRPC (instead of REST)
**Outcome:** gRPC line in CV becomes concrete. mTLS provides automatic SVID rotation.

## [2026-02-15] SVID management: SPIFFE Go SDK

**Reason:** Manual X.509 loading requires extra code for rotation.
**Choice:** SPIFFE Go SDK
**Outcome:** SDK provides automatic rotation, code remains simple.

## [2026-02-15] Trace backend: Jaeger

**Reason:** OTLP native support and easy setup with all-in-one image.
**Choice:** Jaeger (instead of Grafana Tempo)
**Outcome:** docker run sufficient for single container.

## [2026-02-15] Kafka security: Out of Scope

**Reason:** Does not fit into the 10-11 day timeframe.
**Choice:** No Kafka mTLS/SASL in this iteration.
**Outcome:** Document in README: "Kafka transport security is out of scope for this iteration."

## [2026-02-15] DB isolation: Separate schemas

**Reason:** Separate instances create operational complexity.
**Choice:** Same PostgreSQL, different schemas (payments, accounts)
**Outcome:** Separate instances preferred in production, but distinct schemas sufficient for dev.

## [2026-02-15] Cache strategy: Read-aside

**Reason:** Balance heavy on reads, write-through not necessary.
**Choice:** Read-aside cache (Redis)
**Outcome:** If cache miss, fetch from PostgreSQL, write to Redis. TTL: 60s.

## [2026-02-15] OpenTelemetry Collector: None

**Reason:** Exceeds project scope.
**Choice:** Services send OTLP directly to Jaeger.
**Outcome:** Less infrastructure, simpler setup.

## [2026-02-15] Hard-Coded Endpoints: None

**Reason:** We want the project to be organized.
**Choice:** Open `endpoints.go` document and register there.
**Outcome:** More complex endpoint structure.