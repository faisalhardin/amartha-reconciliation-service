CREATE TABLE IF NOT EXISTS amt_mst_transaction (
    id UUID PRIMARY KEY,
    amount NUMERIC(19, 4) NOT NULL,
    type VARCHAR(10) NOT NULL,
    transaction_time BIGINT NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS amt_mst_bank_statement (
    id UUID PRIMARY KEY,
    bank_code VARCHAR(50) NOT NULL,
    unique_identifier VARCHAR(255) NOT NULL,
    amount NUMERIC(19, 4) NOT NULL,
    date BIGINT NOT NULL,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_amt_mst_bank_statement_unique_identifier
    ON amt_mst_bank_statement (unique_identifier);

ALTER TABLE amt_mst_bank_statement
    ADD COLUMN IF NOT EXISTS bank_code VARCHAR(50) NOT NULL DEFAULT '';
