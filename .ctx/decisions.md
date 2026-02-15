# Mimari Kararlar — SecurePay

## [2026-02-15] Polyglot mimari: Go + Java

**Neden:** Her servisin sorumluluğuna göre dil seçimi yapıldı.
**Seçim:**
- Go: API Gateway, Payment Service, Account Service (kritik servisler)
- Java/Spring Boot: Notification Service (basit IO-bound servis)
**Sonuç:**
- CV'de Go ağırlıklı profil öne çıkıyor
- Java tamamen dışarıda kalmıyor, polyglot anlatısı güçleniyor
- Mülakatta "Kritik servisler için Go, basit IO servisleri için Spring Boot" açıklaması yapılabilir

## [2026-02-15] Servis iletişimi: gRPC + mTLS

**Neden:** SPIFFE mTLS ile gRPC daha temiz entegre oluyor.
**Seçim:** gRPC (REST alternatifi yerine)
**Sonuç:** CV'de gRPC satırı somutlaşıyor. mTLS otomatik SVID rotasyonu sağlıyor.

## [2026-02-15] SVID yönetimi: SPIFFE Go SDK

**Neden:** Manuel X.509 yüklemede rotation için ek kod gerekiyor.
**Seçim:** SPIFFE Go SDK
**Sonuç:** SDK otomatik rotasyon sağlıyor, kod basit kalıyor.

## [2026-02-15] Trace backend: Jaeger

**Neden:** OTLP native desteği ve all-in-one image ile kolay setup.
**Seçim:** Jaeger (Grafana Tempo alternatifi yerine)
**Sonuç:** docker run ile tek container yeterli.

## [2026-02-15] Kafka security: Kapsam dışı

**Neden:** 10-11 günlük süreye sığmıyor.
**Seçim:** Kafka mTLS/SASL bu iterasyonda yok.
**Sonuç:** README'de belgelenecek: "Kafka transport security is out of scope for this iteration."

## [2026-02-15] DB izolasyonu: Ayrı schema

**Neden:** Ayrı instance operasyonel kompleksite yaratır.
**Seçim:** Aynı PostgreSQL, farklı schema (payments, accounts)
**Sonuç:** Production'da ayrı instance tercih edilir ama dev için yeterli.

## [2026-02-15] Cache stratejisi: Read-aside

**Neden:** Bakiye okuma ağırlıklı, write-through gerekmez.
**Seçim:** Read-aside cache (Redis)
**Sonuç:** Cache miss ise PostgreSQL'den çek, Redis'e yaz. TTL: 60s.

## [2026-02-15] OpenTelemetry Collector: Yok

**Neden:** Projenin scope'unu aşıyor.
**Seçim:** Servisler doğrudan Jaeger'a OTLP gönderir.
**Sonuç:** Daha az altyapı, daha basit setup.

## [2026-02-15] Hard-Coded Endpoints: Yok

**Neden:** Projenin derli toplu olmasını istiyoruz.
**Seçim:** `endpoints.go` belgesi açılır ve oraya kaydedilir.
**Sonuç:** Daha karmaşık endpoint yapısı.