package http

import (
	"github.com/faisalhardin/amartha-reconciliation-service/internal/http/transaction"
)

type Handlers struct {
	TransactionHandler *transaction.Handler
}
