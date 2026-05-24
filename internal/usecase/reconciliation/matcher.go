package reconciliation

import (
	"strings"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/constant"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
	"github.com/shopspring/decimal"
)

type matchResult struct {
	BankStatementID string
	TransactionID   string
	Status          string
	MatchMethod     string
}

func normalizeBankCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func sameBankCode(tx model.MstTransaction, stmt model.MstBankStatement) bool {
	return normalizeBankCode(tx.BankCode) == normalizeBankCode(stmt.BankCode)
}

func signedAmount(tx model.MstTransaction) decimal.Decimal {
	if tx.Type == model.TransactionTypeDebit {
		return tx.Amount.Neg()
	}
	return tx.Amount
}

func matchStatements(statements []model.MstBankStatement, transactions []model.MstTransaction) []matchResult {
	stmtPool := append([]model.MstBankStatement(nil), statements...)
	txPool := append([]model.MstTransaction(nil), transactions...)
	var results []matchResult

	results = append(results, referencePass(&stmtPool, &txPool)...)
	results = append(results, compositePass(&stmtPool, &txPool)...)
	results = append(results, leftoverPass(stmtPool, txPool)...)

	return results
}

func referencePass(stmtPool *[]model.MstBankStatement, txPool *[]model.MstTransaction) []matchResult {
	var results []matchResult
	remainingStmt := (*stmtPool)[:0]
	usedTx := make(map[string]bool)

	for _, stmt := range *stmtPool {
		ref := strings.TrimSpace(stmt.UniqueIdentifier)
		if ref == "" {
			remainingStmt = append(remainingStmt, stmt)
			continue
		}

		var candidates []model.MstTransaction
		for _, tx := range *txPool {
			if usedTx[tx.ID] {
				continue
			}
			if !sameBankCode(tx, stmt) {
				continue
			}
			if strings.TrimSpace(tx.TransactionReference) == ref {
				candidates = append(candidates, tx)
			}
		}

		if len(candidates) == 0 {
			remainingStmt = append(remainingStmt, stmt)
			continue
		}

		if len(candidates) > 1 {
			results = append(results, matchResult{
				BankStatementID: stmt.ID,
				Status:          constant.ReconciliationStatusAmbiguous,
				MatchMethod:     constant.ReconciliationMatchMethodReference,
			})
			continue
		}

		tx := candidates[0]
		usedTx[tx.ID] = true
		status := constant.ReconciliationStatusMatched
		if !signedAmount(tx).Equal(stmt.Amount) {
			status = constant.ReconciliationStatusMismatch
		}
		results = append(results, matchResult{
			BankStatementID: stmt.ID,
			TransactionID:   tx.ID,
			Status:          status,
			MatchMethod:     constant.ReconciliationMatchMethodReference,
		})
	}

	*stmtPool = remainingStmt
	filteredTx := (*txPool)[:0]
	for _, tx := range *txPool {
		if !usedTx[tx.ID] {
			filteredTx = append(filteredTx, tx)
		}
	}
	*txPool = filteredTx
	return results
}

func compositePass(stmtPool *[]model.MstBankStatement, txPool *[]model.MstTransaction) []matchResult {
	var results []matchResult
	remainingStmt := (*stmtPool)[:0]
	usedTx := make(map[string]bool)

	for _, stmt := range *stmtPool {
		var candidates []model.MstTransaction
		for _, tx := range *txPool {
			if usedTx[tx.ID] {
				continue
			}
			if !sameBankCode(tx, stmt) {
				continue
			}
			if tx.TransactionTime != stmt.Date {
				continue
			}
			if !signedAmount(tx).Equal(stmt.Amount) {
				continue
			}
			candidates = append(candidates, tx)
		}

		if len(candidates) == 0 {
			remainingStmt = append(remainingStmt, stmt)
			continue
		}

		if len(candidates) > 1 {
			results = append(results, matchResult{
				BankStatementID: stmt.ID,
				Status:          constant.ReconciliationStatusAmbiguous,
				MatchMethod:     constant.ReconciliationMatchMethodComposite,
			})
			continue
		}

		tx := candidates[0]
		usedTx[tx.ID] = true
		results = append(results, matchResult{
			BankStatementID: stmt.ID,
			TransactionID:   tx.ID,
			Status:          constant.ReconciliationStatusMatched,
			MatchMethod:     constant.ReconciliationMatchMethodComposite,
		})
	}

	*stmtPool = remainingStmt
	filteredTx := (*txPool)[:0]
	for _, tx := range *txPool {
		if !usedTx[tx.ID] {
			filteredTx = append(filteredTx, tx)
		}
	}
	*txPool = filteredTx
	return results
}

func leftoverPass(stmtPool []model.MstBankStatement, txPool []model.MstTransaction) []matchResult {
	results := make([]matchResult, 0, len(stmtPool)+len(txPool))
	for _, stmt := range stmtPool {
		results = append(results, matchResult{
			BankStatementID: stmt.ID,
			Status:          constant.ReconciliationStatusOnlyInBank,
			MatchMethod:     constant.ReconciliationMatchMethodNone,
		})
	}
	for _, tx := range txPool {
		results = append(results, matchResult{
			TransactionID: tx.ID,
			Status:        constant.ReconciliationStatusOnlyInSystem,
			MatchMethod:   constant.ReconciliationMatchMethodNone,
		})
	}
	return results
}
