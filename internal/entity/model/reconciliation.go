package model

import "time"

const TrxReconciliationTableName = "amt_trx_reconciliation"

type TrxReconciliation struct {
	ID                   string    `xorm:"pk 'id'" json:"id"`
	BankCode             string    `xorm:"bank_code" json:"bankCode"`
	BankStatementFileID  string    `xorm:"id_bank_statement_file" json:"bankStatementFileId"`
	BankStatementID      *string   `xorm:"id_bank_statement null" json:"bankStatementId,omitempty"`
	TransactionID        *string   `xorm:"id_transaction null" json:"transactionId,omitempty"`
	Status               string    `xorm:"status" json:"status"`
	MatchMethod          string    `xorm:"match_method" json:"matchMethod"`
	CreateTime           time.Time `xorm:"created 'create_time'" json:"-"`
	UpdateTime           time.Time `xorm:"updated 'update_time'" json:"-"`
}

func (TrxReconciliation) TableName() string {
	return TrxReconciliationTableName
}

type RunReconciliationResponse struct {
	FileID        string `json:"fileId"`
	Matched       int    `json:"matched"`
	Mismatch      int    `json:"mismatch"`
	OnlyInSystem  int    `json:"onlyInSystem"`
	OnlyInBank    int    `json:"onlyInBank"`
	Ambiguous     int    `json:"ambiguous"`
}

type ListReconciliationParam struct {
	BankCode string `schema:"-" validate:"required"`
	FileID   string `schema:"fileId" validate:"required"`
	Status   string `schema:"status" validate:"omitempty,oneof=MATCHED MISMATCH ONLY_IN_SYSTEM ONLY_IN_BANK AMBIGUOUS"`
	CommonRequestParam
}

type ReconciliationResponse struct {
	ID                  string `json:"id"`
	BankCode            string `json:"bankCode"`
	BankStatementFileID string `json:"bankStatementFileId"`
	BankStatementID     string `json:"bankStatementId,omitempty"`
	TransactionID       string `json:"transactionId,omitempty"`
	Status              string `json:"status"`
	MatchMethod         string `json:"matchMethod"`
}

type ReconciliationSummaryResponse struct {
	FileID                 string                        `json:"fileId"`
	BankCode               string                        `json:"bankCode"`
	StartDate              int64                         `json:"startDate"`
	EndDate                int64                         `json:"endDate"`
	TotalProcessed         int                           `json:"totalProcessed"`
	TotalMatched           int                           `json:"totalMatched"`
	TotalUnmatched         int                           `json:"totalUnmatched"`
	TotalMismatch          int                           `json:"totalMismatch"`
	TotalAmbiguous         int                           `json:"totalAmbiguous"`
	TotalDiscrepancyAmount string                        `json:"totalDiscrepancyAmount"`
	UnmatchedByBank        map[string]UnmatchedByBankGroup `json:"unmatchedByBank"`
}

type UnmatchedByBankGroup struct {
	OnlyInSystem []MstTransactionResponse      `json:"onlyInSystem"`
	OnlyInBank   []BankStatementSummaryDetail  `json:"onlyInBank"`
}

type BankStatementSummaryDetail struct {
	ID               string `json:"id"`
	UniqueIdentifier string `json:"uniqueIdentifier"`
	Amount           string `json:"amount"`
	Date             int64  `json:"date"`
}

func ToReconciliationResponse(r *TrxReconciliation) ReconciliationResponse {
	resp := ReconciliationResponse{
		ID:                  r.ID,
		BankCode:            r.BankCode,
		BankStatementFileID: r.BankStatementFileID,
		Status:              r.Status,
		MatchMethod:         r.MatchMethod,
	}
	if r.BankStatementID != nil {
		resp.BankStatementID = *r.BankStatementID
	}
	if r.TransactionID != nil {
		resp.TransactionID = *r.TransactionID
	}
	return resp
}
