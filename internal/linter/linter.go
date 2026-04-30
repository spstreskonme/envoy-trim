package linter

import (
	"fmt"
	"sort"
	"strings"
)

// Issue represents a linting problem found in a .env file.
type Issue struct {
	Key      string
	Line     int
	Severity string // "warn" or "error"
	Message  string
}

// Result holds all issues found for a given .env file.
type Result struct {
	File   string
	Issues []Issue
}

// Lint checks the parsed entries of a .env file for common problems.
// entries is a map of key -> raw value (as returned by envparser.ParseFile).
func Lint(file string, entries map[string]string) Result {
	result := Result{File: file}

	// Collect keys in deterministic order for stable output.
	keys := make([]string, 0, len(entries))
	for k := range entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		val := entries[key]

		// Warn on empty values.
		if strings.TrimSpace(val) == "" {
			result.Issues = append(result.Issues, Issue{
				Key:      key,
				Severity: "warn",
				Message:  fmt.Sprintf("key %q has an empty value", key),
			})
		}

		// Warn on keys that are not uppercase.
		if key != strings.ToUpper(key) {
			result.Issues = append(result.Issues, Issue{
				Key:      key,
				Severity: "warn",
				Message:  fmt.Sprintf("key %q is not uppercase", key),
			})
		}

		// Warn on values that contain unquoted spaces (possible mistake).
		if !isQuoted(val) && strings.Contains(val, " ") {
			result.Issues = append(result.Issues, Issue{
				Key:      key,
				Severity: "warn",
				Message:  fmt.Sprintf("key %q has an unquoted value containing spaces", key),
			})
		}
	}

	return result
}

// HasErrors returns true if any issue has severity "error".
func (r Result) HasErrors() bool {
	for _, i := range r.Issues {
		if i.Severity == "error" {
			return true
		}
	}
	return false
}

func isQuoted(s string) bool {
	return (strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`)) ||
		(strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'"))
}
