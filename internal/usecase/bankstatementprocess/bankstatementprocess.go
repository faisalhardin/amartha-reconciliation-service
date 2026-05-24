package bankstatementprocess

import (
	"context"
	"os"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/constant"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
	bankstatementrepo "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/repo/bankstatement"
	bankstatementfilerepo "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/repo/bankstatementfile"
	bankstatementprocessuc "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/usecase/bankstatementprocess"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/messaging"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const wrapErrMsg = "BankStatementProcessUC."

type bankStatementProcessUC struct {
	fileRepo            bankstatementfilerepo.BankStatementFileDB
	stmtRepo            bankstatementrepo.BankStatementDB
	reconciliationQueue *messaging.ReconciliationQueue
}

func NewBankStatementProcessUC(
	fileRepo bankstatementfilerepo.BankStatementFileDB,
	stmtRepo bankstatementrepo.BankStatementDB,
	reconciliationQueue *messaging.ReconciliationQueue,
) bankstatementprocessuc.BankStatementProcessUC {
	return &bankStatementProcessUC{
		fileRepo:            fileRepo,
		stmtRepo:            stmtRepo,
		reconciliationQueue: reconciliationQueue,
	}
}

func (u *bankStatementProcessUC) Process(ctx context.Context, fileID string) error {
	file, err := u.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		return errors.Wrap(err, wrapErrMsg+"Process.GetByID")
	}

	if file.Status != constant.BankStatementFileStatusPending {
		return nil
	}

	content, err := os.ReadFile(file.FileLocation)
	if err != nil {
		if updateErr := u.fileRepo.UpdateStatus(ctx, fileID, constant.BankStatementFileStatusInvalid); updateErr != nil {
			return errors.Wrap(updateErr, wrapErrMsg+"Process.UpdateStatus.ReadFile")
		}
		return errors.Wrap(err, wrapErrMsg+"Process.ReadFile")
	}

	rows, err := parseBankStatementCSV(content)
	if err != nil {
		if updateErr := u.fileRepo.UpdateStatus(ctx, fileID, constant.BankStatementFileStatusInvalid); updateErr != nil {
			return errors.Wrap(updateErr, wrapErrMsg+"Process.UpdateStatus.ParseCSV")
		}
		return errors.Wrap(err, wrapErrMsg+"Process.ParseCSV")
	}

	statements := make([]model.MstBankStatement, 0, len(rows))
	for i, row := range rows {
		stmt, err := buildBankStatementFromCSVRow(file, i+2, row)
		if err != nil {
			if updateErr := u.fileRepo.UpdateStatus(ctx, fileID, constant.BankStatementFileStatusInvalid); updateErr != nil {
				return errors.Wrap(updateErr, wrapErrMsg+"Process.UpdateStatus.ValidateRow")
			}
			return errors.Wrap(err, wrapErrMsg+"Process.ValidateRow")
		}

		id, err := uuid.NewV7()
		if err != nil {
			return errors.Wrap(err, wrapErrMsg+"Process.NewV7")
		}
		stmt.ID = id.String()
		statements = append(statements, *stmt)
	}

	if err := u.stmtRepo.InsertBatch(ctx, statements); err != nil {
		if updateErr := u.fileRepo.UpdateStatus(ctx, fileID, constant.BankStatementFileStatusInvalid); updateErr != nil {
			return errors.Wrap(updateErr, wrapErrMsg+"Process.UpdateStatus.InsertBatch")
		}
		return errors.Wrap(err, wrapErrMsg+"Process.InsertBatch")
	}

	if err := u.fileRepo.UpdateStatus(ctx, fileID, constant.BankStatementFileStatusCompleted); err != nil {
		return errors.Wrap(err, wrapErrMsg+"Process.UpdateStatus.Completed")
	}

	u.reconciliationQueue.Publish(fileID)

	return nil
}
