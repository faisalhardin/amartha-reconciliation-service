ALTER TABLE amt_mst_bank_statement_file
    ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'PENDING';

CREATE INDEX IF NOT EXISTS idx_amt_mst_bank_statement_file_status
    ON amt_mst_bank_statement_file (status);

ALTER TABLE amt_mst_bank_statement
    ADD COLUMN IF NOT EXISTS id_bank_statement_file UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000';

CREATE INDEX IF NOT EXISTS idx_amt_mst_bank_statement_id_bank_statement_file
    ON amt_mst_bank_statement (id_bank_statement_file);

DROP INDEX IF EXISTS idx_amt_mst_bank_statement_unique_identifier;
CREATE UNIQUE INDEX IF NOT EXISTS idx_amt_mst_bank_statement_file_unique_identifier
    ON amt_mst_bank_statement (id_bank_statement_file, unique_identifier);
