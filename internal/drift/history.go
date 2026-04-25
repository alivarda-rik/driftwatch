package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// HistoryEntry records a drift check result at a point in time.
type HistoryEntry struct {
	Timestamp time.Time   `json:"timestamp"`
	Service   string      `json:"service"`
	Drifted   bool        `json:"drifted"`
	Entries   []DiffEntry `json:"entries,omitempty"`
}

// HistoryStore persists and retrieves drift history for services.
type HistoryStore struct {
	dir string
}

// NewHistoryStore creates a HistoryStore that writes to the given directory.
func NewHistoryStore(dir string) *HistoryStore {
	return &HistoryStore{dir: dir}
}

func (h *HistoryStore) historyPath(service string) string {
	return filepath.Join(h.dir, fmt.Sprintf("%s.history.json", service))
}

// Append adds a new HistoryEntry for the given service.
func (h *HistoryStore) Append(service string, entry HistoryEntry) error {
	if service == "" {
		return fmt.Errorf("service name must not be empty")
	}
	entries, _ := h.Load(service)
	entries = append(entries, entry)

	if err := os.MkdirAll(h.dir, 0755); err != nil {
		return fmt.Errorf("creating history dir: %w", err)
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling history: %w", err)
	}
	return os.WriteFile(h.historyPath(service), data, 0644)
}

// Load returns all history entries for the given service.
func (h *HistoryStore) Load(service string) ([]HistoryEntry, error) {
	data, err := os.ReadFile(h.historyPath(service))
	if os.IsNotExist(err) {
		return []HistoryEntry{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading history: %w", err)
	}
	var entries []HistoryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("parsing history: %w", err)
	}
	return entries, nil
}

// Recent returns the most recent n entries for the given service, newest first.
func (h *HistoryStore) Recent(service string, n int) ([]HistoryEntry, error) {
	entries, err := h.Load(service)
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.After(entries[j].Timestamp)
	})
	if n > 0 && len(entries) > n {
		entries = entries[:n]
	}
	return entries, nil
}
