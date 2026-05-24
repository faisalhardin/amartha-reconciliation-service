package reconciliation

import (
	"net/http"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
	reconciliationuc "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/usecase/reconciliation"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/library/common/commonerr"
	httpwriter "github.com/faisalhardin/amartha-reconciliation-service/internal/library/common/writer"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/library/util/common/binding"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	ReconciliationUC reconciliationuc.ReconciliationUC
}

func New(uc reconciliationuc.ReconciliationUC) *Handler {
	return &Handler{ReconciliationUC: uc}
}

func (h *Handler) Run(w http.ResponseWriter, r *http.Request) {
	bankCode := chi.URLParam(r, "bankCode")
	fileID := r.URL.Query().Get("fileId")
	if fileID == "" {
		_ = httpwriter.SetError(r.Context(), w, commonerr.SetNewBadRequest("missing_file_id", "fileId query parameter is required"))
		return
	}

	result, err := h.ReconciliationUC.Run(r.Context(), bankCode, fileID)
	if err != nil {
		_ = httpwriter.SetError(r.Context(), w, err)
		return
	}

	_ = httpwriter.SetOKWithData(r.Context(), w, result)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	var param model.ListReconciliationParam
	param.BankCode = chi.URLParam(r, "bankCode")
	if err := binding.Bind(r, &param); err != nil {
		_ = httpwriter.SetError(r.Context(), w, err)
		return
	}

	rows, err := h.ReconciliationUC.List(r.Context(), param)
	if err != nil {
		_ = httpwriter.SetError(r.Context(), w, err)
		return
	}

	_ = httpwriter.SetOKWithData(r.Context(), w, rows)
}
