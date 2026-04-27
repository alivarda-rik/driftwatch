package drift

import (
	"testing"
	"time"
)

func TestCircuitBreaker_InitiallyAllows(t *testing.T) {
	cb := NewCircuitBreaker(3, 100*time.Millisecond)
	if !cb.Allow("svc") {
		t.Error("expected circuit to allow on first call")
	}
	if cb.State("svc") != CircuitClosed {
		t.Errorf("expected closed, got %s", cb.State("svc"))
	}
}

func TestCircuitBreaker_OpensAfterFailureLimit(t *testing.T) {
	cb := NewCircuitBreaker(3, time.Second)
	for i := 0; i < 3; i++ {
		cb.RecordFailure("svc")
	}
	if cb.State("svc") != CircuitOpen {
		t.Errorf("expected open, got %s", cb.State("svc"))
	}
	if cb.Allow("svc") {
		t.Error("expected circuit to block when open")
	}
}

func TestCircuitBreaker_HalfOpenAfterTimeout(t *testing.T) {
	cb := NewCircuitBreaker(2, 50*time.Millisecond)
	cb.RecordFailure("svc")
	cb.RecordFailure("svc")

	time.Sleep(60 * time.Millisecond)

	if !cb.Allow("svc") {
		t.Error("expected circuit to allow after reset timeout")
	}
	if cb.State("svc") != CircuitHalfOpen {
		t.Errorf("expected half-open, got %s", cb.State("svc"))
	}
}

func TestCircuitBreaker_SuccessCloses(t *testing.T) {
	cb := NewCircuitBreaker(2, 50*time.Millisecond)
	cb.RecordFailure("svc")
	cb.RecordFailure("svc")
	time.Sleep(60 * time.Millisecond)
	_ = cb.Allow("svc") // transition to half-open
	cb.RecordSuccess("svc")

	if cb.State("svc") != CircuitClosed {
		t.Errorf("expected closed after success, got %s", cb.State("svc"))
	}
}

func TestCircuitBreaker_IndependentServices(t *testing.T) {
	cb := NewCircuitBreaker(2, time.Second)
	cb.RecordFailure("alpha")
	cb.RecordFailure("alpha")

	if !cb.Allow("beta") {
		t.Error("beta should be unaffected by alpha failures")
	}
	if cb.Allow("alpha") {
		t.Error("alpha circuit should be open")
	}
}

func TestCircuitBreaker_Reset(t *testing.T) {
	cb := NewCircuitBreaker(2, time.Hour)
	cb.RecordFailure("svc")
	cb.RecordFailure("svc")

	cb.Reset("svc")

	if cb.State("svc") != CircuitClosed {
		t.Errorf("expected closed after reset, got %s", cb.State("svc"))
	}
	if !cb.Allow("svc") {
		t.Error("expected allow after reset")
	}
}

func TestCircuitState_String(t *testing.T) {
	cases := []struct {
		state CircuitState
		want  string
	}{
		{CircuitClosed, "closed"},
		{CircuitOpen, "open"},
		{CircuitHalfOpen, "half-open"},
		{CircuitState(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.state.String(); got != tc.want {
			t.Errorf("State(%d).String() = %q, want %q", tc.state, got, tc.want)
		}
	}
}
