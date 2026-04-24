package drift_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/driftwatch/internal/drift"
)

func makeBaseline(name string) *drift.Baseline {
	return &drift.Baseline{
		ServiceName: name,
		CapturedAt:  time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Fields: map[string]string{
			"replicas": "3",
			"image":    "nginx:1.25",
		},
	}
}

func TestBaselineStore_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store := drift.NewBaselineStore(dir)
	b := makeBaseline("svc-alpha")

	if err := store.Save(b); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := store.Load("svc-alpha")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.ServiceName != b.ServiceName {
		t.Errorf("ServiceName: got %q, want %q", got.ServiceName, b.ServiceName)
	}
	if got.Fields["replicas"] != "3" {
		t.Errorf("Fields[replicas]: got %q, want %q", got.Fields["replicas"], "3")
	}
}

func TestBaselineStore_LoadMissing(t *testing.T) {
	store := drift.NewBaselineStore(t.TempDir())
	_, err := store.Load("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing baseline, got nil")
	}
}

func TestBaselineStore_SaveNil(t *testing.T) {
	store := drift.NewBaselineStore(t.TempDir())
	if err := store.Save(nil); err == nil {
		t.Fatal("expected error saving nil baseline")
	}
}

func TestBaselineStore_Delete(t *testing.T) {
	dir := t.TempDir()
	store := drift.NewBaselineStore(dir)
	b := makeBaseline("svc-beta")

	_ = store.Save(b)
	if err := store.Delete("svc-beta"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	path := filepath.Join(dir, "svc-beta.baseline.json")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected baseline file to be deleted")
	}
}

func TestBaselineStore_DeleteMissing(t *testing.T) {
	store := drift.NewBaselineStore(t.TempDir())
	if err := store.Delete("ghost-service"); err != nil {
		t.Errorf("Delete of missing baseline should not error: %v", err)
	}
}
