package model

import "time"

const MstBankStatementFileTableName = "amt_mst_bank_statement_file"

type MstBankStatementFile struct {
	ID           string    `xorm:"pk 'id'" json:"id"`
	BankCode     string    `xorm:"bank_code" json:"bankCode"`
	FileName     string    `xorm:"file_name" json:"fileName"`
	FileLocation string    `xorm:"file_location" json:"fileLocation"`
	StartDate    int64     `xorm:"start_date" json:"startDate"`
	EndDate      int64     `xorm:"end_date" json:"endDate"`
	HashID       string    `xorm:"hash_id" json:"-"`
	CreateTime   time.Time `xorm:"created 'create_time'" json:"-"`
}

func (MstBankStatementFile) TableName() string {
	return MstBankStatementFileTableName
}

type UploadMstBankStatementFileRequest struct {
	BankCode string
	StartDate int64
	EndDate   int64
	FileName  string
	FileContent []byte
}

type UploadMstBankStatementFileResponse struct {
	FileID string `json:"fileId"`
}
