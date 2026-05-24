package transaction

import (
	"context"
	"database/sql"
	"strings"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/constant"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
	transactionrepo "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/repo/transaction"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/library/db/xorm"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const wrapErrMsgPrefix = "TransactionDB."

type transactionDB struct {
	conn *xorm.DBConnect
}

func NewTransactionDB(conn *xorm.DBConnect) transactionrepo.TransactionDB {
	return &transactionDB{conn: conn}
}

func (d *transactionDB) Insert(ctx context.Context, tx *model.MstTransaction) error {
	_, err := d.conn.MasterDB.Context(ctx).Insert(tx)
	if err != nil {
		return errors.Wrap(err, wrapErrMsgPrefix+"Insert")
	}
	return nil
}

func (d *transactionDB) GetByID(ctx context.Context, bankCode string, id uuid.UUID) (*model.MstTransaction, error) {
	var row model.MstTransaction
	has, err := d.conn.MasterDB.Context(ctx).
		Where("id = ? AND bank_code = ?", id.String(), normalizeBankCode(bankCode)).
		Get(&row)
	if err != nil {
		return nil, errors.Wrap(err, wrapErrMsgPrefix+"GetByID")
	}
	if !has {
		return nil, constant.ErrNotFound
	}
	return &row, nil
}

func (d *transactionDB) List(ctx context.Context, bankCode string, limit, offset int) ([]model.MstTransaction, error) {
	var rows []model.MstTransaction
	err := d.conn.MasterDB.Context(ctx).
		Where("bank_code = ?", normalizeBankCode(bankCode)).
		Limit(limit, offset).
		Desc("transaction_time").
		Find(&rows)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, wrapErrMsgPrefix+"List")
	}
	return rows, nil
}

func normalizeBankCode(bankCode string) string {
	return strings.ToUpper(strings.TrimSpace(bankCode))
}
