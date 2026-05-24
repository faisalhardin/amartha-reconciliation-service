package reconciliation

import (
	"context"
	"strings"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/constant"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
	bankstatementrepo "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/repo/bankstatement"
	bankstatementfilerepo "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/repo/bankstatementfile"
	reconciliationrepo "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/repo/reconciliation"
	transactionrepo "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/repo/transaction"
	reconciliationuc "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/usecase/reconciliation"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/library/common/commonerr"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	wrapErrMsg         = "ReconciliationUC."
	defaultListLimit   = 20
	maxListLimit       = 100
)

type reconciliationUC struct {
	fileRepo bankstatementfilerepo.BankStatementFileDB
	stmtRepo bankstatementrepo.BankStatementDB
	txRepo   transactionrepo.TransactionDB
	reconRepo reconciliationrepo.ReconciliationDB
}

func NewReconciliationUC(
	fileRepo bankstatementfilerepo.BankStatementFileDB,
	stmtRepo bankstatementrepo.BankStatementDB,
	txRepo transactionrepo.TransactionDB,
	reconRepo reconciliationrepo.ReconciliationDB,
) reconciliationuc.ReconciliationUC {
	return &reconciliationUC{
		fileRepo:  fileRepo,
		stmtRepo:  stmtRepo,
		txRepo:    txRepo,
		reconRepo: reconRepo,
	}
}

func (u *reconciliationUC) Run(ctx context.Context, bankCode, fileID string) (*model.RunReconciliationResponse, error) {
	file, err := u.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		if errors.Is(err, constant.ErrNotFound) {
			return nil, commonerr.Set404()
		}
		return nil, errors.Wrap(err, wrapErrMsg+"Run.GetByID")
	}

	if bankCode != "" {
		normalized, err := normalizeBankCodeParam(bankCode)
		if err != nil {
			return nil, err
		}
		if normalizeBankCode(file.BankCode) != normalized {
			return nil, commonerr.SetNewBadRequest("bank_code_mismatch", "file does not belong to the given bank code")
		}
	}

	if file.Status != constant.BankStatementFileStatusCompleted {
		return nil, commonerr.SetNewBadRequest("invalid_file_status", "bank statement file must be COMPLETED")
	}

	stmts, err := u.stmtRepo.ListByFileID(ctx, fileID)
	if err != nil {
		return nil, errors.Wrap(err, wrapErrMsg+"Run.ListByFileID")
	}

	fileBankCode := normalizeBankCode(file.BankCode)
	txs, err := u.txRepo.ListByBankCodeAndTimeRange(ctx, fileBankCode, file.StartDate, file.EndDate)
	if err != nil {
		return nil, errors.Wrap(err, wrapErrMsg+"Run.ListByBankCodeAndTimeRange")
	}

	matches := matchStatements(stmts, txs)
	rows := make([]model.TrxReconciliation, 0, len(matches))
	summary := &model.RunReconciliationResponse{FileID: fileID}

	for _, m := range matches {
		id, err := uuid.NewV7()
		if err != nil {
			return nil, errors.Wrap(err, wrapErrMsg+"Run.NewV7")
		}
		row := model.TrxReconciliation{
			ID:                  id.String(),
			BankCode:            fileBankCode,
			BankStatementFileID: fileID,
			Status:              m.Status,
			MatchMethod:         m.MatchMethod,
		}
		if m.BankStatementID != "" {
			stmtID := m.BankStatementID
			row.BankStatementID = &stmtID
		}
		if m.TransactionID != "" {
			txID := m.TransactionID
			row.TransactionID = &txID
		}
		rows = append(rows, row)
		incrementSummary(summary, m.Status)
	}

	if err := u.reconRepo.DeleteByFileID(ctx, fileID); err != nil {
		return nil, errors.Wrap(err, wrapErrMsg+"Run.DeleteByFileID")
	}
	if err := u.reconRepo.InsertBatch(ctx, rows); err != nil {
		return nil, errors.Wrap(err, wrapErrMsg+"Run.InsertBatch")
	}

	return summary, nil
}

func (u *reconciliationUC) List(ctx context.Context, param model.ListReconciliationParam) ([]model.ReconciliationResponse, error) {
	normalizedBankCode, err := normalizeBankCodeParam(param.BankCode)
	if err != nil {
		return nil, err
	}

	file, err := u.fileRepo.GetByID(ctx, param.FileID)
	if err != nil {
		if errors.Is(err, constant.ErrNotFound) {
			return nil, commonerr.Set404()
		}
		return nil, errors.Wrap(err, wrapErrMsg+"List.GetByID")
	}

	if normalizeBankCode(file.BankCode) != normalizedBankCode {
		return nil, commonerr.SetNewBadRequest("bank_code_mismatch", "file does not belong to the given bank code")
	}

	limit, offset := normalizeListParams(param.Limit, param.Offset)
	rows, err := u.reconRepo.ListByFileID(ctx, param.FileID, param.Status, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, wrapErrMsg+"List.ListByFileID")
	}

	resp := make([]model.ReconciliationResponse, 0, len(rows))
	for i := range rows {
		resp = append(resp, model.ToReconciliationResponse(&rows[i]))
	}
	return resp, nil
}

func incrementSummary(s *model.RunReconciliationResponse, status string) {
	switch status {
	case constant.ReconciliationStatusMatched:
		s.Matched++
	case constant.ReconciliationStatusMismatch:
		s.Mismatch++
	case constant.ReconciliationStatusOnlyInSystem:
		s.OnlyInSystem++
	case constant.ReconciliationStatusOnlyInBank:
		s.OnlyInBank++
	case constant.ReconciliationStatusAmbiguous:
		s.Ambiguous++
	}
}

func normalizeBankCodeParam(bankCode string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(bankCode))
	if normalized == "" {
		return "", commonerr.SetNewBadRequest("invalid_bank_code", "bank code is required")
	}
	return normalized, nil
}

func normalizeListParams(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = defaultListLimit
	}
	if limit > maxListLimit {
		limit = maxListLimit
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
