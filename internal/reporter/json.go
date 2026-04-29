package reporter

import (
	"encoding/json"
	"io"
	"sort"

	"github.com/envoy-trim/internal/scanner"
)

// jsonReport is the serialisable representation of a scan report.
type jsonReport struct {
	Summary jsonSummary `json:"summary"`
	Used    []string    `json:"used"`
	Unused  []string    `json:"unused"`
}

type jsonSummary struct {
	Total  int `json:"total"`
	Used   int `json:"used"`
	Unused int `json:"unused"`
}

func printJSON(report *scanner.Report, w io.Writer) error {
	usedKeys := mapToSortedSlice(report.Used)
	unusedKeys := mapToSortedSlice(report.Unused)

	out := jsonReport{
		Summary: jsonSummary{
			Total:  len(usedKeys) + len(unusedKeys),
			Used:   len(usedKeys),
			Unused: len(unusedKeys),
		},
		Used:   usedKeys,
		Unused: unusedKeys,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func mapToSortedSlice(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
