package scanner_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envoy-trim/internal/scanner"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
	return p
}

func TestScanDir_FindsRefs(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "main.go", `package main
import "os"
func main() { os.Getenv("DATABASE_URL"); os.Getenv("SECRET_KEY") }
`)

	s := scanner.New()
	results, err := s.ScanDir(dir, []string{"DATABASE_URL", "SECRET_KEY", "UNUSED_VAR"})
	if err != nil {
		t.Fatalf("ScanDir: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if len(results[0].Refs) != 2 {
		t.Errorf("expected 2 refs, got %d: %v", len(results[0].Refs), results[0].Refs)
	}
}

func TestScanDir_SkipsUnknownExtensions(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "data.csv", "DATABASE_URL,SECRET_KEY\n")

	s := scanner.New()
	results, err := s.ScanDir(dir, []string{"DATABASE_URL"})
	if err != nil {
		t.Fatalf("ScanDir: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results for .csv file, got %d", len(results))
	}
}

func TestScanDir_SkipsNodeModules(t *testing.T) {
	dir := t.TempDir()
	nm := filepath.Join(dir, "node_modules")
	if err := os.Mkdir(nm, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, nm, "index.js", "const x = process.env.DATABASE_URL;")

	s := scanner.New()
	results, err := s.ScanDir(dir, []string{"DATABASE_URL"})
	if err != nil {
		t.Fatalf("ScanDir: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected node_modules to be skipped, got %d results", len(results))
	}
}

func TestScanDir_EmptyKeys(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "app.py", "import os\nval = os.environ['DATABASE_URL']\n")

	s := scanner.New()
	results, err := s.ScanDir(dir, []string{})
	if err != nil {
		t.Fatalf("ScanDir: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results for empty keys, got %d", len(results))
	}
}
