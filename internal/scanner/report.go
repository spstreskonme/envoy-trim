package scanner

import "fmt"

// Report summarises which env keys were found in source files.
type Report struct {
	// Refs maps each env key to the list of file:line locations where it appears.
	Refs map[string][]string
}

// NewReport constructs a Report from a pre-built refs map.
func NewReport(refs map[string][]string) *Report {
	if refs == nil {
		refs = make(map[string][]string)
	}
	return &Report{Refs: refs}
}

// IsUsed returns true if the given key appears at least once in the scanned sources.
func (r *Report) IsUsed(key string) bool {
	locs, ok := r.Refs[key]
	return ok && len(locs) > 0
}

// Unused returns all keys from the provided slice that are NOT present in the report.
func (r *Report) Unused(keys []string) []string {
	var out []string
	for _, k := range keys {
		if !r.IsUsed(k) {
			out = append(out, k)
		}
	}
	return out
}

// BuildReport constructs a Report by scanning the given directory for references
// to the provided env keys.
func BuildReport(dir string, keys []string) (*Report, error) {
	s, err := New(dir)
	if err != nil {
		return nil, fmt.Errorf("creating scanner: %w", err)
	}

	refs, err := s.ScanDir(keys)
	if err != nil {
		return nil, fmt.Errorf("scanning directory: %w", err)
	}

	return NewReport(refs), nil
}
