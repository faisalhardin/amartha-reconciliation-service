package reconciliation

import (
	"testing"

	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/constant"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/entity/model"
	"github.com/shopspring/decimal"
)

func TestReferencePass(t *testing.T) {
	t.Run("MATCHED when bank id equals our id and amounts agree", func(t *testing.T) {
		stmts := []model.MstBankStatement{
			{
				ID:               "stmt-1",
				BankCode:         "BCA",
				UniqueIdentifier: "TX-SAME",
				Amount:           decimal.NewFromInt(-1_500_000),
				Date:             1711000000,
			},
		}
		txs := []model.MstTransaction{
			{
				ID:                   "tx-1",
				BankCode:             "BCA",
				TransactionReference: "TX-SAME",
				Amount:               decimal.NewFromInt(1_500_000),
				Type:                 model.TransactionTypeDebit,
				TransactionTime:      999, // time is ignored in reference pass
			},
		}

		results := referencePass(&stmts, &txs)

		if len(results) != 1 {
			t.Fatalf("want 1 result, got %d", len(results))
		}
		r := results[0]
		if r.Status != constant.ReconciliationStatusMatched {
			t.Errorf("status: got %s want MATCHED", r.Status)
		}
		if r.MatchMethod != constant.ReconciliationMatchMethodReference {
			t.Errorf("method: got %s want REFERENCE", r.MatchMethod)
		}
		if r.BankStatementID != "stmt-1" || r.TransactionID != "tx-1" {
			t.Errorf("ids: stmt=%s tx=%s", r.BankStatementID, r.TransactionID)
		}
		if len(stmts) != 0 {
			t.Errorf("matched stmt should leave pool; pool len=%d", len(stmts))
		}
		if len(txs) != 0 {
			t.Errorf("matched tx should leave pool; pool len=%d", len(txs))
		}
	})

	t.Run("MISMATCH when ids match but bank amount differs", func(t *testing.T) {
		stmts := []model.MstBankStatement{
			{
				ID:               "stmt-1",
				BankCode:         "BCA",
				UniqueIdentifier: "TX-002",
				Amount:           decimal.NewFromInt(-1_400_000), // bank says -1.4M
				Date:             1711003600,
			},
		}
		txs := []model.MstTransaction{
			{
				ID:                   "tx-1",
				BankCode:             "BCA",
				TransactionReference: "TX-002",
				Amount:               decimal.NewFromInt(1_500_000), // we recorded 1.5M DEBIT
				Type:                 model.TransactionTypeDebit,
				TransactionTime:      1711003600,
			},
		}

		results := referencePass(&stmts, &txs)

		if len(results) != 1 {
			t.Fatalf("want 1 result, got %d", len(results))
		}
		if results[0].Status != constant.ReconciliationStatusMismatch {
			t.Errorf("status: got %s want MISMATCH", results[0].Status)
		}
		if len(stmts) != 0 {
			t.Error("mismatch still pairs rows; stmt should not stay in pool")
		}
		if len(txs) != 0 {
			t.Error("paired tx should be removed from pool")
		}
	})

	t.Run("no result when ids differ — stmt stays for composite pass", func(t *testing.T) {
		stmts := []model.MstBankStatement{
			{
				ID:               "stmt-1",
				BankCode:         "BCA",
				UniqueIdentifier: "STMT-001",
				Amount:           decimal.NewFromInt(-100),
				Date:             1711000000,
			},
		}
		txs := []model.MstTransaction{
			{
				ID:                   "tx-1",
				BankCode:             "BCA",
				TransactionReference: "TX-001",
				Amount:               decimal.NewFromInt(100),
				Type:                 model.TransactionTypeDebit,
				TransactionTime:      1711000000,
			},
		}

		results := referencePass(&stmts, &txs)

		if len(results) != 0 {
			t.Fatalf("want no reference results, got %d: %+v", len(results), results)
		}
		if len(stmts) != 1 || stmts[0].ID != "stmt-1" {
			t.Errorf("stmt should remain in pool for next pass: %+v", stmts)
		}
		if len(txs) != 1 {
			t.Errorf("tx should stay in pool: len=%d", len(txs))
		}
	})

	t.Run("AMBIGUOUS when two transactions share the same reference", func(t *testing.T) {
		stmts := []model.MstBankStatement{
			{
				ID:               "stmt-1",
				BankCode:         "BCA",
				UniqueIdentifier: "TX-DUP",
				Amount:           decimal.NewFromInt(-100),
				Date:             1711000000,
			},
		}
		txs := []model.MstTransaction{
			{ID: "tx-a", BankCode: "BCA", TransactionReference: "TX-DUP", Amount: decimal.NewFromInt(100), Type: model.TransactionTypeDebit},
			{ID: "tx-b", BankCode: "BCA", TransactionReference: "TX-DUP", Amount: decimal.NewFromInt(100), Type: model.TransactionTypeDebit},
		}

		results := referencePass(&stmts, &txs)

		if len(results) != 1 {
			t.Fatalf("want 1 result, got %d", len(results))
		}
		r := results[0]
		if r.Status != constant.ReconciliationStatusAmbiguous {
			t.Errorf("status: got %s want AMBIGUOUS", r.Status)
		}
		if r.TransactionID != "" {
			t.Error("ambiguous should not pick a transaction id")
		}
		if len(stmts) != 0 {
			t.Error("ambiguous stmt is removed from pool (only the AMBIGUOUS result remains)")
		}
	})

	t.Run("empty bank unique_identifier goes to remaining pool", func(t *testing.T) {
		stmts := []model.MstBankStatement{
			{ID: "stmt-1", BankCode: "BCA", UniqueIdentifier: "  ", Amount: decimal.NewFromInt(-50)},
		}
		txs := []model.MstTransaction{
			{ID: "tx-1", BankCode: "BCA", TransactionReference: "", Amount: decimal.NewFromInt(50), Type: model.TransactionTypeDebit},
		}

		results := referencePass(&stmts, &txs)

		if len(results) != 0 {
			t.Fatalf("want no results, got %d", len(results))
		}
		if len(stmts) != 1 {
			t.Fatalf("want stmt in pool, got len=%d", len(stmts))
		}
	})

	t.Run("different bank_code never pairs", func(t *testing.T) {
		stmts := []model.MstBankStatement{
			{ID: "stmt-1", BankCode: "BCA", UniqueIdentifier: "TX-1", Amount: decimal.NewFromInt(-100)},
		}
		txs := []model.MstTransaction{
			{ID: "tx-1", BankCode: "MANDIRI", TransactionReference: "TX-1", Amount: decimal.NewFromInt(100), Type: model.TransactionTypeDebit},
		}

		results := referencePass(&stmts, &txs)

		if len(results) != 0 {
			t.Fatalf("cross-bank should not reference-match, got %+v", results)
		}
		if len(stmts) != 1 {
			t.Error("stmt should stay in pool")
		}
	})
}

