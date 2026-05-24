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
	Status       string    `xorm:"status" json:"status"`
	CreateTime   time.Time `xorm:"created 'create_time'" json:"-"`
	UpdateTime   time.Time `xorm:"updated 'update_time'" json:"-"`
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

type ListMstBankStatementFileParam struct {
	BankCode string `schema:"-" validate:"required"`
	CommonRequestParam
}

type ListMstBankStatementFileQuery struct {
	BankCode string
	Limit    int
	Offset   int
}

type MstBankStatementFileResponse struct {
	ID        string `json:"id"`
	BankCode  string `json:"bankCode"`
	FileName  string `json:"fileName"`
	StartDate int64  `json:"startDate"`
	EndDate   int64  `json:"endDate"`
	Status    string `json:"status"`
}

func ToMstBankStatementFileResponse(f *MstBankStatementFile) MstBankStatementFileResponse {
	return MstBankStatementFileResponse{
		ID:        f.ID,
		BankCode:  f.BankCode,
		FileName:  f.FileName,
		StartDate: f.StartDate,
		EndDate:   f.EndDate,
		Status:    f.Status,
	}
}
