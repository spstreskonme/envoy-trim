package differ_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/envoy-trim/internal/differ"
)

func TestFormatText_ShowsAddedAndRemoved(t *testing.T) {
	d := differ.Delta{
		Added:   []string{"NEW_KEY"},
		Removed: []string{"OLD_KEY"},
		Kept:    []string{"STABLE"},
	}

	var buf bytes.Buffer
	differ.FormatText(&buf, d)
	out := buf.String()

	if !strings.Contains(out, "+ NEW_KEY") {
		t.Errorf("expected '+ NEW_KEY' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "- OLD_KEY") {
		t.Errorf("expected '- OLD_KEY' in output, got:\n%s", out)
	}
}

func TestFormatText_NoChanges(t *testing.T) {
	d := differ.Delta{Kept: []string{"FOO"}}
	var buf bytes.Buffer
	differ.FormatText(&buf, d)
	if !strings.Contains(buf.String(), "No changes detected") {
		t.Errorf("expected no-change message, got: %s", buf.String())
	}
}

func TestFormatJSON_ValidStructure(t *testing.T) {
	d := differ.Delta{
		Added:   []string{"A"},
		Removed: []string{"B"},
		Kept:    []string{"C"},
	}

	var buf bytes.Buffer
	if err := differ.FormatJSON(&buf, d); err != nil {
		t.Fatalf("FormatJSON error: %v", err)
	}

	var out map[string][]string
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out["added"]) != 1 || out["added"][0] != "A" {
		t.Errorf("unexpected added: %v", out["added"])
	}
}

func TestSummary(t *testing.T) {
	d := differ.Delta{Added: []string{"X", "Y"}, Removed: []string{"Z"}}
	s := differ.Summary(d)
	if !strings.Contains(s, "+2 added") || !strings.Contains(s, "-1 removed") {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummary_NoChanges(t *testing.T) {
	d := differ.Delta{Kept: []string{"A"}}
	if differ.Summary(d) != "no changes" {
		t.Errorf("expected 'no changes'")
	}
}
