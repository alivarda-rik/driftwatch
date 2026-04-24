package drift_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
)

func TestReporter_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	r := drift.NewReporter(&buf)
	report := &drift.Report{
		ServiceName: "svc",
		HasDrift:    false,
		Results:     []drift.DriftResult{},
	}
	r.Print(report)
	out := buf.String()
	if !strings.Contains(out, "No drift detected") {
		t.Errorf("expected 'No drift detected' in output, got: %s", out)
	}
}

func TestReporter_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	r := drift.NewReporter(&buf)
	report := &drift.Report{
		ServiceName: "svc",
		HasDrift:    true,
		Results: []drift.DriftResult{
			{Field: "version", Expected: "1.0.0", Actual: "2.0.0", Drifted: true},
			{Field: "port", Expected: 8080, Actual: 9090, Drifted: true},
			{Field: "image", Expected: "api:latest", Actual: "api:latest", Drifted: false},
		},
	}
	r.Print(report)
	out := buf.String()
	if !strings.Contains(out, "Drift detected") {
		t.Errorf("expected 'Drift detected' in output, got: %s", out)
	}
	if !strings.Contains(out, "version") {
		t.Errorf("expected 'version' field in output")
	}
	if !strings.Contains(out, "port") {
		t.Errorf("expected 'port' field in output")
	}
	if strings.Contains(out, "image") {
		t.Errorf("did not expect non-drifted 'image' field in output")
	}
	if !strings.Contains(out, "Total drifted fields: 2") {
		t.Errorf("expected drift count summary in output")
	}
}

func TestReporter_NilReport(t *testing.T) {
	var buf bytes.Buffer
	r := drift.NewReporter(&buf)
	r.Print(nil)
	if !strings.Contains(buf.String(), "no report available") {
		t.Errorf("expected fallback message for nil report")
	}
}
