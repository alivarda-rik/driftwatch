package drift

import (
	"testing"
	"time"
)

func TestThrottle_AllowFirstCall(t *testing.T) {
	th := NewThrottle(ThrottleConfig{MinInterval: 1 * time.Second})
	if !th.Allow("svc-a") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestThrottle_BlocksWithinInterval(t *testing.T) {
	th := NewThrottle(ThrottleConfig{MinInterval: 10 * time.Second})
	th.Allow("svc-a")
	if th.Allow("svc-a") {
		t.Fatal("expected second call within interval to be blocked")
	}
}

func TestThrottle_AllowsAfterInterval(t *testing.T) {
	th := NewThrottle(ThrottleConfig{MinInterval: 50 * time.Millisecond})
	th.Allow("svc-b")
	time.Sleep(60 * time.Millisecond)
	if !th.Allow("svc-b") {
		t.Fatal("expected call after interval to be allowed")
	}
}

func TestThrottle_IndependentServices(t *testing.T) {
	th := NewThrottle(ThrottleConfig{MinInterval: 10 * time.Second})
	th.Allow("svc-a")
	if !th.Allow("svc-b") {
		t.Fatal("expected different service to be allowed independently")
	}
}

func TestThrottle_Reset(t *testing.T) {
	th := NewThrottle(ThrottleConfig{MinInterval: 10 * time.Second})
	th.Allow("svc-a")
	th.Reset("svc-a")
	if !th.Allow("svc-a") {
		t.Fatal("expected call to be allowed after reset")
	}
}

func TestThrottle_ResetAll(t *testing.T) {
	th := NewThrottle(ThrottleConfig{MinInterval: 10 * time.Second})
	th.Allow("svc-a")
	th.Allow("svc-b")
	th.ResetAll()
	if !th.Allow("svc-a") || !th.Allow("svc-b") {
		t.Fatal("expected all services to be allowed after ResetAll")
	}
}

func TestThrottle_DefaultInterval(t *testing.T) {
	th := NewThrottle(ThrottleConfig{})
	if th.cfg.MinInterval != 30*time.Second {
		t.Fatalf("expected default interval 30s, got %v", th.cfg.MinInterval)
	}
}

func TestThrottle_LastSeen(t *testing.T) {
	th := NewThrottle(ThrottleConfig{MinInterval: 1 * time.Second})
	_, seen := th.LastSeen("svc-x")
	if seen {
		t.Fatal("expected no last seen before any call")
	}
	before := time.Now()
	th.Allow("svc-x")
	last, seen := th.LastSeen("svc-x")
	if !seen {
		t.Fatal("expected last seen after Allow")
	}
	if last.Before(before) {
		t.Fatal("expected last seen timestamp to be recent")
	}
}
