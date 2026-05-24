ALTER TABLE amt_mst_bank_statement_file
    ADD COLUMN IF NOT EXISTS update_time TIMESTAMPTZ NOT NULL DEFAULT NOW();
