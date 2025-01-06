package wallet

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, currency string) (*Wallet, error)
	Get(ctx context.Context, id string) (*Wallet, error)
	UpdateBalance(ctx context.Context, id string, amount float64) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func generateID() string {
	return uuid.New().String()
}

func (r *repository) Create(ctx context.Context, currency string) (*Wallet, error) {
	if currency == "" {
		return nil, errors.New("currency cannot be empty")
	}

	wallet := &Wallet{
		ID:        generateID(), // Ensure a unique ID is generated
		Currency:  currency,
		Balance:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `INSERT INTO wallets (id, currency, balance, created_at, updated_at) 
              VALUES (@id, @currency, @balance, @created_at, @updated_at)`

	_, err := r.db.ExecContext(ctx, query,
		sql.Named("id", wallet.ID),
		sql.Named("currency", wallet.Currency),
		sql.Named("balance", wallet.Balance),
		sql.Named("created_at", wallet.CreatedAt),
		sql.Named("updated_at", wallet.UpdatedAt),
	)
	if err != nil {
		return nil, errors.New("failed to insert wallet into database: " + err.Error())
	}

	return wallet, nil
}

func (r *repository) Get(ctx context.Context, id string) (*Wallet, error) {
	if id == "" {
		return nil, errors.New("wallet ID cannot be empty")
	}

	wallet := &Wallet{}
	query := `SELECT id, balance, currency, created_at, updated_at 
              FROM wallets WHERE id = @id`

	err := r.db.QueryRowContext(ctx, query,
		sql.Named("id", id),
	).Scan(&wallet.ID, &wallet.Balance, &wallet.Currency,
		&wallet.CreatedAt, &wallet.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWalletNotFound
		}
		return nil, errors.New("failed to retrieve wallet: " + err.Error())
	}

	return wallet, nil
}

func (r *repository) UpdateBalance(ctx context.Context, id string, amount float64) error {
	if id == "" {
		return errors.New("wallet ID cannot be empty")
	}
	if amount == 0 {
		return errors.New("amount must be non-zero")
	}

	query := `UPDATE wallets 
              SET balance = balance + @amount, updated_at = @updated_at 
              WHERE id = @id`

	result, err := r.db.ExecContext(ctx, query,
		sql.Named("amount", amount),
		sql.Named("updated_at", time.Now()),
		sql.Named("id", id),
	)
	if err != nil {
		return errors.New("failed to update wallet balance: " + err.Error())
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.New("failed to check affected rows: " + err.Error())
	}

	if rows == 0 {
		return ErrWalletNotFound
	}

	return nil
}
