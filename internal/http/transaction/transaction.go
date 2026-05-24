package transaction

import (
	"net/http"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
	transactionuc "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/usecase/transaction"
	httpwriter "github.com/faisalhardin/amartha-reconciliation-service/internal/library/common/writer"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/library/util/common/binding"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	TransactionUC transactionuc.TransactionUC
}

func New(uc transactionuc.TransactionUC) *Handler {
	return &Handler{TransactionUC: uc}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateMstTransactionRequest
	req.BankCode = chi.URLParam(r, "bankCode")
	if err := binding.Bind(r, &req); err != nil {
		_ = httpwriter.SetError(r.Context(), w, err)
		return
	}

	tx, err := h.TransactionUC.Create(r.Context(), req)
	if err != nil {
		_ = httpwriter.SetError(r.Context(), w, err)
		return
	}

	_, _ = httpwriter.WriteJSONAPIData(w, r, http.StatusCreated, model.ToMstTransactionResponse(tx))
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	bankCode := chi.URLParam(r, "bankCode")
	id := chi.URLParam(r, "id")

	tx, err := h.TransactionUC.GetByID(r.Context(), bankCode, id)
	if err != nil {
		_ = httpwriter.SetError(r.Context(), w, err)
		return
	}

	_ = httpwriter.SetOKWithData(r.Context(), w, model.ToMstTransactionResponse(tx))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	var param model.ListMstTransactionParam
	param.BankCode = chi.URLParam(r, "bankCode")
	if err := binding.Bind(r, &param); err != nil {
		_ = httpwriter.SetError(r.Context(), w, err)
		return
	}

	rows, err := h.TransactionUC.List(r.Context(), param)
	if err != nil {
		_ = httpwriter.SetError(r.Context(), w, err)
		return
	}

	resp := make([]model.MstTransactionResponse, 0, len(rows))
	for i := range rows {
		resp = append(resp, model.ToMstTransactionResponse(&rows[i]))
	}

	_ = httpwriter.SetOKWithData(r.Context(), w, resp)
}
