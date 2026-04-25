package drift

import (
	"strings"
	"testing"
)

func TestNewSummaryReport_Empty(t *testing.T) {
	report := NewSummaryReport(map[string][]DiffEntry{})
	if report.TotalChecked != 0 {
		t.Errorf("expected 0 total, got %d", report.TotalChecked)
	}
	if report.DriftedCount != 0 || report.CleanCount != 0 {
		t.Error("expected zero drifted and clean counts")
	}
}

func TestNewSummaryReport_NoDrift(t *testing.T) {
	results := map[string][]DiffEntry{
		"svc-a": {},
		"svc-b": {},
	}
	report := NewSummaryReport(results)
	if report.TotalChecked != 2 {
		t.Errorf("expected 2 total, got %d", report.TotalChecked)
	}
	if report.CleanCount != 2 {
		t.Errorf("expected 2 clean, got %d", report.CleanCount)
	}
	if report.DriftedCount != 0 {
		t.Errorf("expected 0 drifted, got %d", report.DriftedCount)
	}
}

func TestNewSummaryReport_WithDrift(t *testing.T) {
	results := map[string][]DiffEntry{
		"svc-a": {
			{Key: "port", Kind: DiffChanged, Declared: "8080", Live: "9090"},
			{Key: "timeout", Kind: DiffMissing, Declared: "30s", Live: ""},
		},
		"svc-b": {},
	}
	report := NewSummaryReport(results)
	if report.DriftedCount != 1 {
		t.Errorf("expected 1 drifted service, got %d", report.DriftedCount)
	}
	if report.CleanCount != 1 {
		t.Errorf("expected 1 clean service, got %d", report.CleanCount)
	}
	var drifted *ServiceSummary
	for i := range report.Services {
		if report.Services[i].ServiceName == "svc-a" {
			drifted = &report.Services[i]
		}
	}
	if drifted == nil {
		t.Fatal("svc-a not found in summary")
	}
	if drifted.ChangedCount != 1 {
		t.Errorf("expected 1 changed, got %d", drifted.ChangedCount)
	}
	if drifted.MissingCount != 1 {
		t.Errorf("expected 1 missing, got %d", drifted.MissingCount)
	}
}

func TestSummaryReport_Format_ContainsServiceName(t *testing.T) {
	results := map[string][]DiffEntry{
		"my-service": {
			{Key: "replicas", Kind: DiffExtra, Declared: "", Live: "5"},
		},
	}
	report := NewSummaryReport(results)
	output := report.Format()
	if !strings.Contains(output, "my-service") {
		t.Error("expected service name in formatted output")
	}
	if !strings.Contains(output, "DRIFT") {
		t.Error("expected DRIFT label in formatted output")
	}
}

func TestSummaryReport_Format_Nil(t *testing.T) {
	var r *SummaryReport
	out := r.Format()
	if out != "no summary available" {
		t.Errorf("unexpected output for nil report: %q", out)
	}
}
