package bankregistry

import "context"

const CommonRouteKey = "COMMON"

type BankRegistryDB interface {
	ListRegistered(ctx context.Context) []string
	IsRegistered(bankCode string) bool
	RouteKey(bankCode string) string
}
