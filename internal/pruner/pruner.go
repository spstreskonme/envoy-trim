package pruner

import (
	"fmt"
	"os"
	"strings"

	"github.com/envoy-trim/internal/scanner"
)

// Result holds the outcome of a prune operation.
type Result struct {
	Removed []string
	Kept    []string
	DryRun  bool
}

// Summary returns a human-readable summary of the prune result.
func (r *Result) Summary() string {
	mode := "applied"
	if r.DryRun {
		mode = "dry-run"
	}
	return fmt.Sprintf("[%s] removed: %d, kept: %d", mode, len(r.Removed), len(r.Kept))
}

// Prune removes unused keys from the given env file based on the report.
// If dryRun is true, no changes are written to disk.
func Prune(envFilePath string, report *scanner.Report, dryRun bool) (*Result, error) {
	data, err := os.ReadFile(envFilePath)
	if err != nil {
		return nil, fmt.Errorf("reading env file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	result := &Result{DryRun: dryRun}
	var kept []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			kept = append(kept, line)
			continue
		}

		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) != 2 {
			kept = append(kept, line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		if report.IsUsed(key) {
			kept = append(kept, line)
			result.Kept = append(result.Kept, key)
		} else {
			result.Removed = append(result.Removed, key)
		}
	}

	if !dryRun {
		output := strings.Join(kept, "\n")
		if err := os.WriteFile(envFilePath, []byte(output), 0644); err != nil {
			return nil, fmt.Errorf("writing pruned env file: %w", err)
		}
	}

	return result, nil
}
