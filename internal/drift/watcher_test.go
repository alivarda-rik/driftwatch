package drift

import (
	"testing"
	"time"

	"github.com/user/driftwatch/internal/config"
)

func watcherConfig() *config.ServiceConfig {
	return &config.ServiceConfig{
		Name: "svc-watch",
		Fields: map[string]string{
			"env":     "production",
			"region":  "us-east-1",
		},
	}
}

func TestWatcher_Check_NoDrift(t *testing.T) {
	detector := NewDetector()
	manager := NewBaselineManager(NewBaselineStore(t.TempDir()))
	w := NewWatcher(detector, manager, time.Second)

	cfg := watcherConfig()
	live := map[string]string{
		"env":    "production",
		"region": "us-east-1",
	}

	result, err := w.Check(cfg, live)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Drifted {
		t.Errorf("expected no drift, got %d entries", len(result.Entries))
	}
	if result.Service != "svc-watch" {
		t.Errorf("expected service svc-watch, got %s", result.Service)
	}
}

func TestWatcher_Check_WithDrift(t *testing.T) {
	detector := NewDetector()
	manager := NewBaselineManager(NewBaselineStore(t.TempDir()))
	w := NewWatcher(detector, manager, time.Second)

	cfg := watcherConfig()
	live := map[string]string{
		"env":    "staging",
		"region": "us-east-1",
	}

	result, err := w.Check(cfg, live)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Drifted {
		t.Error("expected drift but got none")
	}
	if len(result.Entries) == 0 {
		t.Error("expected at least one diff entry")
	}
}

func TestWatcher_Check_NilConfig(t *testing.T) {
	detector := NewDetector()
	manager := NewBaselineManager(NewBaselineStore(t.TempDir()))
	w := NewWatcher(detector, manager, time.Second)

	_, err := w.Check(nil, map[string]string{})
	if err == nil {
		t.Error("expected error for nil config")
	}
}

func TestWatcher_Watch_ReceivesResult(t *testing.T) {
	detector := NewDetector()
	manager := NewBaselineManager(NewBaselineStore(t.TempDir()))
	w := NewWatcher(detector, manager, 50*time.Millisecond)

	cfg := watcherConfig()
	live := map[string]string{"env": "production", "region": "us-east-1"}

	ch := w.Watch(cfg, live)
	select {
	case result := <-ch:
		if result == nil {
			t.Fatal("received nil result")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for watch result")
	}
	w.Stop()
}
