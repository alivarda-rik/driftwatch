package drift_test

import (
	"testing"
	"time"

	"github.com/user/driftwatch/internal/config"
	"github.com/user/driftwatch/internal/drift"
)

func TestWatcherIntegration_DriftDetectedOverTime(t *testing.T) {
	detector := drift.NewDetector()
	store := drift.NewBaselineStore(t.TempDir())
	manager := drift.NewBaselineManager(store)
	watcher := drift.NewWatcher(detector, manager, 40*time.Millisecond)

	cfg := &config.ServiceConfig{
		Name: "integration-svc",
		Fields: map[string]string{
			"log_level": "info",
			"timeout":   "30s",
		},
	}

	// live state that drifts from declared
	live := map[string]string{
		"log_level": "debug",
		"timeout":   "30s",
	}

	ch := watcher.Watch(cfg, live)
	defer watcher.Stop()

	var driftSeen bool
	timeout := time.After(400 * time.Millisecond)
	for {
		select {
		case result, ok := <-ch:
			if !ok {
				return
			}
			if result.Drifted {
				driftSeen = true
				watcher.Stop()
			}
		case <-timeout:
			if !driftSeen {
				t.Error("expected drift to be detected within timeout")
			}
			return
		}
	}
}
