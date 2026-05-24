package bankstatementfile

import (
	"context"

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
