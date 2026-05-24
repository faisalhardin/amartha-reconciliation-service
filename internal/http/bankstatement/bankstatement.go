package bankstatement

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
	bankstatementuc "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/usecase/bankstatement"
	bankstatementfileuc "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/usecase/bankstatementfile"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/library/common/commonerr"
	httpwriter "github.com/faisalhardin/amartha-reconciliation-service/internal/library/common/writer"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/library/util/common/binding"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	BankStatementFileUC bankstatementfileuc.BankStatementFileUC
	BankStatementUC     bankstatementuc.BankStatementUC
}

func New(fileUC bankstatementfileuc.BankStatementFileUC, stmtUC bankstatementuc.BankStatementUC) *Handler {
	return &Handler{
		BankStatementFileUC: fileUC,
		BankStatementUC:     stmtUC,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	var param model.ListMstBankStatementParam
	param.BankCode = chi.URLParam(r, "bankCode")
	if err := binding.Bind(r, &param); err != nil {
		_ = httpwriter.SetError(r.Context(), w, err)
		return
	}

	rows, err := h.BankStatementUC.List(r.Context(), param)
	if err != nil {
		_ = httpwriter.SetError(r.Context(), w, err)
		return
	}

	resp := make([]model.MstBankStatementResponse, 0, len(rows))
	for i := range rows {
		resp = append(resp, model.ToMstBankStatementResponse(&rows[i]))
	}

	_ = httpwriter.SetOKWithData(r.Context(), w, resp)
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		_ = httpwriter.SetError(r.Context(), w, commonerr.SetNewBadRequest("invalid_form", "failed to parse multipart form"))
		return
	}

	bankCode := strings.TrimSpace(r.FormValue("bank_code"))
	if bankCode == "" {
		_ = httpwriter.SetError(r.Context(), w, commonerr.SetNewBadRequest("missing_bank_code", "bank_code is required"))
		return
	}

	startDate, err := parseFormInt64(r.FormValue("start_date"), "start_date")
	if err != nil {
		_ = httpwriter.SetError(r.Context(), w, err)
		return
	}

	endDate, err := parseFormInt64(r.FormValue("end_date"), "end_date")
	if err != nil {
		_ = httpwriter.SetError(r.Context(), w, err)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		_ = httpwriter.SetError(r.Context(), w, commonerr.SetNewBadRequest("missing_file", "file is required"))
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		_ = httpwriter.SetError(r.Context(), w, commonerr.SetNewBadRequest("invalid_file", "failed to read uploaded file"))
		return
	}

	req := model.UploadMstBankStatementFileRequest{
		BankCode:    bankCode,
		StartDate:   startDate,
		EndDate:     endDate,
		FileName:    header.Filename,
		FileContent: content,
	}

	result, err := h.BankStatementFileUC.Upload(r.Context(), req)
	if err != nil {
		_ = httpwriter.SetError(r.Context(), w, err)
		return
	}

	_, _ = httpwriter.WriteJSONAPIData(w, r, http.StatusCreated, result)
}

func parseFormInt64(value, field string) (int64, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, commonerr.SetNewBadRequest("missing_"+field, field+" is required")
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, commonerr.SetNewBadRequest("invalid_"+field, field+" must be a valid unix timestamp")
	}
	return parsed, nil
}
