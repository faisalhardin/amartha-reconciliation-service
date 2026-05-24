package model

import (
	"time"

	"github.com/shopspring/decimal"
)

const MstBankStatementTableName = "amt_mst_bank_statement"

type MstBankStatement struct {
	ID                   string          `xorm:"pk 'id'" json:"id"`
	BankCode             string          `xorm:"bank_code" json:"bankCode"`
	BankStatementFileID  string          `xorm:"id_bank_statement_file" json:"bankStatementFileId"`
	UniqueIdentifier     string          `xorm:"unique_identifier" json:"uniqueIdentifier"`
	Amount           decimal.Decimal `xorm:"amount" json:"amount"`
	Date             int64           `xorm:"date" json:"date"`
	CreateTime       time.Time       `xorm:"created 'create_time'" json:"-"`
	UpdateTime       time.Time       `xorm:"updated 'update_time'" json:"-"`
}

func (MstBankStatement) TableName() string {
	return MstBankStatementTableName
}

type ListMstBankStatementParam struct {
	BankCode string `schema:"-" validate:"required"`
	FileID   string `schema:"fileId" validate:"required"`
	CommonRequestParam
}

type ListMstBankStatementQuery struct {
	BankCode string
	FileID   string
	Limit    int
	Offset   int
}

type MstBankStatementResponse struct {
	ID                  string          `json:"id"`
	BankCode            string          `json:"bankCode"`
	BankStatementFileID string          `json:"bankStatementFileId"`
	UniqueIdentifier    string          `json:"uniqueIdentifier"`
	Amount              decimal.Decimal `json:"amount"`
	Date                int64           `json:"date"`
}

func ToMstBankStatementResponse(s *MstBankStatement) MstBankStatementResponse {
	return MstBankStatementResponse{
		ID:                  s.ID,
		BankCode:            s.BankCode,
		BankStatementFileID: s.BankStatementFileID,
		UniqueIdentifier:    s.UniqueIdentifier,
		Amount:              s.Amount,
		Date:                s.Date,
	}
}
