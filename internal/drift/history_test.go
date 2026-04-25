package drift

import (
	"os"
	"testing"
	"time"
)

func makeHistoryEntry(service string, drifted bool) HistoryEntry {
	return HistoryEntry{
		Timestamp: time.Now().UTC(),
		Service:   service,
		Drifted:   drifted,
		Entries: []DiffEntry{
			{Key: "port", Declared: "8080", Live: "9090"},
		},
	}
}

func TestHistoryStore_AppendAndLoad(t *testing.T) {
	dir := t.TempDir()
	hs := NewHistoryStore(dir)

	e1 := makeHistoryEntry("svc-a", true)
	e2 := makeHistoryEntry("svc-a", false)

	if err := hs.Append("svc-a", e1); err != nil {
		t.Fatalf("first append: %v", err)
	}
	if err := hs.Append("svc-a", e2); err != nil {
		t.Fatalf("second append: %v", err)
	}

	entries, err := hs.Load("svc-a")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestHistoryStore_LoadMissing(t *testing.T) {
	hs := NewHistoryStore(t.TempDir())
	entries, err := hs.Load("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestHistoryStore_AppendEmptyService(t *testing.T) {
	hs := NewHistoryStore(t.TempDir())
	err := hs.Append("", makeHistoryEntry("", false))
	if err == nil {
		t.Error("expected error for empty service name")
	}
}

func TestHistoryStore_Recent(t *testing.T) {
	dir := t.TempDir()
	hs := NewHistoryStore(dir)

	for i := 0; i < 5; i++ {
		e := makeHistoryEntry("svc-b", i%2 == 0)
		e.Timestamp = time.Now().Add(time.Duration(i) * time.Second)
		if err := hs.Append("svc-b", e); err != nil {
			t.Fatalf("append %d: %v", i, err)
		}
	}

	recent, err := hs.Recent("svc-b", 3)
	if err != nil {
		t.Fatalf("recent: %v", err)
	}
	if len(recent) != 3 {
		t.Errorf("expected 3 recent entries, got %d", len(recent))
	}
	if !recent[0].Timestamp.After(recent[1].Timestamp) {
		t.Error("expected entries sorted newest first")
	}
}

func TestHistoryStore_FileCreated(t *testing.T) {
	dir := t.TempDir()
	hs := NewHistoryStore(dir)

	if err := hs.Append("svc-c", makeHistoryEntry("svc-c", false)); err != nil {
		t.Fatalf("append: %v", err)
	}
	path := dir + "/svc-c.history.json"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected history file at %s", path)
	}
}
