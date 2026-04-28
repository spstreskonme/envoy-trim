package scanner

import (
	"fmt"
	"io"
	"sort"
)

// UsageReport summarises which env keys are used and which are unused.
type UsageReport struct {
	Used   []string
	Unused []string
}

// BuildReport cross-references envKeys against scanner Results to produce
// a UsageReport.
func BuildReport(envKeys []string, results []Result) UsageReport {
	usedSet := make(map[string]struct{})
	for _, r := range results {
		for _, ref := range r.Refs {
			usedSet[ref] = struct{}{}
		}
	}

	var used, unused []string
	for _, k := range envKeys {
		if _, ok := usedSet[k]; ok {
			used = append(used, k)
		} else {
			unused = append(unused, k)
		}
	}

	sort.Strings(used)
	sort.Strings(unused)

	return UsageReport{Used: used, Unused: unused}
}

// Print writes a human-readable summary of the report to w.
func (r UsageReport) Print(w io.Writer) {
	fmt.Fprintf(w, "Used variables (%d):\n", len(r.Used))
	for _, k := range r.Used {
		fmt.Fprintf(w, "  ✔ %s\n", k)
	}

	fmt.Fprintf(w, "\nUnused variables (%d):\n", len(r.Unused))
	for _, k := range r.Unused {
		fmt.Fprintf(w, "  ✘ %s\n", k)
	}
}
