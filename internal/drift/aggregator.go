package drift

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// AggregatedReport holds drift results grouped by service.
type AggregatedReport struct {
	GeneratedAt time.Time
	Services    map[string]*ServiceDriftSummary
}

// ServiceDriftSummary holds the drift summary for a single service.
type ServiceDriftSummary struct {
	ServiceName  string
	TotalDrifts  int
	Changed      []DiffEntry
	Missing      []DiffEntry
	Extra        []DiffEntry
	LastChecked  time.Time
}

// NewAggregator creates a new Aggregator.
func NewAggregator() *Aggregator {
	return &Aggregator{
		reports: make(map[string]*ServiceDriftSummary),
	}
}

// Aggregator collects drift results across multiple services.
type Aggregator struct {
	reports map[string]*ServiceDriftSummary
}

// Add records drift entries for a given service.
func (a *Aggregator) Add(serviceName string, entries []DiffEntry) {
	if serviceName == "" {
		return
	}
	summary := &ServiceDriftSummary{
		ServiceName: serviceName,
		LastChecked: time.Now(),
	}
	for _, e := range entries {
		switch e.Status {
		case "changed":
			summary.Changed = append(summary.Changed, e)
		case "missing":
			summary.Missing = append(summary.Missing, e)
		case "extra":
			summary.Extra = append(summary.Extra, e)
		}
	}
	summary.TotalDrifts = len(summary.Changed) + len(summary.Missing) + len(summary.Extra)
	a.reports[serviceName] = summary
}

// Build returns the final AggregatedReport.
func (a *Aggregator) Build() *AggregatedReport {
	copy := make(map[string]*ServiceDriftSummary, len(a.reports))
	for k, v := range a.reports {
		copy[k] = v
	}
	return &AggregatedReport{
		GeneratedAt: time.Now(),
		Services:    copy,
	}
}

// Format returns a human-readable summary of the aggregated report.
func (r *AggregatedReport) Format() string {
	if r == nil || len(r.Services) == 0 {
		return "No drift data aggregated."
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Aggregated Drift Report — %s\n", r.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", 50))

	names := make([]string, 0, len(r.Services))
	for name := range r.Services {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		s := r.Services[name]
		fmt.Fprintf(&sb, "Service: %s | Drifts: %d (changed=%d, missing=%d, extra=%d)\n",
			s.ServiceName, s.TotalDrifts, len(s.Changed), len(s.Missing), len(s.Extra))
	}
	return sb.String()
}
