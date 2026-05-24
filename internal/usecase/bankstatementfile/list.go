package bankstatementfile

import (
	"context"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
	"github.com/pkg/errors"
)

const (
	defaultListLimit = 20
	maxListLimit     = 100
)

func (u *bankStatementFileUC) List(ctx context.Context, param model.ListMstBankStatementFileParam) ([]model.MstBankStatementFile, error) {
	normalizedBankCode, err := normalizeBankCode(param.BankCode)
	if err != nil {
		return nil, err
	}

	limit, offset := normalizeListParams(param.Limit, param.Offset)

	rows, err := u.repo.ListByBankCode(ctx, model.ListMstBankStatementFileQuery{
		BankCode: normalizedBankCode,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, errors.Wrap(err, wrapErrMsg+"List.ListByBankCode")
	}
	if rows == nil {
		return []model.MstBankStatementFile{}, nil
	}
	return rows, nil
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
