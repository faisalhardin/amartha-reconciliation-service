package bankstatementfile

import (
	"context"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
)

type BankStatementFileDB interface {
	Insert(ctx context.Context, file *model.MstBankStatementFile) error
	GetByID(ctx context.Context, fileID string) (*model.MstBankStatementFile, error)
	UpdateStatus(ctx context.Context, fileID, status string) error
	ListByBankCode(ctx context.Context, query model.ListMstBankStatementFileQuery) ([]model.MstBankStatementFile, error)
}
