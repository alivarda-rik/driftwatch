package drift_test

import (
	"testing"

	"github.com/user/driftwatch/internal/config"
	"github.com/user/driftwatch/internal/drift"
)

func baseServiceConfig() *config.ServiceConfig {
	return &config.ServiceConfig{
		Name: "web",
		Settings: map[string]interface{}{
			"replicas": 2,
			"image":    "nginx:1.25",
		},
	}
}

func TestBaselineManager_CaptureAndCompare(t *testing.T) {
	dir := t.TempDir()
	detector := drift.NewDetector()
	mgr := drift.NewBaselineManager(dir, detector)
	cfg := baseServiceConfig()

	if err := mgr.Capture(cfg); err != nil {
		t.Fatalf("Capture: %v", err)
	}

	live := map[string]string{
		"replicas": "2",
		"image":    "nginx:1.25",
	}
	entries, err := mgr.CompareToBaseline(cfg, live)
	if err != nil {
		t.Fatalf("CompareToBaseline: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected no drift, got %d entries", len(entries))
	}
}

func TestBaselineManager_CaptureAndCompare_WithDrift(t *testing.T) {
	dir := t.TempDir()
	mgr := drift.NewBaselineManager(dir, drift.NewDetector())
	cfg := baseServiceConfig()

	_ = mgr.Capture(cfg)

	live := map[string]string{
		"replicas": "5",
		"image":    "nginx:1.25",
	}
	entries, err := mgr.CompareToBaseline(cfg, live)
	if err != nil {
		t.Fatalf("CompareToBaseline: %v", err)
	}
	if len(entries) == 0 {
		t.Error("expected drift entries, got none")
	}
}

func TestBaselineManager_CaptureNil(t *testing.T) {
	mgr := drift.NewBaselineManager(t.TempDir(), drift.NewDetector())
	if err := mgr.Capture(nil); err == nil {
		t.Error("expected error capturing nil config")
	}
}

func TestBaselineManager_CompareNoBaseline(t *testing.T) {
	mgr := drift.NewBaselineManager(t.TempDir(), drift.NewDetector())
	cfg := baseServiceConfig()
	_, err := mgr.CompareToBaseline(cfg, map[string]string{})
	if err == nil {
		t.Error("expected error when no baseline exists")
	}
}
