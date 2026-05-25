package worker

import (
	"context"

	bankstatementprocessuc "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/usecase/bankstatementprocess"
	reconciliationuc "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/usecase/reconciliation"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/messaging"
)

func StartAll(
	ctx context.Context,
	hub *messaging.JobQueueHub,
	processUC bankstatementprocessuc.BankStatementProcessUC,
	reconUC reconciliationuc.ReconciliationUC,
) {
	for key, queue := range hub.ProcessQueues() {
		StartBankStatementWorker(ctx, key, queue, processUC)
	}
	for key, queue := range hub.ReconciliationQueues() {
		StartReconciliationWorker(ctx, key, queue, reconUC)
	}
}
