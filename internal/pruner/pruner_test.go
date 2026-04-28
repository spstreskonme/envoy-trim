package pruner_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-trim/internal/pruner"
	"github.com/envoy-trim/internal/scanner"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestPrune_DryRun_DoesNotModifyFile(t *testing.T) {
	envContent := "DB_HOST=localhost\nUNUSED_KEY=foo\nAPP_PORT=8080\n"
	path := writeTempEnv(t, envContent)

	report := scanner.NewReport(map[string][]string{
		"DB_HOST":  {"main.go:10"},
		"APP_PORT": {"server.go:5"},
	})

	res, err := pruner.Prune(path, report, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Removed) != 1 || res.Removed[0] != "UNUSED_KEY" {
		t.Errorf("expected UNUSED_KEY removed, got %v", res.Removed)
	}

	got, _ := os.ReadFile(path)
	if string(got) != envContent {
		t.Errorf("dry run should not modify file")
	}
}

func TestPrune_WritesFile_WhenNotDryRun(t *testing.T) {
	envContent := "DB_HOST=localhost\nUNUSED_KEY=foo\nAPP_PORT=8080\n"
	path := writeTempEnv(t, envContent)

	report := scanner.NewReport(map[string][]string{
		"DB_HOST":  {"main.go:10"},
		"APP_PORT": {"server.go:5"},
	})

	res, err := pruner.Prune(path, report, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Removed) != 1 {
		t.Errorf("expected 1 removed key, got %d", len(res.Removed))
	}

	got, _ := os.ReadFile(path)
	if contains(string(got), "UNUSED_KEY") {
		t.Errorf("pruned file should not contain UNUSED_KEY")
	}
}

func TestPrune_PreservesComments(t *testing.T) {
	envContent := "# database config\nDB_HOST=localhost\n"
	path := writeTempEnv(t, envContent)

	report := scanner.NewReport(map[string][]string{
		"DB_HOST": {"main.go:1"},
	})

	_, err := pruner.Prune(path, report, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := os.ReadFile(path)
	if !contains(string(got), "# database config") {
		t.Errorf("comments should be preserved after pruning")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
