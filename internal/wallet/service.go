package wallet

import (
	"context"
	"errors"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrWalletNotFound    = errors.New("wallet not found")
	ErrInvalidAmount     = errors.New("invalid amount")
)

type Service interface {
	CreateWallet(ctx context.Context, currency string) (*Wallet, error)
	GetWallet(ctx context.Context, id string) (*Wallet, error)
	Deposit(ctx context.Context, id string, amount float64) error
	Withdraw(ctx context.Context, id string, amount float64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateWallet(ctx context.Context, currency string) (*Wallet, error) {
	return s.repo.Create(ctx, currency)
}

func (s *service) GetWallet(ctx context.Context, id string) (*Wallet, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Deposit(ctx context.Context, id string, amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	return s.repo.UpdateBalance(ctx, id, amount)
}

func (s *service) Withdraw(ctx context.Context, id string, amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	wallet, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	if wallet.Balance < amount {
		return ErrInsufficientFunds
	}

	return s.repo.UpdateBalance(ctx, id, -amount)
}
