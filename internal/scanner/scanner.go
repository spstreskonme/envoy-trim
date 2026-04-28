package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

// Result holds the scan results for a single source file.
type Result struct {
	FilePath string
	Refs     []string
}

// Scanner searches source files for environment variable references.
type Scanner struct {
	Extensions []string
}

// New returns a Scanner with a default set of source file extensions.
func New() *Scanner {
	return &Scanner{
		Extensions: []string{".go", ".js", ".ts", ".py", ".sh", ".yaml", ".yml"},
	}
}

// ScanDir walks dir recursively and returns all env-var references found
// in files whose extensions match s.Extensions.
func (s *Scanner) ScanDir(dir string, keys []string) ([]Result, error) {
	var results []Result

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" || d.Name() == "node_modules" || d.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		if !s.hasTrackedExtension(path) {
			return nil
		}

		refs, err := findRefs(path, keys)
		if err != nil {
			return err
		}
		if len(refs) > 0 {
			results = append(results, Result{FilePath: path, Refs: refs})
		}
		return nil
	})

	return results, err
}

func (s *Scanner) hasTrackedExtension(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, e := range s.Extensions {
		if e == ext {
			return true
		}
	}
	return false
}

// findRefs returns the subset of keys that appear in the file at path.
func findRefs(path string, keys []string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(data)

	var found []string
	for _, key := range keys {
		if strings.Contains(content, key) {
			found = append(found, key)
		}
	}
	return found, nil
}
