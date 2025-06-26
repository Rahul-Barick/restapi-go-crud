-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =========================
-- ACCOUNTS TABLE
-- =========================
CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    account_id INT UNIQUE NOT NULL,
    balance NUMERIC(20,6) NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC')
);

-- =========================
-- TRANSACTIONS TABLE
-- =========================
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source_account_id INT NOT NULL,
    destination_account_id INT NOT NULL,
    amount NUMERIC(20,6) NOT NULL,
    reference_id UUID UNIQUE NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC'),
    
    -- Foreign keys
    CONSTRAINT fk_txn_source_account FOREIGN KEY (source_account_id) REFERENCES accounts(account_id) ON DELETE RESTRICT,
    CONSTRAINT fk_txn_dest_account FOREIGN KEY (destination_account_id) REFERENCES accounts(account_id) ON DELETE RESTRICT
);

-- =========================
-- ENUM TYPE FOR LEDGER
-- =========================
DO $$ BEGIN
    CREATE TYPE ledger_entry_type AS ENUM ('CREDIT', 'DEBIT');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- =========================
-- LEDGER_ENTRIES TABLE
-- =========================
CREATE TABLE IF NOT EXISTS ledger_entries (
    id SERIAL PRIMARY KEY,
    account_id INT NOT NULL,
    transaction_id UUID NOT NULL,
    amount NUMERIC(20,6) NOT NULL,
    type ledger_entry_type NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC'),

    -- Foreign keys
    CONSTRAINT fk_ledger_account FOREIGN KEY (account_id) REFERENCES accounts(account_id) ON DELETE RESTRICT,
    CONSTRAINT fk_ledger_transaction FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE RESTRICT
);

-- =========================
-- INDEXES
-- =========================
CREATE UNIQUE INDEX idx_account_account_id ON accounts(account_id);
CREATE UNIQUE INDEX idx_transactions_reference_id ON transactions(reference_id);
CREATE INDEX idx_transactions_source_account_id ON transactions(source_account_id);
CREATE INDEX idx_transactions_destination_account_id ON transactions(destination_account_id);
CREATE INDEX idx_ledger_account_id ON ledger_entries(account_id);
CREATE INDEX idx_ledger_transaction_id ON ledger_entries(transaction_id);
CREATE INDEX idx_ledger_type ON ledger_entries(type);
