package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a point-in-time capture of a service's live state.
type Snapshot struct {
	ServiceName string                 `json:"service_name"`
	CapturedAt  time.Time              `json:"captured_at"`
	Fields      map[string]interface{} `json:"fields"`
}

// SnapshotStore handles persistence of snapshots to disk.
type SnapshotStore struct {
	Dir string
}

// NewSnapshotStore creates a SnapshotStore rooted at dir.
func NewSnapshotStore(dir string) (*SnapshotStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("snapshot store: create dir: %w", err)
	}
	return &SnapshotStore{Dir: dir}, nil
}

// Save writes a snapshot to disk as <service_name>.json.
func (s *SnapshotStore) Save(snap *Snapshot) error {
	if snap == nil {
		return fmt.Errorf("snapshot store: cannot save nil snapshot")
	}
	path := filepath.Join(s.Dir, snap.ServiceName+".json")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot store: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return fmt.Errorf("snapshot store: encode: %w", err)
	}
	return nil
}

// Load reads a previously saved snapshot for the given service name.
func (s *SnapshotStore) Load(serviceName string) (*Snapshot, error) {
	path := filepath.Join(s.Dir, serviceName+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot store: read file: %w", err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot store: decode: %w", err)
	}
	return &snap, nil
}
