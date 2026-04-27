package envparser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), ".env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParseFile_BasicEntries(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\n")

	entries, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "DB_HOST" || entries[0].Value != "localhost" {
		t.Errorf("entry 0: got key=%q value=%q", entries[0].Key, entries[0].Value)
	}
	if entries[1].Key != "DB_PORT" || entries[1].Value != "5432" {
		t.Errorf("entry 1: got key=%q value=%q", entries[1].Key, entries[1].Value)
	}
}

func TestParseFile_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `API_KEY="secret123"` + "\n" + `MSG='hello world'` + "\n")

	entries, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Value != "secret123" {
		t.Errorf("expected unquoted value, got %q", entries[0].Value)
	}
	if entries[1].Value != "hello world" {
		t.Errorf("expected unquoted value, got %q", entries[1].Value)
	}
}

func TestParseFile_CommentsSkipped(t *testing.T) {
	path := writeTempEnv(t, "# this is a comment\nFOO=bar\n")

	entries, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries (1 comment + 1 var), got %d", len(entries))
	}
	if !entries[0].Comment {
		t.Errorf("expected first entry to be a comment")
	}
}

func TestParseFile_MalformedLine(t *testing.T) {
	path := writeTempEnv(t, "BADLINE\n")

	_, err := ParseFile(path)
	if err == nil {
		t.Fatal("expected error for malformed line, got nil")
	}
}

func TestKeys(t *testing.T) {
	entries := []EnvEntry{
		{Key: "FOO", Comment: false},
		{Comment: true},
		{Key: "BAR", Comment: false},
	}
	keys := Keys(entries)
	if len(keys) != 2 || keys[0] != "FOO" || keys[1] != "BAR" {
		t.Errorf("unexpected keys: %v", keys)
	}
}
