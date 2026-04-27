package drift

import (
	"errors"
	"testing"
	"time"
)

func TestRetryer_SucceedsOnFirstAttempt(t *testing.T) {
	r := NewRetryer(DefaultRetryPolicy())
	calls := 0
	err := r.Do(func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRetryer_RetriesOnFailure(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: 3, Delay: 0, Multiplier: 1.0}
	r := NewRetryer(policy)
	r.sleep = func(time.Duration) {}

	calls := 0
	sentinel := errors.New("transient")
	err := r.Do(func() error {
		calls++
		if calls < 3 {
			return sentinel
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestRetryer_ExhaustsAttempts(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: 3, Delay: 0, Multiplier: 1.0}
	r := NewRetryer(policy)
	r.sleep = func(time.Duration) {}

	sentinel := errors.New("permanent failure")
	calls := 0
	err := r.Do(func() error {
		calls++
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestRetryer_InvalidMaxAttempts(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: 0, Delay: 0, Multiplier: 1.0}
	r := NewRetryer(policy)
	err := r.Do(func() error { return nil })
	if err == nil {
		t.Fatal("expected error for MaxAttempts=0")
	}
}

func TestRetryer_ExponentialBackoff(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: 3, Delay: 100 * time.Millisecond, Multiplier: 2.0}
	r := NewRetryer(policy)

	var delays []time.Duration
	r.sleep = func(d time.Duration) { delays = append(delays, d) }

	_ = r.Do(func() error { return errors.New("fail") })

	if len(delays) != 2 {
		t.Fatalf("expected 2 sleep calls, got %d", len(delays))
	}
	if delays[0] != 100*time.Millisecond {
		t.Errorf("expected first delay 100ms, got %v", delays[0])
	}
	if delays[1] != 200*time.Millisecond {
		t.Errorf("expected second delay 200ms, got %v", delays[1])
	}
}

func TestRetryer_Attempts(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: 5, Delay: 0, Multiplier: 1.0}
	r := NewRetryer(policy)
	if r.Attempts() != 5 {
		t.Errorf("expected 5 attempts, got %d", r.Attempts())
	}
}
