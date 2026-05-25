package messaging

import (
	"context"

	bankregistryrepo "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/repo/bankregistry"
)

type JobQueueHub struct {
	registry             bankregistryrepo.BankRegistryDB
	processQueues        map[string]*BankStatementProcessQueue
	reconciliationQueues map[string]*ReconciliationQueue
}

func NewJobQueueHub(registry bankregistryrepo.BankRegistryDB, buffer int) *JobQueueHub {
	keys := registry.ListRegistered(context.Background())
	keys = append(keys, bankregistryrepo.CommonRouteKey)

	processQueues := make(map[string]*BankStatementProcessQueue, len(keys))
	reconciliationQueues := make(map[string]*ReconciliationQueue, len(keys))
	for _, key := range keys {
		processQueues[key] = NewBankStatementProcessQueue(buffer)
		reconciliationQueues[key] = NewReconciliationQueue(buffer)
	}

	return &JobQueueHub{
		registry:             registry,
		processQueues:        processQueues,
		reconciliationQueues: reconciliationQueues,
	}
}

func (h *JobQueueHub) PublishProcess(fileID, bankCode string) {
	key := h.registry.RouteKey(bankCode)
	h.processQueues[key].Publish(fileID)
}

func (h *JobQueueHub) PublishReconciliation(fileID, bankCode string) {
	key := h.registry.RouteKey(bankCode)
	h.reconciliationQueues[key].Publish(fileID)
}

func (h *JobQueueHub) ProcessQueues() map[string]*BankStatementProcessQueue {
	out := make(map[string]*BankStatementProcessQueue, len(h.processQueues))
	for k, v := range h.processQueues {
		out[k] = v
	}
	return out
}

func (h *JobQueueHub) ReconciliationQueues() map[string]*ReconciliationQueue {
	out := make(map[string]*ReconciliationQueue, len(h.reconciliationQueues))
	for k, v := range h.reconciliationQueues {
		out[k] = v
	}
	return out
}

func (h *JobQueueHub) Close() {
	for _, q := range h.processQueues {
		q.Close()
	}
	for _, q := range h.reconciliationQueues {
		q.Close()
	}
}
