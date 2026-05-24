package reconciliation

import (
	"context"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
)

type ReconciliationDB interface {
	DeleteByFileID(ctx context.Context, fileID string) error
	InsertBatch(ctx context.Context, rows []model.TrxReconciliation) error
	ListByFileID(ctx context.Context, fileID, status string, limit, offset int) ([]model.TrxReconciliation, error)
	ListAllByFileID(ctx context.Context, fileID string) ([]model.TrxReconciliation, error)
}
