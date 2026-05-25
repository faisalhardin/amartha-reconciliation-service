package worker

import (
	"context"
	"log"

	reconciliationuc "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/usecase/reconciliation"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/messaging"
)

func StartReconciliationWorker(
	ctx context.Context,
	bankKey string,
	queue *messaging.ReconciliationQueue,
	uc reconciliationuc.ReconciliationUC,
) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case fileID, ok := <-queue.Channel():
				if !ok {
					return
				}
				if _, err := uc.Run(ctx, "", fileID); err != nil {
					log.Printf("reconciliation worker [%s] error fileId=%s: %v", bankKey, fileID, err)
				}
			}
		}
	}()
}
