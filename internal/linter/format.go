package linter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// FormatText writes a human-readable lint report to w.
func FormatText(w io.Writer, results []Result) {
	any := false
	for _, r := range results {
		if len(r.Issues) == 0 {
			continue
		}
		any = true
		fmt.Fprintf(w, "\n%s\n", r.File)
		fmt.Fprintf(w, "%s\n", strings.Repeat("-", len(r.File)))
		for _, iss := range r.Issues {
			icon := "⚠"
			if iss.Severity == "error" {
				icon = "✖"
			}
			fmt.Fprintf(w, "  %s  [%s] %s\n", icon, iss.Severity, iss.Message)
		}
	}
	if !any {
		fmt.Fprintln(w, "No lint issues found.")
	}
}

// jsonIssue mirrors Issue for JSON serialisation.
type jsonIssue struct {
	Key      string `json:"key"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

// jsonResult mirrors Result for JSON serialisation.
type jsonResult struct {
	File   string      `json:"file"`
	Issues []jsonIssue `json:"issues"`
}

// FormatJSON writes a JSON lint report to w.
func FormatJSON(w io.Writer, results []Result) error {
	out := make([]jsonResult, 0, len(results))
	for _, r := range results {
		jr := jsonResult{File: r.File, Issues: []jsonIssue{}}
		for _, iss := range r.Issues {
			jr.Issues = append(jr.Issues, jsonIssue{
				Key:      iss.Key,
				Severity: iss.Severity,
				Message:  iss.Message,
			})
		}
		out = append(out, jr)
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
