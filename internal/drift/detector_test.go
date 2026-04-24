package drift_test

import (
	"testing"

	"github.com/driftwatch/internal/config"
	"github.com/driftwatch/internal/drift"
)

func baseConfig() *config.ServiceConfig {
	return &config.ServiceConfig{
		Name:     "api",
		Version:  "1.0.0",
		Image:    "api:latest",
		Port:     8080,
		Replicas: 2,
		Env:      map[string]string{"LOG_LEVEL": "info"},
	}
}

func TestDetect_NoDrift(t *testing.T) {
	d := drift.NewDetector()
	live := map[string]interface{}{
		"name":          "api",
		"version":       "1.0.0",
		"image":         "api:latest",
		"port":          8080,
		"replicas":      2,
		"env.LOG_LEVEL": "info",
	}
	report, err := d.Detect(baseConfig(), live)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.HasDrift {
		t.Errorf("expected no drift, but drift was reported")
	}
}

func TestDetect_WithDrift(t *testing.T) {
	d := drift.NewDetector()
	live := map[string]interface{}{
		"name":          "api",
		"version":       "2.0.0", // drifted
		"image":         "api:latest",
		"port":          9090, // drifted
		"replicas":      2,
		"env.LOG_LEVEL": "info",
	}
	report, err := d.Detect(baseConfig(), live)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !report.HasDrift {
		t.Errorf("expected drift, but none was reported")
	}
	driftedFields := map[string]bool{}
	for _, r := range report.Results {
		if r.Drifted {
			driftedFields[r.Field] = true
		}
	}
	if !driftedFields["version"] {
		t.Errorf("expected 'version' to be drifted")
	}
	if !driftedFields["port"] {
		t.Errorf("expected 'port' to be drifted")
	}
}

func TestDetect_MissingField(t *testing.T) {
	d := drift.NewDetector()
	live := map[string]interface{}{} // empty live state
	report, err := d.Detect(baseConfig(), live)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !report.HasDrift {
		t.Errorf("expected drift due to missing fields")
	}
}

func TestDetect_NilConfig(t *testing.T) {
	d := drift.NewDetector()
	_, err := d.Detect(nil, map[string]interface{}{})
	if err == nil {
		t.Errorf("expected error for nil config, got none")
	}
}
