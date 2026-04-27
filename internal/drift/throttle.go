package drift

import (
	"sync"
	"time"
)

// ThrottleConfig controls how often alerts can be emitted per service.
type ThrottleConfig struct {
	// MinInterval is the minimum duration between alerts for the same service.
	MinInterval time.Duration
}

// Throttle prevents duplicate alerts from firing too frequently for the same service.
type Throttle struct {
	mu       sync.Mutex
	lastSeen map[string]time.Time
	cfg      ThrottleConfig
}

// NewThrottle creates a new Throttle with the given configuration.
func NewThrottle(cfg ThrottleConfig) *Throttle {
	if cfg.MinInterval <= 0 {
		cfg.MinInterval = 30 * time.Second
	}
	return &Throttle{
		lastSeen: make(map[string]time.Time),
		cfg:      cfg,
	}
}

// Allow returns true if an alert for the given service should be allowed
// through based on the configured minimum interval.
func (t *Throttle) Allow(service string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	last, seen := t.lastSeen[service]
	if seen && now.Sub(last) < t.cfg.MinInterval {
		return false
	}
	t.lastSeen[service] = now
	return true
}

// Reset clears the throttle state for a specific service.
func (t *Throttle) Reset(service string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.lastSeen, service)
}

// ResetAll clears all throttle state.
func (t *Throttle) ResetAll() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastSeen = make(map[string]time.Time)
}

// LastSeen returns the last time an alert was allowed for the given service,
// and whether it has been seen at all.
func (t *Throttle) LastSeen(service string) (time.Time, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	last, ok := t.lastSeen[service]
	return last, ok
}
