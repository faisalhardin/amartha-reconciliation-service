package reconciliation

import (
	"context"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/constant"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/library/common/commonerr"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func (u *reconciliationUC) GetSummary(ctx context.Context, fileID string) (*model.ReconciliationSummaryResponse, error) {
	file, err := u.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		if errors.Is(err, constant.ErrNotFound) {
			return nil, commonerr.Set404()
		}
		return nil, errors.Wrap(err, wrapErrMsg+"GetSummary.GetByID")
	}

	normalizedBankCode := normalizeBankCode(file.BankCode)

	if file.Status != constant.BankStatementFileStatusCompleted {
		return nil, commonerr.SetNewBadRequest("invalid_file_status", "bank statement file must be COMPLETED")
	}

	rows, err := u.reconRepo.ListAllByFileID(ctx, fileID)
	if err != nil {
		return nil, errors.Wrap(err, wrapErrMsg+"GetSummary.ListAllByFileID")
	}

	summary := &model.ReconciliationSummaryResponse{
		FileID:          fileID,
		BankCode:        normalizedBankCode,
		StartDate:       file.StartDate,
		EndDate:         file.EndDate,
		UnmatchedByBank: make(map[string]model.UnmatchedByBankGroup),
	}

	discrepancyTotal := decimal.Zero
	for i := range rows {
		row := &rows[i]
		incrementSummaryCounts(summary, row.Status)

		switch row.Status {
		case constant.ReconciliationStatusMismatch:
			diff, err := u.mismatchDiscrepancy(ctx, normalizedBankCode, row)
			if err != nil {
				return nil, errors.Wrap(err, wrapErrMsg+"GetSummary.mismatchDiscrepancy")
			}
			discrepancyTotal = discrepancyTotal.Add(diff)

		case constant.ReconciliationStatusOnlyInSystem:
			if row.TransactionID == nil {
				continue
			}
			tx, err := u.loadTransaction(ctx, normalizedBankCode, *row.TransactionID)
			if err != nil {
				return nil, err
			}
			appendOnlyInSystem(summary, normalizedBankCode, model.ToMstTransactionResponse(tx))

		case constant.ReconciliationStatusOnlyInBank:
			if row.BankStatementID == nil {
				continue
			}
			stmt, err := u.loadBankStatement(ctx, *row.BankStatementID)
			if err != nil {
				return nil, err
			}
			appendOnlyInBank(summary, normalizedBankCode, toBankStatementSummaryDetail(stmt))
		}
	}

	summary.TotalDiscrepancyAmount = discrepancyTotal.String()
	return summary, nil
}

func incrementSummaryCounts(s *model.ReconciliationSummaryResponse, status string) {
	s.TotalProcessed++
	switch status {
	case constant.ReconciliationStatusMatched:
		s.TotalMatched++
	case constant.ReconciliationStatusMismatch:
		s.TotalMismatch++
	case constant.ReconciliationStatusOnlyInSystem, constant.ReconciliationStatusOnlyInBank:
		s.TotalUnmatched++
	case constant.ReconciliationStatusAmbiguous:
		s.TotalAmbiguous++
	}
}

func appendOnlyInSystem(s *model.ReconciliationSummaryResponse, bankCode string, tx model.MstTransactionResponse) {
	group := s.UnmatchedByBank[bankCode]
	group.OnlyInSystem = append(group.OnlyInSystem, tx)
	s.UnmatchedByBank[bankCode] = group
}

func appendOnlyInBank(s *model.ReconciliationSummaryResponse, bankCode string, stmt model.BankStatementSummaryDetail) {
	group := s.UnmatchedByBank[bankCode]
	group.OnlyInBank = append(group.OnlyInBank, stmt)
	s.UnmatchedByBank[bankCode] = group
}

func toBankStatementSummaryDetail(stmt *model.MstBankStatement) model.BankStatementSummaryDetail {
	return model.BankStatementSummaryDetail{
		ID:               stmt.ID,
		UniqueIdentifier: stmt.UniqueIdentifier,
		Amount:           stmt.Amount.String(),
		Date:             stmt.Date,
	}
}

func (u *reconciliationUC) mismatchDiscrepancy(ctx context.Context, bankCode string, row *model.TrxReconciliation) (decimal.Decimal, error) {
	if row.BankStatementID == nil || row.TransactionID == nil {
		return decimal.Zero, nil
	}
	stmt, err := u.loadBankStatement(ctx, *row.BankStatementID)
	if err != nil {
		return decimal.Zero, err
	}
	tx, err := u.loadTransaction(ctx, bankCode, *row.TransactionID)
	if err != nil {
		return decimal.Zero, err
	}
	diff := stmt.Amount.Sub(signedAmount(*tx))
	return diff.Abs(), nil
}

func (u *reconciliationUC) loadTransaction(ctx context.Context, bankCode, id string) (*model.MstTransaction, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, commonerr.SetNewBadRequest("invalid_transaction_id", "invalid transaction id in reconciliation row")
	}
	tx, err := u.txRepo.GetByID(ctx, bankCode, parsed)
	if err != nil {
		if errors.Is(err, constant.ErrNotFound) {
			return nil, commonerr.Set404()
		}
		return nil, errors.Wrap(err, wrapErrMsg+"loadTransaction")
	}
	return tx, nil
}

func (u *reconciliationUC) loadBankStatement(ctx context.Context, id string) (*model.MstBankStatement, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, commonerr.SetNewBadRequest("invalid_bank_statement_id", "invalid bank statement id in reconciliation row")
	}
	stmt, err := u.stmtRepo.GetByID(ctx, parsed)
	if err != nil {
		if errors.Is(err, constant.ErrNotFound) {
			return nil, commonerr.Set404()
		}
		return nil, errors.Wrap(err, wrapErrMsg+"loadBankStatement")
	}
	return stmt, nil
}
