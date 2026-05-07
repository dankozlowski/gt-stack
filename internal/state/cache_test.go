package state

import (
	"path/filepath"
	"testing"
	"time"
)

func TestCache_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")

	c := &Cache{
		LastFetch: time.Date(2026, 5, 7, 10, 0, 0, 0, time.UTC),
		PRStates:  map[int]string{100: "MERGED", 102: "OPEN"},
	}
	if err := WriteCache(path, c); err != nil {
		t.Fatalf("write: %v", err)
	}
	got, err := ReadCache(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if got.PRStates[100] != "MERGED" || got.PRStates[102] != "OPEN" {
		t.Errorf("round-trip mismatch: %+v", got)
	}
}

func TestReadCache_MissingFileReturnsEmpty(t *testing.T) {
	got, err := ReadCache(filepath.Join(t.TempDir(), "nope.json"))
	if err != nil {
		t.Fatalf("ReadCache should return (empty, nil) for missing file, got: %v", err)
	}
	if len(got.PRStates) != 0 {
		t.Errorf("expected empty cache, got %+v", got)
	}
}
