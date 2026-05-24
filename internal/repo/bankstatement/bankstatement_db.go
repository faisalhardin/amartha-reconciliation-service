package bankstatement

import (
	"context"
	"database/sql"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/constant"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
	bankstatementrepo "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/repo/bankstatement"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/library/db/xorm"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const wrapErrMsgPrefix = "BankStatementDB."

type bankStatementDB struct {
	conn *xorm.DBConnect
}

func NewBankStatementDB(conn *xorm.DBConnect) bankstatementrepo.BankStatementDB {
	return &bankStatementDB{conn: conn}
}

func (d *bankStatementDB) Insert(ctx context.Context, stmt *model.MstBankStatement) error {
	_, err := d.conn.MasterDB.Context(ctx).Insert(stmt)
	if err != nil {
		return errors.Wrap(err, wrapErrMsgPrefix+"Insert")
	}
	return nil
}

func (d *bankStatementDB) GetByID(ctx context.Context, id uuid.UUID) (*model.MstBankStatement, error) {
	var row model.MstBankStatement
	has, err := d.conn.MasterDB.Context(ctx).ID(id.String()).Get(&row)
	if err != nil {
		return nil, errors.Wrap(err, wrapErrMsgPrefix+"GetByID")
	}
	if !has {
		return nil, constant.ErrNotFound
	}
	return &row, nil
}

func (d *bankStatementDB) GetByUniqueIdentifier(ctx context.Context, uniqueID string) (*model.MstBankStatement, error) {
	var row model.MstBankStatement
	has, err := d.conn.MasterDB.Context(ctx).
		Where("unique_identifier = ?", uniqueID).
		Get(&row)
	if err != nil {
		return nil, errors.Wrap(err, wrapErrMsgPrefix+"GetByUniqueIdentifier")
	}
	if !has {
		return nil, constant.ErrNotFound
	}
	return &row, nil
}

func (d *bankStatementDB) List(ctx context.Context, limit, offset int) ([]model.MstBankStatement, error) {
	var rows []model.MstBankStatement
	err := d.conn.MasterDB.Context(ctx).
		Limit(limit, offset).
		Desc("date").
		Find(&rows)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, wrapErrMsgPrefix+"List")
	}
	return rows, nil
}
