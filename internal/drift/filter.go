package drift

import "strings"

// FilterOptions controls which drift entries are included in results.
type FilterOptions struct {
	// OnlyChanged limits results to entries where the value changed.
	OnlyChanged bool
	// OnlyMissing limits results to entries missing from the live state.
	OnlyMissing bool
	// OnlyExtra limits results to entries present in live but not declared.
	OnlyExtra bool
	// KeyPrefix filters entries whose key starts with the given prefix.
	KeyPrefix string
}

// Filter applies FilterOptions to a slice of DiffEntry values, returning
// only the entries that satisfy all active constraints.
func Filter(entries []DiffEntry, opts FilterOptions) []DiffEntry {
	var out []DiffEntry
	for _, e := range entries {
		if opts.KeyPrefix != "" && !strings.HasPrefix(e.Key, opts.KeyPrefix) {
			continue
		}
		if opts.OnlyChanged && e.Declared == "" && e.Live == "" {
			continue
		}
		if opts.OnlyMissing && e.Live != "" {
			continue
		}
		if opts.OnlyExtra && e.Declared != "" {
			continue
		}
		if opts.OnlyChanged && !opts.OnlyMissing && !opts.OnlyExtra {
			if e.Declared == "" || e.Live == "" {
				continue
			}
		}
		out = append(out, e)
	}
	return out
}

// FilterByKeys returns only the DiffEntry values whose Key is present in the
// provided set.
func FilterByKeys(entries []DiffEntry, keys []string) []DiffEntry {
	set := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		set[k] = struct{}{}
	}
	var out []DiffEntry
	for _, e := range entries {
		if _, ok := set[e.Key]; ok {
			out = append(out, e)
		}
	}
	return out
}
