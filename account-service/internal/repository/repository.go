package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/shopspring/decimal"

	"securepay/account-service/models"

	"go.opentelemetry.io/otel"
	_ "github.com/lib/pq"
)

// Repository defines the interface for database operations
type Repository interface {
	GetAccount(ctx context.Context, accountID string) (*models.Account, error)
	UpsertAccount(ctx context.Context, account *models.Account) error
	ProcessPayment(ctx context.Context, fromAccountID, toAccountID string, amount decimal.Decimal) error
}

// PostgresRepository implements Repository
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgresRepository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// GetAccount fetches account details by ID
func (r *PostgresRepository) GetAccount(ctx context.Context, accountID string) (*models.Account, error) {
	ctx, span := otel.Tracer("account-service").Start(ctx, "postgres.GetAccount")
	defer span.End()

	query := `
		SELECT account_id, balance, currency, created_at, updated_at, version
		FROM accounts.balances
		WHERE account_id = $1
	`

	var acc models.Account
	err := r.db.QueryRowContext(ctx, query, accountID).Scan(
		&acc.ID,
		&acc.Balance,
		&acc.Currency,
		&acc.CreatedAt,
		&acc.UpdatedAt,
		&acc.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("account not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &acc, nil
}

// UpsertAccount inserts or updates an account (used for seeding)
func (r *PostgresRepository) UpsertAccount(ctx context.Context, account *models.Account) error {
	ctx, span := otel.Tracer("account-service").Start(ctx, "postgres.UpsertAccount")
	defer span.End()

	query := `
		INSERT INTO accounts.balances (account_id, balance, currency, created_at, updated_at, version)
		VALUES ($1, $2, $3, NOW(), NOW(), 1)
		ON CONFLICT (account_id) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, query, account.ID, account.Balance, account.Currency)
	if err != nil {
		return fmt.Errorf("failed to upsert account: %w", err)
	}
	return nil
}


// ProcessPayment handles the transactional balance update
func (r *PostgresRepository) ProcessPayment(ctx context.Context, fromAccountID, toAccountID string, amount decimal.Decimal) error {
	ctx, span := otel.Tracer("account-service").Start(ctx, "postgres.ProcessPayment")
	defer span.End()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// 1. Lock From Account and Check Balance (using decimal for precise comparison)
	var fromBalance decimal.Decimal
	err = tx.QueryRowContext(ctx, "SELECT balance FROM accounts.balances WHERE account_id = $1 FOR UPDATE", fromAccountID).Scan(&fromBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("from_account not found")
		}
		return fmt.Errorf("failed to lock from_account: %w", err)
	}

	if fromBalance.LessThan(amount) {
		return fmt.Errorf("insufficient funds")
	}

	// 2. Deduct from From Account (SQL NUMERIC arithmetic is exact)
	_, err = tx.ExecContext(ctx, "UPDATE accounts.balances SET balance = balance - $1, version = version + 1, updated_at = NOW() WHERE account_id = $2", amount, fromAccountID)
	if err != nil {
		return fmt.Errorf("failed to deduct balance: %w", err)
	}

	// 3. Add to To Account
	res, err := tx.ExecContext(ctx, "UPDATE accounts.balances SET balance = balance + $1, version = version + 1, updated_at = NOW() WHERE account_id = $2", amount, toAccountID)
	if err != nil {
		return fmt.Errorf("failed to credit balance: %w", err)
	}
	
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("to_account not found")
	}

	slog.InfoContext(ctx, "Successfully processed payment", "from", fromAccountID, "to", toAccountID, "amount", amount.String())
	return nil
}

