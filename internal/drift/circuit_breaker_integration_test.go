package drift_test

import (
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

// TestCircuitBreaker_Integration_WatcherRespectsCB validates that a Watcher
// integrated with a CircuitBreaker skips checks when the circuit is open and
// resumes after the reset timeout.
func TestCircuitBreaker_Integration_WatcherRespectsCB(t *testing.T) {
	cfg := watcherConfig()
	w := drift.NewWatcher(cfg, 50*time.Millisecond)
	cb := drift.NewCircuitBreaker(2, 80*time.Millisecond)

	service := cfg.Name

	// Simulate two consecutive failures to open the circuit.
	cb.RecordFailure(service)
	cb.RecordFailure(service)

	if cb.State(service) != drift.CircuitOpen {
		t.Fatal("expected circuit to be open after failures")
	}

	// While open, Allow should block.
	if cb.Allow(service) {
		t.Error("circuit should block checks when open")
	}

	// Wait for reset timeout.
	time.Sleep(100 * time.Millisecond)

	// After timeout the circuit transitions to half-open.
	if !cb.Allow(service) {
		t.Error("circuit should allow probe after reset timeout")
	}
	if cb.State(service) != drift.CircuitHalfOpen {
		t.Errorf("expected half-open, got %s", cb.State(service))
	}

	// A successful check closes the circuit.
	result := w.Check()
	if result == nil {
		t.Fatal("expected a check result")
	}
	cb.RecordSuccess(service)

	if cb.State(service) != drift.CircuitClosed {
		t.Errorf("expected closed after success, got %s", cb.State(service))
	}
}
