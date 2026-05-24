package transaction

import (
	"context"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
	"github.com/google/uuid"
)

type TransactionDB interface {
	Insert(ctx context.Context, tx *model.MstTransaction) error
	InsertBatch(ctx context.Context, txs []model.MstTransaction) error
	GetByID(ctx context.Context, bankCode string, id uuid.UUID) (*model.MstTransaction, error)
	List(ctx context.Context, bankCode string, limit, offset int) ([]model.MstTransaction, error)
}
