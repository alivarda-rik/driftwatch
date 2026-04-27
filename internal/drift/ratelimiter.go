package drift

import (
	"sync"
	"time"
)

// RateLimiter controls how frequently drift checks can be triggered per service,
// using a token bucket approach with a fixed refill interval.
type RateLimiter struct {
	mu       sync.Mutex
	tokens   map[string]int
	lastFill map[string]time.Time
	max      int
	interval time.Duration
	now      func() time.Time
}

// NewRateLimiter creates a RateLimiter allowing up to max checks per interval per service.
func NewRateLimiter(max int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:   make(map[string]int),
		lastFill: make(map[string]time.Time),
		max:      max,
		interval: interval,
		now:      time.Now,
	}
}

// Allow returns true if a drift check for the given service is permitted.
// It refills tokens if the interval has elapsed since the last fill.
func (r *RateLimiter) Allow(service string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.now()

	last, seen := r.lastFill[service]
	if !seen || now.Sub(last) >= r.interval {
		r.tokens[service] = r.max
		r.lastFill[service] = now
	}

	if r.tokens[service] <= 0 {
		return false
	}

	r.tokens[service]--
	return true
}

// Remaining returns the number of tokens left for the given service.
func (r *RateLimiter) Remaining(service string) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.now()
	last, seen := r.lastFill[service]
	if !seen || now.Sub(last) >= r.interval {
		return r.max
	}
	return r.tokens[service]
}

// Reset clears the rate limit state for a given service.
func (r *RateLimiter) Reset(service string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tokens, service)
	delete(r.lastFill, service)
}
