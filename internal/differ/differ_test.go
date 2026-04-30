package differ_test

import (
	"testing"

	"github.com/yourorg/envoy-trim/internal/differ"
)

func TestCompare_AddsAndRemoves(t *testing.T) {
	base := []string{"DB_HOST", "DB_PORT", "SECRET_KEY"}
	next := []string{"DB_HOST", "API_KEY"}

	delta := differ.Compare(base, next)

	if len(delta.Added) != 1 || delta.Added[0] != "API_KEY" {
		t.Errorf("expected Added=[API_KEY], got %v", delta.Added)
	}
	if len(delta.Removed) != 2 {
		t.Errorf("expected 2 removed keys, got %v", delta.Removed)
	}
	if len(delta.Kept) != 1 || delta.Kept[0] != "DB_HOST" {
		t.Errorf("expected Kept=[DB_HOST], got %v", delta.Kept)
	}
}

func TestCompare_NoDifference(t *testing.T) {
	keys := []string{"FOO", "BAR", "BAZ"}
	delta := differ.Compare(keys, keys)

	if len(delta.Added) != 0 {
		t.Errorf("expected no additions, got %v", delta.Added)
	}
	if len(delta.Removed) != 0 {
		t.Errorf("expected no removals, got %v", delta.Removed)
	}
	if len(delta.Kept) != 3 {
		t.Errorf("expected 3 kept keys, got %v", delta.Kept)
	}
}

func TestCompare_EmptyBase(t *testing.T) {
	delta := differ.Compare([]string{}, []string{"NEW_KEY"})

	if len(delta.Added) != 1 || delta.Added[0] != "NEW_KEY" {
		t.Errorf("expected Added=[NEW_KEY], got %v", delta.Added)
	}
	if len(delta.Removed) != 0 {
		t.Errorf("expected no removals, got %v", delta.Removed)
	}
}

func TestCompare_EmptyNext(t *testing.T) {
	delta := differ.Compare([]string{"OLD_KEY"}, []string{})

	if len(delta.Removed) != 1 || delta.Removed[0] != "OLD_KEY" {
		t.Errorf("expected Removed=[OLD_KEY], got %v", delta.Removed)
	}
	if len(delta.Added) != 0 {
		t.Errorf("expected no additions, got %v", delta.Added)
	}
}
