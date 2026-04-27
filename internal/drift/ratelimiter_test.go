package drift

import (
	"testing"
	"time"
)

func TestRateLimiter_AllowFirstCall(t *testing.T) {
	rl := NewRateLimiter(3, time.Second)
	if !rl.Allow("svc-a") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestRateLimiter_ExhaustsTokens(t *testing.T) {
	rl := NewRateLimiter(2, time.Second)
	if !rl.Allow("svc-a") {
		t.Fatal("expected call 1 to be allowed")
	}
	if !rl.Allow("svc-a") {
		t.Fatal("expected call 2 to be allowed")
	}
	if rl.Allow("svc-a") {
		t.Fatal("expected call 3 to be denied")
	}
}

func TestRateLimiter_RefillsAfterInterval(t *testing.T) {
	now := time.Now()
	rl := NewRateLimiter(1, time.Second)
	rl.now = func() time.Time { return now }

	if !rl.Allow("svc-a") {
		t.Fatal("expected first call to be allowed")
	}
	if rl.Allow("svc-a") {
		t.Fatal("expected second call to be denied")
	}

	// Advance time past the interval
	rl.now = func() time.Time { return now.Add(2 * time.Second) }
	if !rl.Allow("svc-a") {
		t.Fatal("expected call after interval to be allowed")
	}
}

func TestRateLimiter_IndependentServices(t *testing.T) {
	rl := NewRateLimiter(1, time.Second)
	if !rl.Allow("svc-a") {
		t.Fatal("expected svc-a to be allowed")
	}
	if !rl.Allow("svc-b") {
		t.Fatal("expected svc-b to be allowed independently")
	}
	if rl.Allow("svc-a") {
		t.Fatal("expected svc-a second call to be denied")
	}
}

func TestRateLimiter_Remaining(t *testing.T) {
	rl := NewRateLimiter(3, time.Second)
	if got := rl.Remaining("svc-a"); got != 3 {
		t.Fatalf("expected 3 remaining, got %d", got)
	}
	rl.Allow("svc-a")
	if got := rl.Remaining("svc-a"); got != 2 {
		t.Fatalf("expected 2 remaining after one call, got %d", got)
	}
}

func TestRateLimiter_Reset(t *testing.T) {
	rl := NewRateLimiter(1, time.Second)
	rl.Allow("svc-a")
	if rl.Allow("svc-a") {
		t.Fatal("expected call to be denied before reset")
	}
	rl.Reset("svc-a")
	if !rl.Allow("svc-a") {
		t.Fatal("expected call to be allowed after reset")
	}
}
