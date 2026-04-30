package differ

import (
	"sort"
)

// Delta represents the difference between two sets of env keys.
type Delta struct {
	Added   []string // keys present in next but not in base
	Removed []string // keys present in base but not in next
	Kept    []string // keys present in both
}

// Compare computes the delta between a base set of keys and a next set of keys.
// base is typically the current .env file keys; next is the pruned/updated set.
func Compare(base, next []string) Delta {
	baseSet := toSet(base)
	nextSet := toSet(next)

	var added, removed, kept []string

	for k := range nextSet {
		if _, ok := baseSet[k]; ok {
			kept = append(kept, k)
		} else {
			added = append(added, k)
		}
	}

	for k := range baseSet {
		if _, ok := nextSet[k]; !ok {
			removed = append(removed, k)
		}
	}

	sort.Strings(added)
	sort.Strings(removed)
	sort.Strings(kept)

	return Delta{
		Added:   added,
		Removed: removed,
		Kept:    kept,
	}
}

func toSet(keys []string) map[string]struct{} {
	s := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		s[k] = struct{}{}
	}
	return s
}
