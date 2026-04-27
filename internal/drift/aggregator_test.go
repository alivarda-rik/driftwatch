package drift

import (
	"strings"
	"testing"
)

func sampleEntries() []DiffEntry {
	return []DiffEntry{
		{Key: "port", Declared: "8080", Live: "9090", Status: "changed"},
		{Key: "timeout", Declared: "30s", Live: "", Status: "missing"},
		{Key: "debug", Declared: "", Live: "true", Status: "extra"},
	}
}

func TestAggregator_AddAndBuild(t *testing.T) {
	a := NewAggregator()
	a.Add("serviceA", sampleEntries())

	report := a.Build()
	if report == nil {
		t.Fatal("expected non-nil report")
	}
	if len(report.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(report.Services))
	}
	s := report.Services["serviceA"]
	if s.TotalDrifts != 3 {
		t.Errorf("expected 3 total drifts, got %d", s.TotalDrifts)
	}
	if len(s.Changed) != 1 {
		t.Errorf("expected 1 changed, got %d", len(s.Changed))
	}
	if len(s.Missing) != 1 {
		t.Errorf("expected 1 missing, got %d", len(s.Missing))
	}
	if len(s.Extra) != 1 {
		t.Errorf("expected 1 extra, got %d", len(s.Extra))
	}
}

func TestAggregator_MultipleServices(t *testing.T) {
	a := NewAggregator()
	a.Add("serviceA", sampleEntries())
	a.Add("serviceB", []DiffEntry{
		{Key: "replicas", Declared: "3", Live: "1", Status: "changed"},
	})

	report := a.Build()
	if len(report.Services) != 2 {
		t.Fatalf("expected 2 services, got %d", len(report.Services))
	}
	if report.Services["serviceB"].TotalDrifts != 1 {
		t.Errorf("expected 1 drift for serviceB")
	}
}

func TestAggregator_EmptyServiceName(t *testing.T) {
	a := NewAggregator()
	a.Add("", sampleEntries())

	report := a.Build()
	if len(report.Services) != 0 {
		t.Errorf("expected no services for empty name, got %d", len(report.Services))
	}
}

func TestAggregator_NoDriftEntries(t *testing.T) {
	a := NewAggregator()
	a.Add("serviceC", []DiffEntry{})

	report := a.Build()
	s := report.Services["serviceC"]
	if s.TotalDrifts != 0 {
		t.Errorf("expected 0 drifts, got %d", s.TotalDrifts)
	}
}

func TestAggregatedReport_Format(t *testing.T) {
	a := NewAggregator()
	a.Add("serviceA", sampleEntries())
	report := a.Build()

	out := report.Format()
	if !strings.Contains(out, "serviceA") {
		t.Errorf("expected output to contain 'serviceA', got: %s", out)
	}
	if !strings.Contains(out, "Drifts: 3") {
		t.Errorf("expected output to contain 'Drifts: 3', got: %s", out)
	}
}

func TestAggregatedReport_Format_Nil(t *testing.T) {
	var r *AggregatedReport
	out := r.Format()
	if !strings.Contains(out, "No drift data") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestAggregatedReport_Format_Empty(t *testing.T) {
	a := NewAggregator()
	report := a.Build()
	out := report.Format()
	if !strings.Contains(out, "No drift data") {
		t.Errorf("expected empty message for no services, got: %s", out)
	}
}
