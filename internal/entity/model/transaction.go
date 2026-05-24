package model

import (
	"time"

	"github.com/shopspring/decimal"
)

const MstTransactionTableName = "amt_mst_transaction"

type TransactionType string

const (
	TransactionTypeDebit  TransactionType = "DEBIT"
	TransactionTypeCredit TransactionType = "CREDIT"
)

func (t TransactionType) IsValid() bool {
	switch t {
	case TransactionTypeDebit, TransactionTypeCredit:
		return true
	default:
		return false
	}
}

type MstTransaction struct {
	ID                    string          `xorm:"pk 'id'" json:"id"`
	BankCode              string          `xorm:"bank_code" json:"bankCode"`
	TransactionReference  string          `xorm:"transaction_reference" json:"transactionReference"`
	Amount                decimal.Decimal `xorm:"amount" json:"amount"`
	Type                  TransactionType `xorm:"type" json:"type"`
	TransactionTime       int64           `xorm:"transaction_time" json:"transactionTime"`
	CreateTime            time.Time       `xorm:"created 'create_time'" json:"-"`
	UpdateTime            time.Time       `xorm:"updated 'update_time'" json:"-"`
}

func (MstTransaction) TableName() string {
	return MstTransactionTableName
}

type CreateMstTransactionRequest struct {
	BankCode        string          `json:"-" validate:"required"`
	Amount          float64         `json:"amount" validate:"required,gt=0"`
	Type            TransactionType `json:"type" validate:"required,oneof=DEBIT CREDIT"`
	TransactionTime int64           `json:"transactionTime" validate:"required,gt=0"`
}

type ListMstTransactionParam struct {
	BankCode string `schema:"-" validate:"required"`
	CommonRequestParam
}

type MstTransactionResponse struct {
	ID                   string          `json:"id"`
	BankCode             string          `json:"bankCode"`
	TransactionReference string          `json:"transactionReference,omitempty"`
	Amount               decimal.Decimal `json:"amount"`
	Type                 TransactionType `json:"type"`
	TransactionTime      int64           `json:"transactionTime"`
}

type UploadMstTransactionCSVResponse struct {
	Inserted int `json:"inserted"`
}

func ToMstTransactionResponse(t *MstTransaction) MstTransactionResponse {
	return MstTransactionResponse{
		ID:                   t.ID,
		BankCode:             t.BankCode,
		TransactionReference: t.TransactionReference,
		Amount:               t.Amount,
		Type:                 t.Type,
		TransactionTime:      t.TransactionTime,
	}
}
