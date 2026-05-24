package model

import (
	"time"

	"github.com/shopspring/decimal"
)

const MstBankStatementTableName = "amt_mst_bank_statement"

type MstBankStatement struct {
	ID               string          `xorm:"pk 'id'" json:"id"`
	BankCode         string          `xorm:"bank_code" json:"bankCode"`
	UniqueIdentifier string          `xorm:"unique_identifier" json:"uniqueIdentifier"`
	Amount           decimal.Decimal `xorm:"amount" json:"amount"`
	Date             int64           `xorm:"date" json:"date"`
	CreateTime       time.Time       `xorm:"created 'create_time'" json:"-"`
	UpdateTime       time.Time       `xorm:"updated 'update_time'" json:"-"`
}

func (MstBankStatement) TableName() string {
	return MstBankStatementTableName
}
