package messaging

import (
	"context"
	"testing"
	"time"

	bankregistryrepo "github.com/faisalhardin/amartha-reconciliation-service/internal/entity/repo/bankregistry"
	bankregistryimpl "github.com/faisalhardin/amartha-reconciliation-service/internal/repo/bankregistry"
)

func TestJobQueueHubPublishProcessRoutesByBank(t *testing.T) {
	reg := bankregistryimpl.New()
	hub := NewJobQueueHub(reg, 10)
	defer hub.Close()

	hub.PublishProcess("file-bni", "BNI")
	hub.PublishProcess("file-unknown", "UNKNOWN")

	select {
	case id := <-hub.ProcessQueues()["BNI"].Channel():
		if id != "file-bni" {
			t.Fatalf("BNI queue got %q, want file-bni", id)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting on BNI process queue")
	}

	select {
	case id := <-hub.ProcessQueues()[bankregistryrepo.CommonRouteKey].Channel():
		if id != "file-unknown" {
			t.Fatalf("COMMON queue got %q, want file-unknown", id)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting on COMMON process queue")
	}
}

func TestJobQueueHubPublishReconciliationRoutesByBank(t *testing.T) {
	reg := bankregistryimpl.New()
	hub := NewJobQueueHub(reg, 10)
	defer hub.Close()

	hub.PublishReconciliation("file-bca", "BCA")
	hub.PublishReconciliation("file-unknown", "UNKNOWN")

	select {
	case id := <-hub.ReconciliationQueues()["BCA"].Channel():
		if id != "file-bca" {
			t.Fatalf("BCA queue got %q, want file-bca", id)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting on BCA reconciliation queue")
	}

	select {
	case id := <-hub.ReconciliationQueues()[bankregistryrepo.CommonRouteKey].Channel():
		if id != "file-unknown" {
			t.Fatalf("COMMON queue got %q, want file-unknown", id)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting on COMMON reconciliation queue")
	}
}

func TestJobQueueHubHasAllRouteKeys(t *testing.T) {
	reg := bankregistryimpl.New()
	hub := NewJobQueueHub(reg, 10)
	defer hub.Close()

	wantKeys := append(reg.ListRegistered(context.Background()), bankregistryrepo.CommonRouteKey)
	if len(hub.ProcessQueues()) != len(wantKeys) {
		t.Fatalf("process queues count = %d, want %d", len(hub.ProcessQueues()), len(wantKeys))
	}
	for _, key := range wantKeys {
		if _, ok := hub.ProcessQueues()[key]; !ok {
			t.Fatalf("missing process queue for %q", key)
		}
		if _, ok := hub.ReconciliationQueues()[key]; !ok {
			t.Fatalf("missing reconciliation queue for %q", key)
		}
	}
}
