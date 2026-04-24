package drift

import (
	"fmt"
	"time"

	"github.com/user/driftwatch/internal/config"
)

// BaselineManager orchestrates capturing and comparing baselines.
type BaselineManager struct {
	store    *BaselineStore
	detector *Detector
}

// NewBaselineManager creates a BaselineManager using the given store directory.
func NewBaselineManager(storeDir string, detector *Detector) *BaselineManager {
	return &BaselineManager{
		store:    NewBaselineStore(storeDir),
		detector: detector,
	}
}

// Capture records the current declared config as a baseline for later comparison.
func (m *BaselineManager) Capture(cfg *config.ServiceConfig) error {
	if cfg == nil {
		return fmt.Errorf("cannot capture nil config")
	}
	fields := flattenDeclared(cfg)
	b := &Baseline{
		ServiceName: cfg.Name,
		CapturedAt:  time.Now().UTC(),
		Fields:      fields,
	}
	return m.store.Save(b)
}

// CompareToBaseline loads the stored baseline for cfg.Name and diffs it
// against the provided live fields, returning any drift entries.
func (m *BaselineManager) CompareToBaseline(cfg *config.ServiceConfig, live map[string]string) ([]DiffEntry, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config must not be nil")
	}
	b, err := m.store.Load(cfg.Name)
	if err != nil {
		return nil, fmt.Errorf("loading baseline: %w", err)
	}
	entries := Diff(b.Fields, live)
	return entries, nil
}

// DeleteBaseline removes the stored baseline for the named service.
func (m *BaselineManager) DeleteBaseline(serviceName string) error {
	return m.store.Delete(serviceName)
}
