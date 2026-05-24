CREATE TABLE IF NOT EXISTS amt_trx_reconciliation (
    id UUID PRIMARY KEY,
    bank_code VARCHAR(50) NOT NULL,
    id_bank_statement_file UUID NOT NULL,
    id_bank_statement UUID,
    id_transaction UUID,
    status VARCHAR(30) NOT NULL,
    match_method VARCHAR(20) NOT NULL DEFAULT 'NONE',
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_amt_trx_reconciliation_file
    ON amt_trx_reconciliation (id_bank_statement_file);

CREATE INDEX IF NOT EXISTS idx_amt_trx_reconciliation_file_status
    ON amt_trx_reconciliation (id_bank_statement_file, status);

CREATE INDEX IF NOT EXISTS idx_amt_trx_reconciliation_file_statement
    ON amt_trx_reconciliation (id_bank_statement_file, id_bank_statement);

CREATE INDEX IF NOT EXISTS idx_amt_trx_reconciliation_file_transaction
    ON amt_trx_reconciliation (id_bank_statement_file, id_transaction);
