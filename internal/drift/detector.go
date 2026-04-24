package drift

import (
	"fmt"

	"github.com/driftwatch/internal/config"
)

// DriftResult holds the result of a drift check for a single field.
type DriftResult struct {
	Field    string
	Expected interface{}
	Actual   interface{}
	Drifted  bool
}

// Report aggregates all drift results for a service.
type Report struct {
	ServiceName string
	Results     []DriftResult
	HasDrift    bool
}

// Detector compares a declared config against a live state map.
type Detector struct{}

// NewDetector creates a new Detector instance.
func NewDetector() *Detector {
	return &Detector{}
}

// Detect compares the declared config with the live state and returns a Report.
func (d *Detector) Detect(declared *config.ServiceConfig, live map[string]interface{}) (*Report, error) {
	if declared == nil {
		return nil, fmt.Errorf("declared config must not be nil")
	}

	report := &Report{
		ServiceName: declared.Name,
	}

	declaredMap := flattenDeclared(declared)

	for field, expectedVal := range declaredMap {
		actualVal, exists := live[field]
		result := DriftResult{
			Field:    field,
			Expected: expectedVal,
		}
		if !exists {
			result.Actual = nil
			result.Drifted = true
		} else {
			result.Actual = actualVal
			result.Drifted = fmt.Sprintf("%v", expectedVal) != fmt.Sprintf("%v", actualVal)
		}
		if result.Drifted {
			report.HasDrift = true
		}
		report.Results = append(report.Results, result)
	}

	return report, nil
}

// flattenDeclared converts a ServiceConfig into a flat key-value map.
func flattenDeclared(cfg *config.ServiceConfig) map[string]interface{} {
	m := map[string]interface{}{
		"name":    cfg.Name,
		"version": cfg.Version,
		"image":   cfg.Image,
		"port":    cfg.Port,
		"replicas": cfg.Replicas,
	}
	for k, v := range cfg.Env {
		m["env."+k] = v
	}
	return m
}
