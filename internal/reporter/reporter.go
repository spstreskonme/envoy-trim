package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/envoy-trim/internal/scanner"
)

// Format represents the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Options configures report output.
type Options struct {
	Format  Format
	Writer  io.Writer
	Verbose bool
}

// DefaultOptions returns options writing to stdout in text format.
func DefaultOptions() Options {
	return Options{
		Format:  FormatText,
		Writer:  os.Stdout,
		Verbose: false,
	}
}

// Print writes a human-readable or structured report from a BuildReport result.
func Print(report *scanner.Report, opts Options) error {
	switch opts.Format {
	case FormatJSON:
		return printJSON(report, opts.Writer)
	default:
		return printText(report, opts.Writer, opts.Verbose)
	}
}

func printText(report *scanner.Report, w io.Writer, verbose bool) error {
	unused := sortedKeys(report.Unused)
	used := sortedKeys(report.Used)

	fmt.Fprintf(w, "envoy-trim report\n")
	fmt.Fprintf(w, "==================\n")
	fmt.Fprintf(w, "Total keys : %d\n", len(used)+len(unused))
	fmt.Fprintf(w, "Used       : %d\n", len(used))
	fmt.Fprintf(w, "Unused     : %d\n\n", len(unused))

	if len(unused) == 0 {
		fmt.Fprintln(w, "✓ No unused environment variables found.")
		return nil
	}

	fmt.Fprintln(w, "Unused keys:")
	for _, k := range unused {
		fmt.Fprintf(w, "  - %s\n", k)
	}

	if verbose && len(used) > 0 {
		fmt.Fprintln(w, "\nUsed keys:")
		for _, k := range used {
			fmt.Fprintf(w, "  + %s\n", k)
		}
	}
	return nil
}

// jsonReport is the structure used for JSON-formatted output.
type jsonReport struct {
	TotalKeys int      `json:"total_keys"`
	Used      []string `json:"used"`
	Unused    []string `json:"unused"`
}

func printJSON(report *scanner.Report, w io.Writer) error {
	unused := sortedKeys(report.Unused)
	used := sortedKeys(report.Used)

	payload := jsonReport{
		TotalKeys: len(used) + len(unused),
		Used:      used,
		Unused:    unused,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		return fmt.Errorf("reporter: failed to encode JSON: %w", err)
	}
	return nil
}

func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
