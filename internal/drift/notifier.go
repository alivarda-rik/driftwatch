package drift

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// NotifyLevel represents the severity of a notification.
type NotifyLevel int

const (
	NotifyInfo NotifyLevel = iota
	NotifyWarn
	NotifyError
)

func (l NotifyLevel) String() string {
	switch l {
	case NotifyInfo:
		return "INFO"
	case NotifyWarn:
		return "WARN"
	case NotifyError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Notification holds a single notification event.
type Notification struct {
	Service   string
	Level     NotifyLevel
	Message   string
	Timestamp time.Time
}

func (n Notification) String() string {
	return fmt.Sprintf("[%s] %s (%s): %s",
		n.Timestamp.Format(time.RFC3339),
		n.Level,
		n.Service,
		n.Message,
	)
}

// Notifier dispatches notifications for drift events.
type Notifier struct {
	out    io.Writer
	prefix string
}

// NewNotifier creates a Notifier writing to the given writer.
// If w is nil, os.Stdout is used.
func NewNotifier(w io.Writer, prefix string) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w, prefix: prefix}
}

// Notify emits a notification for the given service and diff entries.
// Returns the number of notifications emitted.
func (n *Notifier) Notify(service string, entries []DiffEntry) int {
	if len(entries) == 0 {
		return 0
	}

	level := NotifyWarn
	for _, e := range entries {
		if e.Live == "" {
			level = NotifyError
			break
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d drift(s) detected", len(entries)))
	for _, e := range entries {
		sb.WriteString("; " + e.String())
	}

	notif := Notification{
		Service:   service,
		Level:     level,
		Message:   sb.String(),
		Timestamp: time.Now().UTC(),
	}

	line := notif.String()
	if n.prefix != "" {
		line = n.prefix + " " + line
	}
	fmt.Fprintln(n.out, line)
	return 1
}
