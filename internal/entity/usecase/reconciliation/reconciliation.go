package reconciliation

import (
	"context"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
)

type ReconciliationUC interface {
	Run(ctx context.Context, bankCode, fileID string) (*model.RunReconciliationResponse, error)
	List(ctx context.Context, param model.ListReconciliationParam) ([]model.ReconciliationResponse, error)
}
