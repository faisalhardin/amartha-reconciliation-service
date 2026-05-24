package bankstatementfile

import (
	"context"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
)

type BankStatementFileUC interface {
	Upload(ctx context.Context, req model.UploadMstBankStatementFileRequest) (*model.UploadMstBankStatementFileResponse, error)
}
