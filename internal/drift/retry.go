package drift

import (
	"errors"
	"time"
)

// RetryPolicy defines the configuration for retry behaviour.
type RetryPolicy struct {
	MaxAttempts int
	Delay       time.Duration
	Multiplier  float64
}

// DefaultRetryPolicy returns a sensible default retry policy.
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxAttempts: 3,
		Delay:       200 * time.Millisecond,
		Multiplier:  2.0,
	}
}

// Retryer executes an operation with retry logic using exponential backoff.
type Retryer struct {
	policy RetryPolicy
	sleep  func(time.Duration)
}

// NewRetryer creates a Retryer with the given policy.
func NewRetryer(policy RetryPolicy) *Retryer {
	return &Retryer{
		policy: policy,
		sleep:  time.Sleep,
	}
}

// Do executes fn up to MaxAttempts times, backing off between failures.
// Returns the last error if all attempts fail.
func (r *Retryer) Do(fn func() error) error {
	if r.policy.MaxAttempts <= 0 {
		return errors.New("retry: MaxAttempts must be greater than 0")
	}

	delay := r.policy.Delay
	var lastErr error

	for attempt := 1; attempt <= r.policy.MaxAttempts; attempt++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if attempt < r.policy.MaxAttempts {
			r.sleep(delay)
			delay = time.Duration(float64(delay) * r.policy.Multiplier)
		}
	}

	return lastErr
}

// Attempts returns the configured max attempts.
func (r *Retryer) Attempts() int {
	return r.policy.MaxAttempts
}
