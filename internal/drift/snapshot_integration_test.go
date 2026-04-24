package drift_test

import (
	"testing"
	"time"

	"github.com/user/driftwatch/internal/drift"
)

// TestSnapshotRoundtrip verifies that a snapshot saved and reloaded
// produces identical data that can be fed into the Detector.
func TestSnapshotRoundtrip(t *testing.T) {
	dir := t.TempDir()
	store, err := drift.NewSnapshotStore(dir)
	if err != nil {
		t.Fatalf("NewSnapshotStore: %v", err)
	}

	original := &drift.Snapshot{
		ServiceName: "payment-svc",
		CapturedAt:  time.Now().UTC().Truncate(time.Second),
		Fields: map[string]interface{}{
			"replicas":    float64(2),
			"image":       "payments:v3.1",
			"port":        float64(8080),
			"healthcheck": "/health",
		},
	}

	if err := store.Save(original); err != nil {
		t.Fatalf("Save: %v", err)
	}

	reloaded, err := store.Load("payment-svc")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if reloaded.ServiceName != original.ServiceName {
		t.Errorf("ServiceName mismatch: got %q want %q", reloaded.ServiceName, original.ServiceName)
	}

	for k, want := range original.Fields {
		got, ok := reloaded.Fields[k]
		if !ok {
			t.Errorf("field %q missing after roundtrip", k)
			continue
		}
		if got != want {
			t.Errorf("field %q: got %v (%T), want %v (%T)", k, got, got, want, want)
		}
	}
}
