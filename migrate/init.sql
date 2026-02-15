-- Create schemas
CREATE SCHEMA IF NOT EXISTS payments;
CREATE SCHEMA IF NOT EXISTS accounts;

-- Create payments.transactions table
CREATE TABLE IF NOT EXISTS payments.transactions (
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

-- Create accounts.balances table
CREATE TABLE IF NOT EXISTS accounts.balances (
    account_id  UUID PRIMARY KEY,
    balance     NUMERIC(18,2) NOT NULL DEFAULT 0,
    currency    VARCHAR(3) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version     INT NOT NULL DEFAULT 1
);

-- Seed data for accounts
INSERT INTO accounts.balances (account_id, balance, currency) VALUES
('11111111-1111-1111-1111-111111111111', 1000.00, 'TRY'),
('22222222-2222-2222-2222-222222222222', 500.00, 'TRY')
ON CONFLICT (account_id) DO NOTHING;
