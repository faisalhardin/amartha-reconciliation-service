package transaction

import (
	"context"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
)

type TransactionUC interface {
	Create(ctx context.Context, req model.CreateMstTransactionRequest) (*model.MstTransaction, error)
	UploadFromCSV(ctx context.Context, bankCode string, file []byte) (*model.UploadMstTransactionCSVResponse, error)
	GetByID(ctx context.Context, bankCode, id string) (*model.MstTransaction, error)
	List(ctx context.Context, param model.ListMstTransactionParam) ([]model.MstTransaction, error)
}
