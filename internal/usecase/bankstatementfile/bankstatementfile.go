package bankstatementfile

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/constant"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
	bankstatementfilerepo "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/repo/bankstatementfile"
	bankstatementfileuc "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/usecase/bankstatementfile"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/library/common/commonerr"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/messaging"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const wrapErrMsg = "BankStatementFileUC."

type bankStatementFileUC struct {
	repo  bankstatementfilerepo.BankStatementFileDB
	queue *messaging.BankStatementProcessQueue
}

func NewBankStatementFileUC(
	repo bankstatementfilerepo.BankStatementFileDB,
	queue *messaging.BankStatementProcessQueue,
) bankstatementfileuc.BankStatementFileUC {
	return &bankStatementFileUC{
		repo:  repo,
		queue: queue,
	}
}

func (u *bankStatementFileUC) Upload(ctx context.Context, req model.UploadMstBankStatementFileRequest) (*model.UploadMstBankStatementFileResponse, error) {
	bankCode, err := normalizeBankCode(req.BankCode)
	if err != nil {
		return nil, err
	}

	if err := validateUploadRequest(req); err != nil {
		return nil, err
	}

	fileID, err := uuid.NewV7()
	if err != nil {
		return nil, errors.Wrap(err, wrapErrMsg+"Upload.NewV7")
	}

	safeFileName := sanitizeFileName(req.FileName)
	fileLocation := filepath.Join(os.TempDir(), fmt.Sprintf("%s_%s", fileID.String(), safeFileName))

	if err := os.WriteFile(fileLocation, req.FileContent, 0o644); err != nil {
		return nil, errors.Wrap(err, wrapErrMsg+"Upload.WriteFile")
	}

	record := &model.MstBankStatementFile{
		ID:           fileID.String(),
		BankCode:     bankCode,
		FileName:     req.FileName,
		FileLocation: fileLocation,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Status:       constant.BankStatementFileStatusPending,
	}

	if err := u.repo.Insert(ctx, record); err != nil {
		_ = os.Remove(fileLocation)
		return nil, errors.Wrap(err, wrapErrMsg+"Upload.Insert")
	}

	u.queue.Publish(record.ID)

	return &model.UploadMstBankStatementFileResponse{FileID: record.ID}, nil
}

func validateUploadRequest(req model.UploadMstBankStatementFileRequest) error {
	if len(req.FileContent) == 0 {
		return commonerr.SetNewBadRequest("empty_file", "uploaded file is empty")
	}
	if strings.TrimSpace(req.FileName) == "" {
		return commonerr.SetNewBadRequest("invalid_file_name", "file name is required")
	}
	if req.StartDate <= 0 {
		return commonerr.SetNewBadRequest("invalid_start_date", "start_date must be a positive unix timestamp")
	}
	if req.EndDate <= 0 {
		return commonerr.SetNewBadRequest("invalid_end_date", "end_date must be a positive unix timestamp")
	}
	if req.StartDate > req.EndDate {
		return commonerr.SetNewBadRequest("invalid_date_range", "start_date must be less than or equal to end_date")
	}
	return nil
}

func normalizeBankCode(bankCode string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(bankCode))
	if normalized == "" {
		return "", commonerr.SetNewBadRequest("invalid_bank_code", "bank_code is required")
	}
	return normalized, nil
}

func sanitizeFileName(name string) string {
	base := filepath.Base(strings.TrimSpace(name))
	if base == "." || base == string(filepath.Separator) {
		return "upload.csv"
	}
	return base
}
