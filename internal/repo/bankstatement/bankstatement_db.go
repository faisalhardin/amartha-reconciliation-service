package bankstatement

import (
	"context"
	"database/sql"
	"strings"

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

func (d *bankStatementDB) InsertBatch(ctx context.Context, stmts []model.MstBankStatement) error {
	if len(stmts) == 0 {
		return nil
	}
	_, err := d.conn.MasterDB.Context(ctx).Insert(&stmts)
	if err != nil {
		return errors.Wrap(err, wrapErrMsgPrefix+"InsertBatch")
	}
	return nil
}

func (d *bankStatementDB) GetByID(ctx context.Context, id uuid.UUID) (*model.MstBankStatement, error) {
	var row model.MstBankStatement
	has, err := d.conn.MasterDB.Context(ctx).
		Where("id = ?", id.String()).
		NoAutoCondition().
		Get(&row)
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

func (d *bankStatementDB) ListByFileID(ctx context.Context, fileID string) ([]model.MstBankStatement, error) {
	var rows []model.MstBankStatement
	err := d.conn.MasterDB.Context(ctx).
		Where("id_bank_statement_file = ?", fileID).
		Find(&rows)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, wrapErrMsgPrefix+"ListByFileID")
	}
	return rows, nil
}

func (d *bankStatementDB) ListByBankCodeAndFileID(ctx context.Context, query model.ListMstBankStatementQuery) ([]model.MstBankStatement, error) {
	var rows []model.MstBankStatement
	err := d.conn.MasterDB.Context(ctx).
		Where("bank_code = ? AND id_bank_statement_file = ?", normalizeBankCode(query.BankCode), query.FileID).
		Limit(query.Limit, query.Offset).
		Desc("date").
		Find(&rows)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, wrapErrMsgPrefix+"ListByBankCodeAndFileID")
	}
	return rows, nil
}

func normalizeBankCode(bankCode string) string {
	return strings.ToUpper(strings.TrimSpace(bankCode))
}
