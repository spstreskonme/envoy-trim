package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/envoy-trim/internal/reporter"
	"github.com/envoy-trim/internal/scanner"
)

func buildReport(used, unused []string) *scanner.Report {
	r := &scanner.Report{
		Used:   make(map[string]bool),
		Unused: make(map[string]bool),
	}
	for _, k := range used {
		r.Used[k] = true
	}
	for _, k := range unused {
		r.Unused[k] = true
	}
	return r
}

func TestPrintText_ShowsUnusedKeys(t *testing.T) {
	var buf bytes.Buffer
	opts := reporter.DefaultOptions()
	opts.Writer = &buf

	rep := buildReport([]string{"DB_HOST"}, []string{"OLD_KEY", "LEGACY_TOKEN"})
	if err := reporter.Print(rep, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "OLD_KEY") {
		t.Errorf("expected OLD_KEY in output, got:\n%s", out)
	}
	if !strings.Contains(out, "LEGACY_TOKEN") {
		t.Errorf("expected LEGACY_TOKEN in output, got:\n%s", out)
	}
}

func TestPrintText_NoneUnused(t *testing.T) {
	var buf bytes.Buffer
	opts := reporter.DefaultOptions()
	opts.Writer = &buf

	rep := buildReport([]string{"DB_HOST"}, []string{})
	if err := reporter.Print(rep, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "No unused") {
		t.Errorf("expected no-unused message, got:\n%s", buf.String())
	}
}

func TestPrintJSON_ValidStructure(t *testing.T) {
	var buf bytes.Buffer
	opts := reporter.Options{
		Format: reporter.FormatJSON,
		Writer: &buf,
	}

	rep := buildReport([]string{"DB_HOST", "PORT"}, []string{"OLD_KEY"})
	if err := reporter.Print(rep, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	summary, ok := result["summary"].(map[string]interface{})
	if !ok {
		t.Fatal("missing summary field")
	}
	if int(summary["total"].(float64)) != 3 {
		t.Errorf("expected total=3, got %v", summary["total"])
	}
	if int(summary["unused"].(float64)) != 1 {
		t.Errorf("expected unused=1, got %v", summary["unused"])
	}
}
