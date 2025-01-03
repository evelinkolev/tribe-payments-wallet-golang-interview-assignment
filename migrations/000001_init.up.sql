--- Your forward (up) migrations go here
CREATE TABLE wallets (
    id VARCHAR(36) PRIMARY KEY,
    balance DECIMAL(20,2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(3) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
