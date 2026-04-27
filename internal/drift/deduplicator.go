package drift

import (
	"fmt"
	"sync"
	"time"
)

// DedupKey uniquely identifies a drift event for deduplication purposes.
type DedupKey struct {
	Service string
	Key     string
	Kind    string
}

// DedupEntry holds metadata about a previously seen drift event.
type DedupEntry struct {
	FirstSeen time.Time
	LastSeen  time.Time
	Count     int
}

// Deduplicator suppresses repeated identical drift events within a TTL window.
type Deduplicator struct {
	mu      sync.Mutex
	ttl     time.Duration
	seen    map[DedupKey]*DedupEntry
	nowFunc func() time.Time
}

// NewDeduplicator creates a Deduplicator that suppresses repeated events within ttl.
func NewDeduplicator(ttl time.Duration) *Deduplicator {
	return &Deduplicator{
		ttl:     ttl,
		seen:    make(map[DedupKey]*DedupEntry),
		nowFunc: time.Now,
	}
}

// IsDuplicate returns true if the given DiffEntry for a service was already
// seen within the TTL window. It also records the event if it is new.
func (d *Deduplicator) IsDuplicate(service string, entry DiffEntry) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.nowFunc()
	key := DedupKey{
		Service: service,
		Key:     entry.Key,
		Kind:    fmt.Sprintf("%T", entry),
	}

	if e, ok := d.seen[key]; ok {
		if now.Sub(e.LastSeen) < d.ttl {
			e.LastSeen = now
			e.Count++
			return true
		}
	}

	d.seen[key] = &DedupEntry{
		FirstSeen: now,
		LastSeen:  now,
		Count:     1,
	}
	return false
}

// Evict removes all entries whose LastSeen time is older than the TTL.
func (d *Deduplicator) Evict() int {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.nowFunc()
	removed := 0
	for k, e := range d.seen {
		if now.Sub(e.LastSeen) >= d.ttl {
			delete(d.seen, k)
			removed++
		}
	}
	return removed
}

// Size returns the number of tracked dedup entries.
func (d *Deduplicator) Size() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.seen)
}
