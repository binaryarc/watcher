package keystore

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewStoreLoadsExistingKeys(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "keys.json")

	existing := []*KeyInfo{{
		Key:         "existing-key",
		Description: "preloaded",
		CreatedAt:   time.Now().UTC(),
	}}
	data, err := json.Marshal(existing)
	if err != nil {
		t.Fatalf("json.Marshal error = %v", err)
	}
	if err := os.WriteFile(filePath, data, 0600); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	store, err := NewStore(filePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	if !store.Validate("existing-key") {
		t.Fatalf("Validate(existing-key) = false, want true")
	}

	keys := store.List()
	if len(keys) != 1 || keys[0].Key != "existing-key" {
		t.Fatalf("List() = %+v, want single existing key", keys)
	}
}

func TestAddAndRemoveKey(t *testing.T) {
	t.Parallel()

	store := newTempStore(t)

	if err := store.Add("key1", "desc"); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	if !store.Validate("key1") {
		t.Fatalf("Validate(key1) = false, want true")
	}

	if err := store.Remove("key1"); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}

	if store.Validate("key1") {
		t.Fatalf("Validate(key1) = true after removal, want false")
	}
}

func TestAddDuplicateKeyFails(t *testing.T) {
	t.Parallel()

	store := newTempStore(t)

	if err := store.Add("dup", "first"); err != nil {
		t.Fatalf("Add(dup) error = %v", err)
	}

	if err := store.Add("dup", "second"); err == nil {
		t.Fatalf("Add(dup) second call expected error, got nil")
	}
}

func TestClearAndIsEmpty(t *testing.T) {
	t.Parallel()

	store := newTempStore(t)
	if err := store.Add("key", "desc"); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	if err := store.Clear(); err != nil {
		t.Fatalf("Clear() error = %v", err)
	}

	if !store.IsEmpty() {
		t.Fatalf("IsEmpty() = false, want true after Clear")
	}

	if store.Validate("key") {
		t.Fatalf("Validate(key) = true after Clear, want false")
	}
}

func TestValidateNoKeys(t *testing.T) {
	t.Parallel()

	store, err := NewStore("")
	if err != nil {
		t.Fatalf("NewStore(\"\") error = %v", err)
	}

	if store.Validate("any") {
		t.Fatalf("Validate(any) = true with no keys, want false")
	}
}

func newTempStore(t *testing.T) *Store {
	t.Helper()

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "keys.json")

	store, err := NewStore(filePath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	return store
}
