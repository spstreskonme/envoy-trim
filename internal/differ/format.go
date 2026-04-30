package differ

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// FormatText writes a human-readable diff summary to w.
func FormatText(w io.Writer, d Delta) {
	if len(d.Added) == 0 && len(d.Removed) == 0 {
		fmt.Fprintln(w, "No changes detected.")
		return
	}

	if len(d.Added) > 0 {
		fmt.Fprintf(w, "+ Added (%d):\n", len(d.Added))
		for _, k := range d.Added {
			fmt.Fprintf(w, "  + %s\n", k)
		}
	}

	if len(d.Removed) > 0 {
		fmt.Fprintf(w, "- Removed (%d):\n", len(d.Removed))
		for _, k := range d.Removed {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	fmt.Fprintf(w, "  Kept: %d key(s) unchanged\n", len(d.Kept))
}

// FormatJSON writes a JSON-encoded diff to w.
func FormatJSON(w io.Writer, d Delta) error {
	type jsonDelta struct {
		Added   []string `json:"added"`
		Removed []string `json:"removed"`
		Kept    []string `json:"kept"`
	}

	out := jsonDelta{
		Added:   nullSafe(d.Added),
		Removed: nullSafe(d.Removed),
		Kept:    nullSafe(d.Kept),
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

// Summary returns a one-line summary string of the delta.
func Summary(d Delta) string {
	parts := []string{}
	if len(d.Added) > 0 {
		parts = append(parts, fmt.Sprintf("+%d added", len(d.Added)))
	}
	if len(d.Removed) > 0 {
		parts = append(parts, fmt.Sprintf("-%d removed", len(d.Removed)))
	}
	if len(parts) == 0 {
		return "no changes"
	}
	return strings.Join(parts, ", ")
}

func nullSafe(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
