package wallet

import (
	"time"
)

type Wallet struct {
	ID        string    `json:"id" db:"id"`
	Balance   float64   `json:"balance" db:"balance"`
	Currency  string    `json:"currency" db:"currency"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
