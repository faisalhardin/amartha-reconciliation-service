package reconciliation

import (
	"context"
	"database/sql"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
	reconciliationrepo "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/repo/reconciliation"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/library/db/xorm"
	"github.com/pkg/errors"
)

const wrapErrMsgPrefix = "ReconciliationDB."

type reconciliationDB struct {
	conn *xorm.DBConnect
}

func NewReconciliationDB(conn *xorm.DBConnect) reconciliationrepo.ReconciliationDB {
	return &reconciliationDB{conn: conn}
}

func (d *reconciliationDB) DeleteByFileID(ctx context.Context, fileID string) error {
	_, err := d.conn.MasterDB.Context(ctx).
		Where("id_bank_statement_file = ?", fileID).
		Delete(&model.TrxReconciliation{})
	if err != nil {
		return errors.Wrap(err, wrapErrMsgPrefix+"DeleteByFileID")
	}
	return nil
}

func (d *reconciliationDB) InsertBatch(ctx context.Context, rows []model.TrxReconciliation) error {
	if len(rows) == 0 {
		return nil
	}
	_, err := d.conn.MasterDB.Context(ctx).Insert(&rows)
	if err != nil {
		return errors.Wrap(err, wrapErrMsgPrefix+"InsertBatch")
	}
	return nil
}

func (d *reconciliationDB) ListAllByFileID(ctx context.Context, fileID string) ([]model.TrxReconciliation, error) {
	var rows []model.TrxReconciliation
	err := d.conn.MasterDB.Context(ctx).
		Where("id_bank_statement_file = ?", fileID).
		Find(&rows)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, wrapErrMsgPrefix+"ListAllByFileID")
	}
	return rows, nil
}

func (d *reconciliationDB) ListByFileID(ctx context.Context, fileID, status string, limit, offset int) ([]model.TrxReconciliation, error) {
	sess := d.conn.MasterDB.Context(ctx).
		Where("id_bank_statement_file = ?", fileID)
	if status != "" {
		sess = sess.And("status = ?", status)
	}
	var rows []model.TrxReconciliation
	err := sess.Limit(limit, offset).Desc("create_time").Find(&rows)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, wrapErrMsgPrefix+"ListByFileID")
	}
	return rows, nil
}
