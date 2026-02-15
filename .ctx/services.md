# SecurePay — Service Architecture

## Payment Service

**Responsibility:** Owner and orchestrator of payment process.

### File Structure

```
payment-service/
  main.go          → start gRPC server, port 8081
  handler.go       → implement InitiatePayment, GetPayment
  validator.go     → request validation
  repository.go    → PostgreSQL operations
  kafka.go         → publish event
  spiffe.go        → mTLS server credential
  telemetry.go     → OpenTelemetry setup
  state.go         → state machine
```

### What It Does

**When InitiatePayment arrives:**

```
1. Validation
   - amount > 0?
   - from_account and to_account in UUID format?
   - currency valid? (TRY, USD, EUR)
   - from_account == to_account? (self-transfer forbidden)

2. Idempotency check
   - Exists idempotency:{key} in Redis?
   - If yes: return saved response, do not repeat operation

3. Balance check
   - gRPC to account-service: CheckBalance(from_account)
   - If returned balance < amount: return FAILED

4. Write to database
   - Save to payments.transactions table as PENDING
   - version: 1 (for optimistic locking)

5. Write event to Kafka
   - topic: payment.initiated
   - payload: payment_id, from_account, to_account, amount, currency, timestamp

6. Write idempotency record to Redis
   - TTL: 86400 seconds (24 hours)

7. Return response
   - payment_id, status: PENDING, message: "Payment initiated"
```

**When GetPayment arrives:**

```
1. Fetch from payments.transactions table with payment_id
2. If not found: return gRPC NotFound error
3. If found: return all fields
```

### State Machine

```
PENDING → COMPLETED  (when account-service deducts balance)
PENDING → FAILED     (insufficient balance or operation error)
```

### PostgreSQL Table

```sql
CREATE TABLE payments.transactions (
    id              UUID PRIMARY KEY,
    from_account    UUID NOT NULL,
    to_account      UUID NOT NULL,
    amount          NUMERIC(18,2) NOT NULL,
    currency        VARCHAR(3) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    idempotency_key VARCHAR(255) UNIQUE NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version         INT NOT NULL DEFAULT 1
);
```

---

## Account Service

**Responsibility:** Balance store. Only reads and updates balance, unaware of payment logic.

### File Structure

```
account-service/
  main.go          → start gRPC server, port 50051
  handler.go       → implement CheckBalance
  repository.go    → PostgreSQL operations
  cache.go         → Redis cache
  kafka.go         → Kafka consumer (update balance)
  spiffe.go        → mTLS server credential
  telemetry.go     → OpenTelemetry setup
```

### What It Does

**When CheckBalance arrives:**

```
1. Check Redis
   - Exists balance:{account_id} key?
   - If yes: return from Redis (TTL: 60 seconds)

2. If cache miss, go to PostgreSQL
   - Fetch from accounts.balances table with account_id
   - If not found: return gRPC NotFound error

3. Write to Redis
   - balance:{account_id} → balance value
   - TTL: 60 seconds

4. Return response
   - account_id, balance, currency
```

**When payment.initiated arrives from Kafka:**

```
1. Parse Event
   - payment_id, from_account, to_account, amount, currency

2. Start PostgreSQL transaction
   - Read from_account balance (lock with FOR UPDATE)
   - Check if Balance >= amount
   - Deduct from_account balance: balance - amount
   - Increase to_account balance: balance + amount
   - Increment version (optimistic locking)
   - Commit Transaction

3. Clear Redis cache
   - delete balance:{from_account}
   - delete balance:{to_account}
   (so next CheckBalance fetches fresh data from DB)
```

### Seed Data

```go
// inside main.go, before starting gRPC server
accounts := []Account{
    {
        ID:       "11111111-1111-1111-1111-111111111111",
        Balance:  1000.00,
        Currency: "TRY",
    },
    {
        ID:       "22222222-2222-2222-2222-222222222222",
        Balance:  500.00,
        Currency: "TRY",
    },
}
// INSERT ... ON CONFLICT DO NOTHING
```

### PostgreSQL Table

```sql
CREATE TABLE accounts.balances (
    account_id  UUID PRIMARY KEY,
    balance     NUMERIC(18,2) NOT NULL DEFAULT 0,
    currency    VARCHAR(3) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version     INT NOT NULL DEFAULT 1
);
```

---

## Notification Service

**Responsibility:** Receives event from Kafka, records notification. Logging is sufficient for this project.

### File Structure

```
notification-service/
  src/main/java/com/securepay/notification/
    consumer/
      PaymentEventConsumer.java  → Kafka listener
    model/
      PaymentEvent.java          → event model
    config/
      KafkaConfig.java           → Kafka configuration
    telemetry/
      TelemetryConfig.java       → OpenTelemetry setup
```

### What It Does

**When event arrives from payment.initiated topic:**

```
1. Deserialize Event
   - Convert to PaymentEvent object

2. Log
   - payment_id
   - from_account → to_account
   - amount + currency
   - timestamp

3. Operation finished
   (in real system: push notification, SMS, email sent)
```

### Code

```java
@KafkaListener(
    topics = "payment.initiated",
    groupId = "notification-service-group"
)
public void handlePaymentEvent(PaymentEvent event) {
    log.info("Payment notification received: " +
        "payment_id={}, from={}, to={}, amount={} {}",
        event.getPaymentId(),
        event.getFromAccount(),
        event.getToAccount(),
        event.getAmount(),
        event.getCurrency()
    );
}
```

---

## Three Services Working Together

```
Client: Ahmet → Mehmet, 150 TL

1.  payment-service: request arrived, validation complete
2.  payment-service: ask account-service → balance 1000 TL, sufficient
3.  payment-service: write PENDING to DB
4.  payment-service: write event to Kafka
5.  payment-service: return PENDING to client

Background:
6.  account-service: received event from Kafka
7.  account-service: made Ahmet 850 TL, Mehmet 650 TL
8.  account-service: cleared Redis cache
9.  notification-service: received event from Kafka
10. notification-service: logged "150 TL transferred"
```

---

## Summary Table

| Service              | gRPC Server                    | Kafka Producer     | Kafka Consumer     | PostgreSQL   | Redis        |
|----------------------|--------------------------------|--------------------|--------------------|--------------|--------------|
| payment-service      | InitiatePayment, GetPayment    | payment.initiated  | No                 | transactions | idempotency  |
| account-service      | CheckBalance                   | No                 | payment.initiated  | balances     | balance cache|
| notification-service | No                             | No                 | payment.initiated  | No           | No           |