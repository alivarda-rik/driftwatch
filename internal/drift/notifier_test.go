package drift

import (
	"strings"
	"testing"
)

func TestNotifier_NoDrift(t *testing.T) {
	var buf strings.Builder
	n := NewNotifier(&buf, "")

	count := n.Notify("svc-a", nil)

	if count != 0 {
		t.Errorf("expected 0 notifications, got %d", count)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got %q", buf.String())
	}
}

func TestNotifier_WithDrift_WarnLevel(t *testing.T) {
	var buf strings.Builder
	n := NewNotifier(&buf, "")

	entries := []DiffEntry{
		{Key: "replicas", Declared: "3", Live: "2"},
	}

	count := n.Notify("svc-b", entries)

	if count != 1 {
		t.Errorf("expected 1 notification, got %d", count)
	}
	out := buf.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN level in output, got: %s", out)
	}
	if !strings.Contains(out, "svc-b") {
		t.Errorf("expected service name in output, got: %s", out)
	}
}

func TestNotifier_MissingLive_ErrorLevel(t *testing.T) {
	var buf strings.Builder
	n := NewNotifier(&buf, "")

	entries := []DiffEntry{
		{Key: "port", Declared: "8080", Live: ""},
	}

	count := n.Notify("svc-c", entries)

	if count != 1 {
		t.Errorf("expected 1 notification, got %d", count)
	}
	out := buf.String()
	if !strings.Contains(out, "ERROR") {
		t.Errorf("expected ERROR level in output, got: %s", out)
	}
}

func TestNotifier_Prefix(t *testing.T) {
	var buf strings.Builder
	n := NewNotifier(&buf, "[driftwatch]")

	entries := []DiffEntry{
		{Key: "image", Declared: "v1", Live: "v2"},
	}

	n.Notify("svc-d", entries)

	out := buf.String()
	if !strings.HasPrefix(out, "[driftwatch]") {
		t.Errorf("expected prefix in output, got: %s", out)
	}
}

func TestNotifier_NilWriter_DefaultsToStdout(t *testing.T) {
	// Should not panic when w is nil
	n := NewNotifier(nil, "")
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
	if n.out == nil {
		t.Error("expected default writer to be set")
	}
}

func TestNotifyLevel_String(t *testing.T) {
	cases := []struct {
		level NotifyLevel
		want  string
	}{
		{NotifyInfo, "INFO"},
		{NotifyWarn, "WARN"},
		{NotifyError, "ERROR"},
		{NotifyLevel(99), "UNKNOWN"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("level %d: got %q, want %q", tc.level, got, tc.want)
		}
	}
}
