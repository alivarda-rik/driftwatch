package drift

import (
	"testing"
	"time"
)

func makeDiffEntry(key string) DiffEntry {
	return DiffEntry{
		Key:      key,
		Declared: "v1",
		Live:     "v2",
	}
}

func TestDeduplicator_FirstCallNotDuplicate(t *testing.T) {
	d := NewDeduplicator(5 * time.Second)
	entry := makeDiffEntry("port")
	if d.IsDuplicate("svc-a", entry) {
		t.Error("expected first call to not be a duplicate")
	}
}

func TestDeduplicator_SecondCallWithinTTLIsDuplicate(t *testing.T) {
	d := NewDeduplicator(5 * time.Second)
	entry := makeDiffEntry("port")
	d.IsDuplicate("svc-a", entry)
	if !d.IsDuplicate("svc-a", entry) {
		t.Error("expected second call within TTL to be a duplicate")
	}
}

func TestDeduplicator_AfterTTLExpiry_NotDuplicate(t *testing.T) {
	now := time.Now()
	d := NewDeduplicator(1 * time.Second)
	d.nowFunc = func() time.Time { return now }

	entry := makeDiffEntry("replicas")
	d.IsDuplicate("svc-b", entry)

	// Advance time past TTL
	d.nowFunc = func() time.Time { return now.Add(2 * time.Second) }

	if d.IsDuplicate("svc-b", entry) {
		t.Error("expected call after TTL expiry to not be a duplicate")
	}
}

func TestDeduplicator_DifferentServicesAreIndependent(t *testing.T) {
	d := NewDeduplicator(5 * time.Second)
	entry := makeDiffEntry("image")

	d.IsDuplicate("svc-a", entry)

	if d.IsDuplicate("svc-b", entry) {
		t.Error("expected different service to not be treated as duplicate")
	}
}

func TestDeduplicator_Evict_RemovesExpiredEntries(t *testing.T) {
	now := time.Now()
	d := NewDeduplicator(1 * time.Second)
	d.nowFunc = func() time.Time { return now }

	d.IsDuplicate("svc-a", makeDiffEntry("cpu"))
	d.IsDuplicate("svc-b", makeDiffEntry("memory"))

	if d.Size() != 2 {
		t.Fatalf("expected 2 entries, got %d", d.Size())
	}

	d.nowFunc = func() time.Time { return now.Add(2 * time.Second) }
	removed := d.Evict()

	if removed != 2 {
		t.Errorf("expected 2 evictions, got %d", removed)
	}
	if d.Size() != 0 {
		t.Errorf("expected 0 entries after eviction, got %d", d.Size())
	}
}

func TestDeduplicator_Evict_KeepsActiveEntries(t *testing.T) {
	now := time.Now()
	d := NewDeduplicator(10 * time.Second)
	d.nowFunc = func() time.Time { return now }

	d.IsDuplicate("svc-a", makeDiffEntry("port"))

	d.nowFunc = func() time.Time { return now.Add(2 * time.Second) }
	removed := d.Evict()

	if removed != 0 {
		t.Errorf("expected 0 evictions, got %d", removed)
	}
	if d.Size() != 1 {
		t.Errorf("expected 1 active entry, got %d", d.Size())
	}
}
