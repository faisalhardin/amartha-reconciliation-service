package transaction

import (
	"context"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/constant"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
	transactionrepo "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/repo/transaction"
	transactionuc "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/usecase/transaction"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/library/common/commonerr"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	wrapErrMsg       = "TransactionUC."
	defaultListLimit = 20
	maxListLimit     = 50
)

type transactionUC struct {
	repo transactionrepo.TransactionDB
}

func NewTransactionUC(repo transactionrepo.TransactionDB) transactionuc.TransactionUC {
	return &transactionUC{repo: repo}
}

func (u *transactionUC) Create(ctx context.Context, req model.CreateMstTransactionRequest) (*model.MstTransaction, error) {

	id, err := uuid.NewV7()
	if err != nil {
		return nil, errors.Wrap(err, wrapErrMsg+"Create.NewV7")
	}

	tx := &model.MstTransaction{
		ID:              id.String(),
		Amount:          req.Amount,
		Type:            req.Type,
		TransactionTime: req.TransactionTime,
	}

	if err := u.repo.Insert(ctx, tx); err != nil {
		return nil, errors.Wrap(err, wrapErrMsg+"Create.Insert")
	}

	return tx, nil
}

func (u *transactionUC) GetByID(ctx context.Context, id string) (*model.MstTransaction, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, commonerr.SetNewBadRequest("invalid_id", "transaction id must be a valid UUID")
	}

	tx, err := u.repo.GetByID(ctx, parsed)
	if err != nil {
		if errors.Is(err, constant.ErrNotFound) {
			return nil, commonerr.Set404()
		}
		return nil, errors.Wrap(err, wrapErrMsg+"GetByID")
	}

	return tx, nil
}

func (u *transactionUC) List(ctx context.Context, limit, offset int) ([]model.MstTransaction, error) {
	limit, offset = normalizeListParams(limit, offset)

	rows, err := u.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, wrapErrMsg+"List")
	}
	if rows == nil {
		return []model.MstTransaction{}, nil
	}
	return rows, nil
}

func normalizeListParams(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = defaultListLimit
	}
	if limit > maxListLimit {
		limit = maxListLimit
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
