package messaging

import "log"

type BankStatementProcessQueue struct {
	ch chan string
}

func NewBankStatementProcessQueue(buffer int) *BankStatementProcessQueue {
	if buffer <= 0 {
		buffer = 100
	}
	return &BankStatementProcessQueue{
		ch: make(chan string, buffer),
	}
}

func (q *BankStatementProcessQueue) Publish(fileID string) {
	select {
	case q.ch <- fileID:
	default:
		log.Printf("bank statement process queue is full, fileId=%s", fileID)
	}
}

func (q *BankStatementProcessQueue) Channel() <-chan string {
	return q.ch
}

func (q *BankStatementProcessQueue) Close() {
	close(q.ch)
}
