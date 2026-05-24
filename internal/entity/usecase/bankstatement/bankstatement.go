package bankstatement

import (
	"context"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
)

type BankStatementUC interface {
	List(ctx context.Context, param model.ListMstBankStatementParam) ([]model.MstBankStatement, error)
}
