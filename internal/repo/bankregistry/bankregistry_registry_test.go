package bankregistry

import (
	"context"
	"testing"

	bankregistryrepo "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/repo/bankregistry"
)

func TestListRegistered(t *testing.T) {
	reg := New()
	got := reg.ListRegistered(context.Background())
	want := []string{"BCA", "BNI", "BRI", "MANDIRI"}
	if len(got) != len(want) {
		t.Fatalf("ListRegistered len = %d, want %d", len(got), len(want))
	}
	for i, code := range want {
		if got[i] != code {
			t.Fatalf("ListRegistered[%d] = %q, want %q", i, got[i], code)
		}
	}
}

func TestIsRegistered(t *testing.T) {
	reg := New()
	if !reg.IsRegistered("bca") {
		t.Fatal("expected BCA to be registered")
	}
	if reg.IsRegistered("UNKNOWN") {
		t.Fatal("expected UNKNOWN to be unregistered")
	}
}

func TestRouteKey(t *testing.T) {
	reg := New()
	if got := reg.RouteKey("BCA"); got != "BCA" {
		t.Fatalf("RouteKey(BCA) = %q, want BCA", got)
	}
	if got := reg.RouteKey("UNKNOWN"); got != bankregistryrepo.CommonRouteKey {
		t.Fatalf("RouteKey(UNKNOWN) = %q, want %s", got, bankregistryrepo.CommonRouteKey)
	}
}
