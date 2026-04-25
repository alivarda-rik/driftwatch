package drift

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestExporter_JSON_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	ex := NewExporter(FormatJSON, &buf)
	if err := ex.Export("svc-a", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var rec ExportRecord
	if err := json.Unmarshal(buf.Bytes(), &rec); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if rec.Service != "svc-a" {
		t.Errorf("expected service svc-a, got %s", rec.Service)
	}
	if rec.DriftCount != 0 {
		t.Errorf("expected drift_count 0, got %d", rec.DriftCount)
	}
}

func TestExporter_JSON_WithDrift(t *testing.T) {
	entries := []DiffEntry{
		{Key: "port", Declared: "8080", Live: "9090", Kind: KindChanged},
	}
	var buf bytes.Buffer
	ex := NewExporter(FormatJSON, &buf)
	if err := ex.Export("svc-b", entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var rec ExportRecord
	if err := json.Unmarshal(buf.Bytes(), &rec); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if rec.DriftCount != 1 {
		t.Errorf("expected drift_count 1, got %d", rec.DriftCount)
	}
	if len(rec.Entries) != 1 || rec.Entries[0].Key != "port" {
		t.Errorf("unexpected entries: %+v", rec.Entries)
	}
}

func TestExporter_Text_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	ex := NewExporter(FormatText, &buf)
	if err := ex.Export("svc-c", []DiffEntry{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "svc-c") {
		t.Errorf("expected service name in output")
	}
	if !strings.Contains(out, "no drift detected") {
		t.Errorf("expected no-drift message in output")
	}
}

func TestExporter_Text_WithDrift(t *testing.T) {
	entries := []DiffEntry{
		{Key: "replicas", Declared: "3", Live: "1", Kind: KindChanged},
	}
	var buf bytes.Buffer
	ex := NewExporter(FormatText, &buf)
	if err := ex.Export("svc-d", entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "replicas") {
		t.Errorf("expected key 'replicas' in text output")
	}
}

func TestExporter_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	ex := NewExporter(ExportFormat("xml"), &buf)
	if err := ex.Export("svc-e", nil); err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestExporter_Nil(t *testing.T) {
	var ex *Exporter
	if err := ex.Export("svc-f", nil); err == nil {
		t.Error("expected error for nil exporter")
	}
}
