package envparser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvEntry represents a single key-value pair parsed from a .env file.
type EnvEntry struct {
	Key      string
	Value    string
	Line     int
	Comment  bool
	Raw      string
}

// ParseFile reads a .env file and returns a slice of EnvEntry records.
func ParseFile(path string) ([]EnvEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("envparser: open %q: %w", path, err)
	}
	defer f.Close()

	var entries []EnvEntry
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)

		// Skip blank lines
		if trimmed == "" {
			continue
		}

		// Comment lines
		if strings.HasPrefix(trimmed, "#") {
			entries = append(entries, EnvEntry{
				Line:    lineNum,
				Comment: true,
				Raw:     raw,
			})
			continue
		}

		// KEY=VALUE or KEY="VALUE"
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("envparser: %q line %d: malformed entry %q", path, lineNum, trimmed)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = stripQuotes(value)

		entries = append(entries, EnvEntry{
			Key:   key,
			Value: value,
			Line:  lineNum,
			Raw:   raw,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("envparser: scan %q: %w", path, err)
	}

	return entries, nil
}

// Keys returns only the variable keys from a list of entries (non-comment).
func Keys(entries []EnvEntry) []string {
	var keys []string
	for _, e := range entries {
		if !e.Comment {
			keys = append(keys, e.Key)
		}
	}
	return keys
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
