package reconciliation

import (
	"testing"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/constant"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
	"github.com/shopspring/decimal"
)

func TestIncrementSummaryCounts(t *testing.T) {
	s := &model.ReconciliationSummaryResponse{}
	incrementSummaryCounts(s, constant.ReconciliationStatusMatched)
	incrementSummaryCounts(s, constant.ReconciliationStatusMismatch)
	incrementSummaryCounts(s, constant.ReconciliationStatusOnlyInSystem)
	incrementSummaryCounts(s, constant.ReconciliationStatusOnlyInBank)
	incrementSummaryCounts(s, constant.ReconciliationStatusAmbiguous)

	if s.TotalProcessed != 5 {
		t.Fatalf("totalProcessed: got %d want 5", s.TotalProcessed)
	}
	if s.TotalMatched != 1 || s.TotalMismatch != 1 || s.TotalUnmatched != 2 || s.TotalAmbiguous != 1 {
		t.Fatalf("counts: matched=%d mismatch=%d unmatched=%d ambiguous=%d",
			s.TotalMatched, s.TotalMismatch, s.TotalUnmatched, s.TotalAmbiguous)
	}
}

func TestMismatchDiscrepancyAmount(t *testing.T) {
	stmt := &model.MstBankStatement{Amount: decimal.NewFromInt(-1400000)}
	tx := &model.MstTransaction{
		Amount: decimal.NewFromInt(1500000),
		Type:   model.TransactionTypeDebit,
	}
	diff := stmt.Amount.Sub(signedAmount(*tx)).Abs()
	if !diff.Equal(decimal.NewFromInt(100000)) {
		t.Fatalf("discrepancy: got %s want 100000", diff)
	}
}
