package bankstatement

import (
	"context"
	"strings"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/constant"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
	bankstatementrepo "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/repo/bankstatement"
	bankstatementfilerepo "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/repo/bankstatementfile"
	bankstatementuc "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/usecase/bankstatement"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/library/common/commonerr"
	"github.com/pkg/errors"
)

const (
	wrapErrMsg         = "BankStatementUC."
	defaultListLimit   = 20
	maxListLimit       = 100
)

type bankStatementUC struct {
	stmtRepo bankstatementrepo.BankStatementDB
	fileRepo bankstatementfilerepo.BankStatementFileDB
}

func NewBankStatementUC(
	stmtRepo bankstatementrepo.BankStatementDB,
	fileRepo bankstatementfilerepo.BankStatementFileDB,
) bankstatementuc.BankStatementUC {
	return &bankStatementUC{
		stmtRepo: stmtRepo,
		fileRepo: fileRepo,
	}
}

func (u *bankStatementUC) List(ctx context.Context, param model.ListMstBankStatementParam) ([]model.MstBankStatement, error) {
	normalizedBankCode, err := normalizeBankCode(param.BankCode)
	if err != nil {
		return nil, err
	}

	fileID := strings.TrimSpace(param.FileID)
	if fileID == "" {
		return nil, commonerr.SetNewBadRequest("missing_file_id", "fileId query parameter is required")
	}

	file, err := u.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		if errors.Is(err, constant.ErrNotFound) {
			return nil, commonerr.Set404()
		}
		return nil, errors.Wrap(err, wrapErrMsg+"List.GetByID")
	}

	if strings.ToUpper(strings.TrimSpace(file.BankCode)) != normalizedBankCode {
		return nil, commonerr.SetNewBadRequest("bank_code_mismatch", "file does not belong to the given bank code")
	}

	limit, offset := normalizeListParams(param.Limit, param.Offset)

	rows, err := u.stmtRepo.ListByBankCodeAndFileID(ctx, model.ListMstBankStatementQuery{
		BankCode: normalizedBankCode,
		FileID:   fileID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, errors.Wrap(err, wrapErrMsg+"List.ListByBankCodeAndFileID")
	}
	if rows == nil {
		return []model.MstBankStatement{}, nil
	}
	return rows, nil
}

func normalizeBankCode(bankCode string) (string, error) {
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
