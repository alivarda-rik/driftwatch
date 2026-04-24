package drift

import (
	"testing"
)

func sampleEntries() []DiffEntry {
	return []DiffEntry{
		{Key: "app.port", Declared: "8080", Live: "9090"},
		{Key: "app.name", Declared: "svc", Live: ""},
		{Key: "app.debug", Declared: "", Live: "true"},
		{Key: "db.host", Declared: "localhost", Live: "remotehost"},
	}
}

func TestFilter_NoOptions(t *testing.T) {
	entries := sampleEntries()
	got := Filter(entries, FilterOptions{})
	if len(got) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(got))
	}
}

func TestFilter_KeyPrefix(t *testing.T) {
	got := Filter(sampleEntries(), FilterOptions{KeyPrefix: "app."})
	if len(got) != 3 {
		t.Fatalf("expected 3 entries with prefix 'app.', got %d", len(got))
	}
	for _, e := range got {
		if e.Key[:4] != "app." {
			t.Errorf("unexpected key %q", e.Key)
		}
	}
}

func TestFilter_OnlyMissing(t *testing.T) {
	got := Filter(sampleEntries(), FilterOptions{OnlyMissing: true})
	if len(got) != 1 {
		t.Fatalf("expected 1 missing entry, got %d", len(got))
	}
	if got[0].Key != "app.name" {
		t.Errorf("expected key 'app.name', got %q", got[0].Key)
	}
}

func TestFilter_OnlyExtra(t *testing.T) {
	got := Filter(sampleEntries(), FilterOptions{OnlyExtra: true})
	if len(got) != 1 {
		t.Fatalf("expected 1 extra entry, got %d", len(got))
	}
	if got[0].Key != "app.debug" {
		t.Errorf("expected key 'app.debug', got %q", got[0].Key)
	}
}

func TestFilter_EmptyInput(t *testing.T) {
	got := Filter(nil, FilterOptions{KeyPrefix: "app."})
	if got != nil && len(got) != 0 {
		t.Errorf("expected empty result for nil input, got %v", got)
	}
}

func TestFilterByKeys(t *testing.T) {
	keys := []string{"app.port", "db.host"}
	got := FilterByKeys(sampleEntries(), keys)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestFilterByKeys_NoMatch(t *testing.T) {
	got := FilterByKeys(sampleEntries(), []string{"nonexistent"})
	if len(got) != 0 {
		t.Errorf("expected 0 entries, got %d", len(got))
	}
}
