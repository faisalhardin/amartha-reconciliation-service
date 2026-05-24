package bankstatementprocess

import "context"

type BankStatementProcessUC interface {
	Process(ctx context.Context, fileID string) error
}
