package drift

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// ExportFormat represents the output format for drift reports.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatText ExportFormat = "text"
)

// ExportRecord holds a drift report snapshot suitable for export.
type ExportRecord struct {
	Service   string      `json:"service"`
	Timestamp time.Time   `json:"timestamp"`
	DriftCount int        `json:"drift_count"`
	Entries   []DiffEntry `json:"entries"`
}

// Exporter writes drift results to an io.Writer in a specified format.
type Exporter struct {
	format ExportFormat
	writer io.Writer
}

// NewExporter creates an Exporter that writes to w using the given format.
func NewExporter(format ExportFormat, w io.Writer) *Exporter {
	return &Exporter{format: format, writer: w}
}

// Export serialises the given DiffEntry slice for the named service.
func (e *Exporter) Export(service string, entries []DiffEntry) error {
	if e == nil || e.writer == nil {
		return fmt.Errorf("exporter: nil exporter or writer")
	}
	record := ExportRecord{
		Service:    service,
		Timestamp:  time.Now().UTC(),
		DriftCount: len(entries),
		Entries:    entries,
	}
	switch e.format {
	case FormatJSON:
		return e.writeJSON(record)
	case FormatText:
		return e.writeText(record)
	default:
		return fmt.Errorf("exporter: unsupported format %q", e.format)
	}
}

func (e *Exporter) writeJSON(r ExportRecord) error {
	enc := json.NewEncoder(e.writer)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func (e *Exporter) writeText(r ExportRecord) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Service:    %s\n", r.Service))
	sb.WriteString(fmt.Sprintf("Timestamp:  %s\n", r.Timestamp.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Drift Count: %d\n", r.DriftCount))
	if len(r.Entries) == 0 {
		sb.WriteString("Status:     no drift detected\n")
	} else {
		sb.WriteString("Diffs:\n")
		for _, entry := range r.Entries {
			sb.WriteString(fmt.Sprintf("  - %s\n", entry.String()))
		}
	}
	_, err := fmt.Fprint(e.writer, sb.String())
	return err
}
