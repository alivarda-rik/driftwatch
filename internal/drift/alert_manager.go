package drift

import (
	"fmt"
	"io"
	"os"
)

// AlertManager orchestrates drift evaluation and alerting across services.
type AlertManager struct {
	alerter  *Alerter
	detector *Detector
	out      io.Writer
}

// NewAlertManager creates an AlertManager using the given alerter and detector.
// If out is nil, os.Stdout is used.
func NewAlertManager(alerter *Alerter, detector *Detector, out io.Writer) *AlertManager {
	if out == nil {
		out = os.Stdout
	}
	return &AlertManager{
		alerter:  alerter,
		detector: detector,
		out:      out,
	}
}

// CheckAndAlert detects drift for a service and emits an alert if drift is found.
// Returns the alert (or nil if no drift) and any detection error.
func (m *AlertManager) CheckAndAlert(service string, live map[string]string) (*Alert, error) {
	if m.detector == nil {
		return nil, fmt.Errorf("alert manager: detector is nil")
	}
	if m.alerter == nil {
		return nil, fmt.Errorf("alert manager: alerter is nil")
	}

	report, err := m.detector.Detect(live)
	if err != nil {
		return nil, fmt.Errorf("alert manager: detect error: %w", err)
	}

	var entries []DiffEntry
	if report != nil {
		entries = report.Diffs
	}

	alert := m.alerter.Evaluate(service, entries)
	m.alerter.Emit(alert)
	return alert, nil
}

// CheckAndAlertMulti runs CheckAndAlert for each service in the provided map
// and returns a slice of alerts for services where drift was detected.
// Errors encountered during detection are collected and returned as a combined error.
func (m *AlertManager) CheckAndAlertMulti(services map[string]map[string]string) ([]*Alert, error) {
	var alerts []*Alert
	var errs []error

	for service, live := range services {
		alert, err := m.CheckAndAlert(service, live)
		if err != nil {
			errs = append(errs, fmt.Errorf("service %q: %w", service, err))
			continue
		}
		if alert != nil {
			alerts = append(alerts, alert)
		}
	}

	if len(errs) > 0 {
		return alerts, fmt.Errorf("alert manager: %d error(s) during multi-check: %v", len(errs), errs)
	}
	return alerts, nil
}

// SummaryLine returns a one-line summary suitable for logging.
func SummaryLine(alert *Alert) string {
	if alert == nil {
		return "no drift detected"
	}
	return fmt.Sprintf("service=%s level=%s drifts=%d",
		alert.Service, alert.Level, len(alert.Drifts))
}
