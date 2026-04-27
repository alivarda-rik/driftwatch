package drift_test

import (
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
)

func TestAggregator_Integration_MultiServiceReport(t *testing.T) {
	agg := drift.NewAggregator()

	services := map[string][]drift.DiffEntry{
		"auth-service": {
			{Key: "port", Declared: "8080", Live: "9090", Status: "changed"},
			{Key: "log_level", Declared: "info", Live: "", Status: "missing"},
		},
		"payment-service": {
			{Key: "timeout", Declared: "30s", Live: "60s", Status: "changed"},
		},
		"notification-service": {},
	}

	for svc, entries := range services {
		agg.Add(svc, entries)
	}

	report := agg.Build()

	if len(report.Services) != 3 {
		t.Fatalf("expected 3 services, got %d", len(report.Services))
	}

	authSummary, ok := report.Services["auth-service"]
	if !ok {
		t.Fatal("expected auth-service in report")
	}
	if authSummary.TotalDrifts != 2 {
		t.Errorf("expected 2 drifts for auth-service, got %d", authSummary.TotalDrifts)
	}

	paymentSummary := report.Services["payment-service"]
	if paymentSummary.TotalDrifts != 1 {
		t.Errorf("expected 1 drift for payment-service, got %d", paymentSummary.TotalDrifts)
	}

	notifSummary := report.Services["notification-service"]
	if notifSummary.TotalDrifts != 0 {
		t.Errorf("expected 0 drifts for notification-service, got %d", notifSummary.TotalDrifts)
	}

	formatted := report.Format()
	for _, name := range []string{"auth-service", "payment-service", "notification-service"} {
		if !strings.Contains(formatted, name) {
			t.Errorf("expected formatted output to contain %q", name)
		}
	}
}
