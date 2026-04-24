package drift

import "fmt"

// DiffEntry represents a single difference between declared and live state.
type DiffEntry struct {
	Key      string
	Declared interface{}
	Live     interface{}
	Kind     DiffKind
}

// DiffKind categorises the nature of a drift difference.
type DiffKind string

const (
	DiffKindChanged DiffKind = "changed"
	DiffKindMissing DiffKind = "missing" // present in declared, absent in live
	DiffKindExtra   DiffKind = "extra"   // present in live, absent in declared
)

// String returns a human-readable representation of a DiffEntry.
func (d DiffEntry) String() string {
	switch d.Kind {
	case DiffKindChanged:
		return fmt.Sprintf("[changed] %s: declared=%v live=%v", d.Key, d.Declared, d.Live)
	case DiffKindMissing:
		return fmt.Sprintf("[missing] %s: declared=%v (not found in live)", d.Key, d.Declared)
	case DiffKindExtra:
		return fmt.Sprintf("[extra]   %s: live=%v (not in declared)", d.Key, d.Live)
	default:
		return fmt.Sprintf("[unknown] %s", d.Key)
	}
}

// Diff compares two flat key-value maps and returns a slice of DiffEntry
// describing every discrepancy. Both maps must use string keys.
func Diff(declared, live map[string]interface{}) []DiffEntry {
	var entries []DiffEntry

	for k, dv := range declared {
		lv, ok := live[k]
		if !ok {
			entries = append(entries, DiffEntry{Key: k, Declared: dv, Kind: DiffKindMissing})
			continue
		}
		if fmt.Sprintf("%v", dv) != fmt.Sprintf("%v", lv) {
			entries = append(entries, DiffEntry{Key: k, Declared: dv, Live: lv, Kind: DiffKindChanged})
		}
	}

	for k, lv := range live {
		if _, ok := declared[k]; !ok {
			entries = append(entries, DiffEntry{Key: k, Live: lv, Kind: DiffKindExtra})
		}
	}

	return entries
}
