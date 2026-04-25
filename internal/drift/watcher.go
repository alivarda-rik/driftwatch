package drift

import (
	"fmt"
	"time"

	"github.com/user/driftwatch/internal/config"
)

// WatchResult holds the outcome of a single watch cycle.
type WatchResult struct {
	Service   string
	Drifted   bool
	Entries   []DiffEntry
	CheckedAt time.Time
}

// Watcher periodically checks for drift between a declared config and live state.
type Watcher struct {
	detector *Detector
	manager  *BaselineManager
	interval time.Duration
	stop     chan struct{}
}

// NewWatcher creates a Watcher with the given detector, baseline manager, and poll interval.
func NewWatcher(detector *Detector, manager *BaselineManager, interval time.Duration) *Watcher {
	return &Watcher{
		detector: detector,
		manager:  manager,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Check performs a single drift check for the given service config and live state.
func (w *Watcher) Check(cfg *config.ServiceConfig, live map[string]string) (*WatchResult, error) {
	if cfg == nil {
		return nil, fmt.Errorf("watcher: service config must not be nil")
	}

	report, err := w.detector.Detect(cfg, live)
	if err != nil {
		return nil, fmt.Errorf("watcher: detect failed: %w", err)
	}

	return &WatchResult{
		Service:   cfg.Name,
		Drifted:   len(report.Diffs) > 0,
		Entries:   report.Diffs,
		CheckedAt: time.Now().UTC(),
	}, nil
}

// Watch runs drift checks on the given interval, sending results to the returned channel.
// Call Stop to terminate the loop.
func (w *Watcher) Watch(cfg *config.ServiceConfig, live map[string]string) <-chan *WatchResult {
	results := make(chan *WatchResult, 1)
	go func() {
		defer close(results)
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				result, err := w.Check(cfg, live)
				if err == nil {
					results <- result
				}
			case <-w.stop:
				return
			}
		}
	}()
	return results
}

// Stop signals the Watch goroutine to exit.
// It is safe to call Stop only once; subsequent calls will panic on a closed channel.
func (w *Watcher) Stop() {
	close(w.stop)
}

// IsStopped reports whether the watcher has been stopped.
func (w *Watcher) IsStopped() bool {
	select {
	case <-w.stop:
		return true
	default:
		return false
	}
}
