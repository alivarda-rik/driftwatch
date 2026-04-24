package drift

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeSnapshot(name string) *Snapshot {
	return &Snapshot{
		ServiceName: name,
		CapturedAt:  time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
		Fields: map[string]interface{}{
			"replicas": float64(3),
			"image":    "nginx:1.25",
		},
	}
}

func TestSnapshotStore_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSnapshotStore(dir)
	if err != nil {
		t.Fatalf("NewSnapshotStore: %v", err)
	}

	snap := makeSnapshot("web-api")
	if err := store.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := store.Load("web-api")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.ServiceName != snap.ServiceName {
		t.Errorf("ServiceName: got %q, want %q", got.ServiceName, snap.ServiceName)
	}
	if got.Fields["image"] != snap.Fields["image"] {
		t.Errorf("Fields[image]: got %v, want %v", got.Fields["image"], snap.Fields["image"])
	}
}

func TestSnapshotStore_LoadMissing(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewSnapshotStore(dir)
	_, err := store.Load("nonexistent")
	if err == nil {
		t.Error("expected error loading nonexistent snapshot, got nil")
	}
}

func TestSnapshotStore_SaveNil(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewSnapshotStore(dir)
	if err := store.Save(nil); err == nil {
		t.Error("expected error saving nil snapshot, got nil")
	}
}

func TestSnapshotStore_FileCreated(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewSnapshotStore(dir)
	snap := makeSnapshot("auth-service")
	_ = store.Save(snap)

	expected := filepath.Join(dir, "auth-service.json")
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Errorf("expected snapshot file at %s, not found", expected)
	}
}

func TestNewSnapshotStore_CreatesDir(t *testing.T) {
	base := t.TempDir()
	nested := filepath.Join(base, "snapshots", "v1")
	_, err := NewSnapshotStore(nested)
	if err != nil {
		t.Fatalf("NewSnapshotStore: %v", err)
	}
	if _, err := os.Stat(nested); os.IsNotExist(err) {
		t.Error("expected nested directory to be created")
	}
}
