package drift

import (
	"fmt"
	"io"
	"strings"
)

// Reporter formats and writes drift reports to an output writer.
type Reporter struct {
	out io.Writer
}

// NewReporter creates a Reporter that writes to the given writer.
func NewReporter(out io.Writer) *Reporter {
	return &Reporter{out: out}
}

// Print writes a human-readable drift report to the output writer.
func (r *Reporter) Print(report *Report) {
	if report == nil {
		fmt.Fprintln(r.out, "no report available")
		return
	}

	fmt.Fprintf(r.out, "Service: %s\n", report.ServiceName)
	fmt.Fprintf(r.out, "%s\n", strings.Repeat("-", 40))

	if !report.HasDrift {
		fmt.Fprintln(r.out, "✓ No drift detected")
		return
	}

	fmt.Fprintln(r.out, "✗ Drift detected:")
	for _, res := range report.Results {
		if res.Drifted {
			fmt.Fprintf(r.out, "  [DRIFT] %s: expected=%v, actual=%v\n",
				res.Field, res.Expected, res.Actual)
		}
	}

	fmt.Fprintf(r.out, "%s\n", strings.Repeat("-", 40))
	driftCount := 0
	for _, res := range report.Results {
		if res.Drifted {
			driftCount++
		}
	}
	fmt.Fprintf(r.out, "Total drifted fields: %d\n", driftCount)
}
