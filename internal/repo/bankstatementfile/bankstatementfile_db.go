package bankstatementfile

import (
	"context"
	"database/sql"
	"strings"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/constant"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
	bankstatementfilerepo "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/repo/bankstatementfile"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/library/db/xorm"
	"github.com/pkg/errors"
)

const wrapErrMsgPrefix = "BankStatementFileDB."

type bankStatementFileDB struct {
	conn *xorm.DBConnect
}

func NewBankStatementFileDB(conn *xorm.DBConnect) bankstatementfilerepo.BankStatementFileDB {
	return &bankStatementFileDB{conn: conn}
}

func (d *bankStatementFileDB) Insert(ctx context.Context, file *model.MstBankStatementFile) error {
	_, err := d.conn.MasterDB.Context(ctx).Insert(file)
	if err != nil {
		return errors.Wrap(err, wrapErrMsgPrefix+"Insert")
	}
	return nil
}

func (d *bankStatementFileDB) GetByID(ctx context.Context, fileID string) (*model.MstBankStatementFile, error) {
	var row model.MstBankStatementFile
	has, err := d.conn.MasterDB.Context(ctx).ID(fileID).Get(&row)
	if err != nil {
		return nil, errors.Wrap(err, wrapErrMsgPrefix+"GetByID")
	}
	if !has {
		return nil, constant.ErrNotFound
	}
	return &row, nil
}

func (d *bankStatementFileDB) UpdateStatus(ctx context.Context, fileID, status string) error {
	_, err := d.conn.MasterDB.Context(ctx).
		ID(fileID).
		Cols("status", "update_time").
		Update(&model.MstBankStatementFile{Status: status})
	if err != nil {
		return errors.Wrap(err, wrapErrMsgPrefix+"UpdateStatus")
	}
	return nil
}

func (d *bankStatementFileDB) ListByBankCode(ctx context.Context, query model.ListMstBankStatementFileQuery) ([]model.MstBankStatementFile, error) {
	var rows []model.MstBankStatementFile
	err := d.conn.MasterDB.Context(ctx).
		Where("bank_code = ?", normalizeBankCode(query.BankCode)).
		Limit(query.Limit, query.Offset).
		Desc("create_time").
		Find(&rows)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, wrapErrMsgPrefix+"ListByBankCode")
	}
	return rows, nil
}

func normalizeBankCode(bankCode string) string {
	return strings.ToUpper(strings.TrimSpace(bankCode))
}
