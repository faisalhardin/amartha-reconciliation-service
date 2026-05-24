CREATE TABLE IF NOT EXISTS amt_mst_bank_statement_file (
    id UUID PRIMARY KEY,
    bank_code VARCHAR(50) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_location TEXT NOT NULL,
    start_date BIGINT NOT NULL,
    end_date BIGINT NOT NULL,
    hash_id VARCHAR(255) NOT NULL DEFAULT '',
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_amt_mst_bank_statement_file_bank_code
    ON amt_mst_bank_statement_file (bank_code);
