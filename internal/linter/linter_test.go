package linter

import (
	"testing"
)

func TestLint_EmptyValue_Warns(t *testing.T) {
	entries := map[string]string{
		"DATABASE_URL": "",
	}
	res := Lint(".env", entries)
	if len(res.Issues) == 0 {
		t.Fatal("expected at least one issue for empty value")
	}
	found := false
	for _, iss := range res.Issues {
		if iss.Key == "DATABASE_URL" && iss.Severity == "warn" {
			found = true
		}
	}
	if !found {
		t.Error("expected warn issue for DATABASE_URL empty value")
	}
}

func TestLint_LowercaseKey_Warns(t *testing.T) {
	entries := map[string]string{
		"myKey": "value",
	}
	res := Lint(".env", entries)
	found := false
	for _, iss := range res.Issues {
		if iss.Key == "myKey" && iss.Severity == "warn" {
			found = true
		}
	}
	if !found {
		t.Error("expected warn issue for non-uppercase key")
	}
}

func TestLint_UnquotedSpaces_Warns(t *testing.T) {
	entries := map[string]string{
		"GREETING": "hello world",
	}
	res := Lint(".env", entries)
	found := false
	for _, iss := range res.Issues {
		if iss.Key == "GREETING" && iss.Severity == "warn" {
			found = true
		}
	}
	if !found {
		t.Error("expected warn issue for unquoted spaces")
	}
}

func TestLint_QuotedSpaces_NoWarn(t *testing.T) {
	entries := map[string]string{
		"GREETING": `"hello world"`,
	}
	res := Lint(".env", entries)
	for _, iss := range res.Issues {
		if iss.Key == "GREETING" {
			t.Errorf("unexpected issue for quoted value: %s", iss.Message)
		}
	}
}

func TestLint_CleanEntry_NoIssues(t *testing.T) {
	entries := map[string]string{
		"API_KEY": "abc123",
	}
	res := Lint(".env", entries)
	if len(res.Issues) != 0 {
		t.Errorf("expected no issues, got %d", len(res.Issues))
	}
}

func TestResult_HasErrors_False(t *testing.T) {
	res := Result{
		Issues: []Issue{{Severity: "warn"}},
	}
	if res.HasErrors() {
		t.Error("expected HasErrors to be false when only warnings present")
	}
}