func TestMatchStatements_compositeMatch(t *testing.T) {
	stmts := []model.MstBankStatement{
		{ID: "s1", BankCode: "BCA", UniqueIdentifier: "STMT-001", Amount: decimal.NewFromInt(-100), Date: 1710000000},
	}
	txs := []model.MstTransaction{
		{ID: "t1", BankCode: "BCA", TransactionReference: "TX-001", Amount: decimal.NewFromInt(100), Type: model.TransactionTypeDebit, TransactionTime: 1710000000},
	}

	results := matchStatements(stmts, txs)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != constant.ReconciliationStatusMatched {
		t.Fatalf("expected MATCHED, got %s", results[0].Status)
	}
	if results[0].MatchMethod != constant.ReconciliationMatchMethodComposite {
		t.Fatalf("expected COMPOSITE, got %s", results[0].MatchMethod)
	}
}

func TestMatchStatements_referenceMismatch(t *testing.T) {
	stmts := []model.MstBankStatement{
		{ID: "s1", BankCode: "BCA", UniqueIdentifier: "TX-002", Amount: decimal.NewFromInt(-140), Date: 1710000000},
	}
	txs := []model.MstTransaction{
		{ID: "t1", BankCode: "BCA", TransactionReference: "TX-002", Amount: decimal.NewFromInt(150), Type: model.TransactionTypeDebit, TransactionTime: 1710000000},
	}

	results := matchStatements(stmts, txs)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != constant.ReconciliationStatusMismatch {
		t.Fatalf("expected MISMATCH, got %s", results[0].Status)
	}
}

func TestMatchStatements_crossBankNoMatch(t *testing.T) {
	stmts := []model.MstBankStatement{
		{ID: "s1", BankCode: "BCA", UniqueIdentifier: "STMT-001", Amount: decimal.NewFromInt(-100), Date: 1710000000},
	}
	txs := []model.MstTransaction{
		{ID: "t1", BankCode: "MANDIRI", TransactionReference: "TX-001", Amount: decimal.NewFromInt(100), Type: model.TransactionTypeDebit, TransactionTime: 1710000000},
	}

	results := matchStatements(stmts, txs)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}
