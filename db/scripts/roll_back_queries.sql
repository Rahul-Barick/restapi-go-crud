-- Rollback script for internal_transfers schema

-- Drop dependent tables first
DROP TABLE IF EXISTS ledger_entries;

-- Drop ENUM type
DO $$ BEGIN
    DROP TYPE IF EXISTS ledger_entry_type;
EXCEPTION
    WHEN undefined_object THEN null;
END $$;

-- Drop transactions table
DROP TABLE IF EXISTS transactions;

-- Drop accounts table
DROP TABLE IF EXISTS accounts;

-- Optionally remove UUID extension
DROP EXTENSION IF EXISTS "uuid-ossp";
