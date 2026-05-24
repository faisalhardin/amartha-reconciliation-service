package transaction

import (
	"context"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
)

type TransactionUC interface {
	Create(ctx context.Context, req model.CreateMstTransactionRequest) (*model.MstTransaction, error)
	GetByID(ctx context.Context, id string) (*model.MstTransaction, error)
	List(ctx context.Context, limit, offset int) ([]model.MstTransaction, error)
}
