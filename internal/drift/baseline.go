package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Baseline represents a saved reference point for a service's configuration.
type Baseline struct {
	ServiceName string            `json:"service_name"`
	CapturedAt  time.Time         `json:"captured_at"`
	Fields      map[string]string `json:"fields"`
}

// BaselineStore manages persistence of baselines to disk.
type BaselineStore struct {
	dir string
}

// NewBaselineStore creates a BaselineStore that persists baselines under dir.
func NewBaselineStore(dir string) *BaselineStore {
	return &BaselineStore{dir: dir}
}

// Save writes a baseline to disk as JSON.
func (s *BaselineStore) Save(b *Baseline) error {
	if b == nil {
		return fmt.Errorf("baseline must not be nil")
	}
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return fmt.Errorf("creating baseline dir: %w", err)
	}
	path := filepath.Join(s.dir, b.ServiceName+".baseline.json")
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling baseline: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// Load reads a baseline for the given service name from disk.
func (s *BaselineStore) Load(serviceName string) (*Baseline, error) {
	path := filepath.Join(s.dir, serviceName+".baseline.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no baseline found for service %q", serviceName)
		}
		return nil, fmt.Errorf("reading baseline: %w", err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("unmarshalling baseline: %w", err)
	}
	return &b, nil
}

// Delete removes the baseline file for the given service name.
func (s *BaselineStore) Delete(serviceName string) error {
	path := filepath.Join(s.dir, serviceName+".baseline.json")
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("deleting baseline: %w", err)
	}
	return nil
}
