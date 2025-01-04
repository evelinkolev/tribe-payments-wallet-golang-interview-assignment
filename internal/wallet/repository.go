package wallet

import (
	"context"
	"database/sql"
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

func (r *repository) Create(ctx context.Context, currency string) (*Wallet, error) {
	wallet := &Wallet{
		ID:        uuid.New().String(),
		Currency:  currency,
		Balance:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO wallets (id, balance, currency, created_at, updated_at) 
         VALUES ($1, $2, $3, $4, $5)`,
		wallet.ID, wallet.Balance, wallet.Currency, wallet.CreatedAt, wallet.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (r *repository) Get(ctx context.Context, id string) (*Wallet, error) {
	wallet := &Wallet{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, balance, currency, created_at, updated_at 
         FROM wallets WHERE id = $1`,
		id).Scan(&wallet.ID, &wallet.Balance, &wallet.Currency,
		&wallet.CreatedAt, &wallet.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrWalletNotFound
	}

	return wallet, err
}

func (r *repository) UpdateBalance(ctx context.Context, id string, amount float64) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE wallets 
         SET balance = balance + $1, updated_at = $2 
         WHERE id = $3`,
		amount, time.Now(), id)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrWalletNotFound
	}

	return nil
}
