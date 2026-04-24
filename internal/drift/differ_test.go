package drift

import (
	"strings"
	"testing"
)

func TestDiff_NoDifferences(t *testing.T) {
	declared := map[string]interface{}{"port": 8080, "host": "localhost"}
	live := map[string]interface{}{"port": 8080, "host": "localhost"}

	entries := Diff(declared, live)
	if len(entries) != 0 {
		t.Fatalf("expected no diff entries, got %d", len(entries))
	}
}

func TestDiff_ChangedValue(t *testing.T) {
	declared := map[string]interface{}{"port": 8080}
	live := map[string]interface{}{"port": 9090}

	entries := Diff(declared, live)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Kind != DiffKindChanged {
		t.Errorf("expected DiffKindChanged, got %s", entries[0].Kind)
	}
	if entries[0].Key != "port" {
		t.Errorf("unexpected key: %s", entries[0].Key)
	}
}

func TestDiff_MissingInLive(t *testing.T) {
	declared := map[string]interface{}{"replicas": 3}
	live := map[string]interface{}{}

	entries := Diff(declared, live)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Kind != DiffKindMissing {
		t.Errorf("expected DiffKindMissing, got %s", entries[0].Kind)
	}
}

func TestDiff_ExtraInLive(t *testing.T) {
	declared := map[string]interface{}{}
	live := map[string]interface{}{"debug": true}

	entries := Diff(declared, live)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Kind != DiffKindExtra {
		t.Errorf("expected DiffKindExtra, got %s", entries[0].Kind)
	}
}

func TestDiffEntry_String(t *testing.T) {
	tests := []struct {
		entry    DiffEntry
		contains string
	}{
		{DiffEntry{Key: "port", Declared: 8080, Live: 9090, Kind: DiffKindChanged}, "changed"},
		{DiffEntry{Key: "host", Declared: "localhost", Kind: DiffKindMissing}, "missing"},
		{DiffEntry{Key: "debug", Live: true, Kind: DiffKindExtra}, "extra"},
	}

	for _, tc := range tests {
		s := tc.entry.String()
		if !strings.Contains(s, tc.contains) {
			t.Errorf("String() = %q, want it to contain %q", s, tc.contains)
		}
		if !strings.Contains(s, tc.entry.Key) {
			t.Errorf("String() = %q, want it to contain key %q", s, tc.entry.Key)
		}
	}
}
