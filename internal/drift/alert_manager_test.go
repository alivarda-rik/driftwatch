package drift

import (
	"strings"
	"testing"

	"github.com/user/driftwatch/internal/config"
)

func alertManagerConfig() *config.ServiceConfig {
	return &config.ServiceConfig{
		Name: "alert-svc",
		Params: map[string]string{
			"replicas": "3",
			"image":    "app:v2",
		},
	}
}

func TestAlertManager_NoDrift(t *testing.T) {
	var buf strings.Builder
	cfg := alertManagerConfig()
	detector := NewDetector(cfg)
	alerter := NewAlerter(&buf, 5)
	mgr := NewAlertManager(alerter, detector, &buf)

	live := map[string]string{"replicas": "3", "image": "app:v2"}
	alert, err := mgr.CheckAndAlert("alert-svc", live)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alert != nil {
		t.Errorf("expected no alert, got %v", alert)
	}
}

func TestAlertManager_WithDrift(t *testing.T) {
	var buf strings.Builder
	cfg := alertManagerConfig()
	detector := NewDetector(cfg)
	alerter := NewAlerter(&buf, 5)
	mgr := NewAlertManager(alerter, detector, &buf)

	live := map[string]string{"replicas": "1", "image": "app:v1"}
	alert, err := mgr.CheckAndAlert("alert-svc", live)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alert == nil {
		t.Fatal("expected alert, got nil")
	}
	if alert.Level != AlertLevelWarning {
		t.Errorf("expected WARNING, got %s", alert.Level)
	}
}

func TestAlertManager_NilDetector(t *testing.T) {
	var buf strings.Builder
	alerter := NewAlerter(&buf, 5)
	mgr := NewAlertManager(alerter, nil, &buf)

	_, err := mgr.CheckAndAlert("svc", map[string]string{})
	if err == nil {
		t.Error("expected error for nil detector")
	}
}

func TestSummaryLine_NoAlert(t *testing.T) {
	s := SummaryLine(nil)
	if s != "no drift detected" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummaryLine_WithAlert(t *testing.T) {
	alert := &Alert{
		Service: "svc-x",
		Level:   AlertLevelCritical,
		Drifts:  sampleDiffEntries(),
	}
	s := SummaryLine(alert)
	if !strings.Contains(s, "svc-x") {
		t.Errorf("expected svc-x in summary: %s", s)
	}
	if !strings.Contains(s, "CRITICAL") {
		t.Errorf("expected CRITICAL in summary: %s", s)
	}
}
