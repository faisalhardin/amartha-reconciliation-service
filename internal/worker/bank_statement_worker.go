package worker

import (
	"context"
	"log"

	bankstatementprocessuc "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/usecase/bankstatementprocess"
	"github.com/faisalhardin/amartha-reconciliation-service/internal/messaging"
)

func StartBankStatementWorker(
	ctx context.Context,
	queue *messaging.BankStatementProcessQueue,
	uc bankstatementprocessuc.BankStatementProcessUC,
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
				if err := uc.Process(ctx, fileID); err != nil {
					log.Printf("bank statement worker process error fileId=%s: %v", fileID, err)
				}
			}
		}
	}()
}
