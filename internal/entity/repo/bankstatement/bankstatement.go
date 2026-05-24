package bankstatement

import (
	"context"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
	"github.com/google/uuid"
)

type BankStatementDB interface {
	Insert(ctx context.Context, stmt *model.MstBankStatement) error
	InsertBatch(ctx context.Context, stmts []model.MstBankStatement) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.MstBankStatement, error)
	GetByUniqueIdentifier(ctx context.Context, uniqueID string) (*model.MstBankStatement, error)
	List(ctx context.Context, limit, offset int) ([]model.MstBankStatement, error)
}
