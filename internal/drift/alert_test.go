package drift

import (
	"strings"
	"testing"
	"time"
)

func sampleDiffEntries() []DiffEntry {
	return []DiffEntry{
		{Key: "replicas", Declared: "3", Live: "2"},
		{Key: "image", Declared: "app:v2", Live: "app:v1"},
	}
}

func TestAlerter_NoDrift(t *testing.T) {
	a := NewAlerter(nil, 5)
	alert := a.Evaluate("svc-a", []DiffEntry{})
	if alert != nil {
		t.Errorf("expected nil alert for no drift, got %v", alert)
	}
}

func TestAlerter_WarningLevel(t *testing.T) {
	a := NewAlerter(nil, 5)
	entries := sampleDiffEntries()
	alert := a.Evaluate("svc-b", entries)
	if alert == nil {
		t.Fatal("expected alert, got nil")
	}
	if alert.Level != AlertLevelWarning {
		t.Errorf("expected WARNING, got %s", alert.Level)
	}
	if alert.Service != "svc-b" {
		t.Errorf("expected service svc-b, got %s", alert.Service)
	}
	if len(alert.Drifts) != 2 {
		t.Errorf("expected 2 drifts, got %d", len(alert.Drifts))
	}
}

func TestAlerter_CriticalLevel(t *testing.T) {
	a := NewAlerter(nil, 2)
	entries := sampleDiffEntries()
	alert := a.Evaluate("svc-c", entries)
	if alert == nil {
		t.Fatal("expected alert, got nil")
	}
	if alert.Level != AlertLevelCritical {
		t.Errorf("expected CRITICAL, got %s", alert.Level)
	}
}

func TestAlerter_Emit(t *testing.T) {
	var buf strings.Builder
	a := NewAlerter(&buf, 5)
	alert := &Alert{
		Service:   "svc-d",
		Level:     AlertLevelWarning,
		Drifts:    sampleDiffEntries(),
		Timestamp: time.Now().UTC(),
	}
	a.Emit(alert)
	out := buf.String()
	if !strings.Contains(out, "svc-d") {
		t.Errorf("expected service name in output, got: %s", out)
	}
	if !strings.Contains(out, "WARNING") {
		t.Errorf("expected WARNING in output, got: %s", out)
	}
}

func TestAlerter_EmitNil(t *testing.T) {
	var buf strings.Builder
	a := NewAlerter(&buf, 5)
	a.Emit(nil)
	if buf.Len() != 0 {
		t.Errorf("expected no output for nil alert, got: %s", buf.String())
	}
}

func TestAlert_String(t *testing.T) {
	alert := &Alert{
		Service:   "svc-e",
		Level:     AlertLevelCritical,
		Drifts:    sampleDiffEntries(),
		Timestamp: time.Now().UTC(),
	}
	s := alert.String()
	if !strings.Contains(s, "CRITICAL") {
		t.Errorf("expected CRITICAL in string, got: %s", s)
	}
	if !strings.Contains(s, "svc-e") {
		t.Errorf("expected svc-e in string, got: %s", s)
	}
}
