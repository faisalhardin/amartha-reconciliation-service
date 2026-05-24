package messaging

import "log"

type ReconciliationQueue struct {
	ch chan string
}

func NewReconciliationQueue(buffer int) *ReconciliationQueue {
	if buffer <= 0 {
		buffer = 100
	}
	return &ReconciliationQueue{
		ch: make(chan string, buffer),
	}
}

func (q *ReconciliationQueue) Publish(fileID string) {
	select {
	case q.ch <- fileID:
	default:
		log.Printf("reconciliation queue is full, fileId=%s", fileID)
	}
}

func (q *ReconciliationQueue) Channel() <-chan string {
	return q.ch
}

func (q *ReconciliationQueue) Close() {
	close(q.ch)
}
