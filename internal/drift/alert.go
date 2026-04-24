package drift

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// AlertLevel represents the severity of a drift alert.
type AlertLevel string

const (
	AlertLevelInfo    AlertLevel = "INFO"
	AlertLevelWarning AlertLevel = "WARNING"
	AlertLevelCritical AlertLevel = "CRITICAL"
)

// Alert represents a drift alert for a service.
type Alert struct {
	Service   string
	Level     AlertLevel
	Drifts    []DiffEntry
	Timestamp time.Time
}

// String returns a human-readable representation of the alert.
func (a *Alert) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] %s drift alert for service %q at %s\n",
		a.Level, len(a.Drifts), a.Service, a.Timestamp.Format(time.RFC3339)))
	for _, d := range a.Drifts {
		sb.WriteString(fmt.Sprintf("  - %s\n", d.String()))
	}
	return sb.String()
}

// Alerter evaluates drift results and emits alerts.
type Alerter struct {
	out            io.Writer
	criticalThreshold int
}

// NewAlerter creates an Alerter that writes to the given writer.
// criticalThreshold defines how many drift entries trigger a CRITICAL level.
func NewAlerter(out io.Writer, criticalThreshold int) *Alerter {
	if out == nil {
		out = os.Stdout
	}
	if criticalThreshold <= 0 {
		criticalThreshold = 5
	}
	return &Alerter{out: out, criticalThreshold: criticalThreshold}
}

// Evaluate inspects the diff entries and returns an Alert, or nil if no drift.
func (a *Alerter) Evaluate(service string, entries []DiffEntry) *Alert {
	if len(entries) == 0 {
		return nil
	}
	level := AlertLevelWarning
	if len(entries) >= a.criticalThreshold {
		level = AlertLevelCritical
	}
	return &Alert{
		Service:   service,
		Level:     level,
		Drifts:    entries,
		Timestamp: time.Now().UTC(),
	}
}

// Emit writes the alert to the configured writer.
func (a *Alerter) Emit(alert *Alert) {
	if alert == nil {
		return
	}
	fmt.Fprint(a.out, alert.String())
}
