package drift

import (
	"fmt"
	"strings"
	"time"
)

// SummaryReport holds a high-level overview of drift detection results
// for one or more services at a point in time.
type SummaryReport struct {
	GeneratedAt  time.Time
	TotalChecked int
	DriftedCount int
	CleanCount   int
	Services     []ServiceSummary
}

// ServiceSummary holds drift summary data for a single service.
type ServiceSummary struct {
	ServiceName  string
	HasDrift     bool
	ChangedCount int
	MissingCount int
	ExtraCount   int
}

// NewSummaryReport builds a SummaryReport from a map of service name to diff entries.
func NewSummaryReport(results map[string][]DiffEntry) *SummaryReport {
	report := &SummaryReport{
		GeneratedAt:  time.Now().UTC(),
		TotalChecked: len(results),
	}

	for svc, entries := range results {
		ss := ServiceSummary{ServiceName: svc}
		for _, e := range entries {
			switch e.Kind {
			case DiffChanged:
				ss.ChangedCount++
			case DiffMissing:
				ss.MissingCount++
			case DiffExtra:
				ss.ExtraCount++
			}
		}
		ss.HasDrift = ss.ChangedCount+ss.MissingCount+ss.ExtraCount > 0
		if ss.HasDrift {
			report.DriftedCount++
		} else {
			report.CleanCount++
		}
		report.Services = append(report.Services, ss)
	}
	return report
}

// Format returns a human-readable multi-line summary string.
func (r *SummaryReport) Format() string {
	if r == nil {
		return "no summary available"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Drift Summary — %s\n", r.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(&sb, "Services checked: %d | drifted: %d | clean: %d\n",
		r.TotalChecked, r.DriftedCount, r.CleanCount)
	for _, s := range r.Services {
		status := "OK"
		if s.HasDrift {
			status = fmt.Sprintf("DRIFT changed=%d missing=%d extra=%d",
				s.ChangedCount, s.MissingCount, s.ExtraCount)
		}
		fmt.Fprintf(&sb, "  %-30s %s\n", s.ServiceName, status)
	}
	return sb.String()
}
