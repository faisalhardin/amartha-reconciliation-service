package bankregistry

import (
	"context"
	"sort"
	"strings"

	bankregistryrepo "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/repo/bankregistry"
)

var registeredBanks = map[string]struct{}{
	"BNI":     {},
	"BRI":     {},
	"MANDIRI": {},
	"BCA":     {},
}

type bankRegistryDB struct{}

func New() bankregistryrepo.BankRegistryDB {
	return &bankRegistryDB{}
}

func (d *bankRegistryDB) ListRegistered(_ context.Context) []string {
	out := make([]string, 0, len(registeredBanks))
	for code := range registeredBanks {
		out = append(out, code)
	}
	sort.Strings(out)
	return out
}

func (d *bankRegistryDB) IsRegistered(bankCode string) bool {
	_, ok := registeredBanks[normalize(bankCode)]
	return ok
}

func (d *bankRegistryDB) RouteKey(bankCode string) string {
	code := normalize(bankCode)
	if _, ok := registeredBanks[code]; ok {
		return code
	}
	return bankregistryrepo.CommonRouteKey
}

func normalize(bankCode string) string {
	return strings.ToUpper(strings.TrimSpace(bankCode))
}
