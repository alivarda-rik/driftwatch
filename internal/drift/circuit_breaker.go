package drift

import (
	"fmt"
	"sync"
	"time"
)

// CircuitState represents the state of a circuit breaker.
type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

func (s CircuitState) String() string {
	switch s {
	case CircuitClosed:
		return "closed"
	case CircuitOpen:
		return "open"
	case CircuitHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// CircuitBreaker prevents repeated drift checks for consistently failing services.
type CircuitBreaker struct {
	mu           sync.Mutex
	states        map[string]*circuitEntry
	failureLimit  int
	resetTimeout  time.Duration
}

type circuitEntry struct {
	state     CircuitState
	failures  int
	openedAt  time.Time
}

// NewCircuitBreaker creates a CircuitBreaker that opens after failureLimit
// consecutive failures and attempts reset after resetTimeout.
func NewCircuitBreaker(failureLimit int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		states:       make(map[string]*circuitEntry),
		failureLimit: failureLimit,
		resetTimeout: resetTimeout,
	}
}

// Allow returns true if the service is allowed to proceed with a drift check.
func (cb *CircuitBreaker) Allow(service string) bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	e := cb.entry(service)
	switch e.state {
	case CircuitClosed:
		return true
	case CircuitOpen:
		if time.Since(e.openedAt) >= cb.resetTimeout {
			e.state = CircuitHalfOpen
			return true
		}
		return false
	case CircuitHalfOpen:
		return true
	}
	return false
}

// RecordSuccess resets failure count and closes the circuit for a service.
func (cb *CircuitBreaker) RecordSuccess(service string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	e := cb.entry(service)
	e.failures = 0
	e.state = CircuitClosed
}

// RecordFailure increments the failure count and may open the circuit.
func (cb *CircuitBreaker) RecordFailure(service string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	e := cb.entry(service)
	e.failures++
	if e.failures >= cb.failureLimit {
		e.state = CircuitOpen
		e.openedAt = time.Now()
	}
}

// State returns the current CircuitState for a service.
func (cb *CircuitBreaker) State(service string) CircuitState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.entry(service).state
}

// Reset forcibly closes the circuit for a service.
func (cb *CircuitBreaker) Reset(service string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.states[service] = &circuitEntry{state: CircuitClosed}
}

func (cb *CircuitBreaker) entry(service string) *circuitEntry {
	if _, ok := cb.states[service]; !ok {
		cb.states[service] = &circuitEntry{state: CircuitClosed}
	}
	return cb.states[service]
}

// ErrCircuitOpen is returned when a check is blocked by an open circuit.
var ErrCircuitOpen = fmt.Errorf("circuit open: service check blocked")
