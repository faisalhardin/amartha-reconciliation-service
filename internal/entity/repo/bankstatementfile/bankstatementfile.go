package bankstatementfile

import (
	"context"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
)

type BankStatementFileDB interface {
	Insert(ctx context.Context, file *model.MstBankStatementFile) error
}
