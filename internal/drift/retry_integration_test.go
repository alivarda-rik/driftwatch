package drift_test

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/driftwatch/internal/drift"
)

// TestRetryer_Integration_WatcherUsesRetry verifies that a Retryer integrates
// correctly with a live check function that eventually succeeds.
func TestRetryer_Integration_WatcherUsesRetry(t *testing.T) {
	policy := drift.RetryPolicy{
		MaxAttempts: 4,
		Delay:       10 * time.Millisecond,
		Multiplier:  1.5,
	}
	r := drift.NewRetryer(policy)

	var callCount int32
	err := r.Do(func() error {
		n := atomic.AddInt32(&callCount, 1)
		if n < 3 {
			return errors.New("not ready")
		}
		return nil
	})

	if err != nil {
		t.Fatalf("expected eventual success, got: %v", err)
	}
	if callCount < 3 {
		t.Fatalf("expected at least 3 calls, got %d", callCount)
	}
}

// TestRetryer_Integration_AlwaysFails ensures the retryer surfaces the last
// error after exhausting all attempts.
func TestRetryer_Integration_AlwaysFails(t *testing.T) {
	policy := drift.RetryPolicy{
		MaxAttempts: 3,
		Delay:       5 * time.Millisecond,
		Multiplier:  1.0,
	}
	r := drift.NewRetryer(policy)

	expected := errors.New("service unavailable")
	var calls int32
	err := r.Do(func() error {
		atomic.AddInt32(&calls, 1)
		return expected
	})

	if !errors.Is(err, expected) {
		t.Fatalf("expected %v, got %v", expected, err)
	}
	if calls != 3 {
		t.Fatalf("expected exactly 3 calls, got %d", calls)
	}
}
