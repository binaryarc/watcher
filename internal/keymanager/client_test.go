package keymanager

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewManagerCreatesDir(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	keysDir := filepath.Join(tempDir, "keys")

	manager, err := NewManager(keysDir)
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}

	if manager.keysDir != keysDir {
		t.Fatalf("manager.keysDir = %s, want %s", manager.keysDir, keysDir)
	}

	if info, err := os.Stat(keysDir); err != nil || !info.IsDir() {
		t.Fatalf("keys directory not created: %v", err)
	}
}

func TestGenerateProducesPrefixedKey(t *testing.T) {
	t.Parallel()

	manager, err := NewManager(t.TempDir())
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}

	key, err := manager.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.HasPrefix(key, KeyPrefix) {
		t.Fatalf("key %q does not start with prefix %q", key, KeyPrefix)
	}

	if len(key) <= len(KeyPrefix) {
		t.Fatalf("key %q too short", key)
	}
}

func TestSaveLoadDelete(t *testing.T) {
	t.Parallel()

	manager, err := NewManager(t.TempDir())
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}

	const name = "my-key"
	const value = "watcher_test_key"

	if err := manager.Save(name, value); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, err := manager.Load(name)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if got != value {
		t.Fatalf("Load() = %q, want %q", got, value)
	}

	if err := manager.Delete(name); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	if _, err := manager.Load(name); err == nil {
		t.Fatalf("Load() after Delete() expected error, got nil")
	}
}

func TestListKeys(t *testing.T) {
	t.Parallel()

	manager, err := NewManager(t.TempDir())
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}

	keys := []string{"alpha", "beta", "gamma"}
	for _, name := range keys {
		if err := manager.Save(name, "value-"+name); err != nil {
			t.Fatalf("Save(%s) error = %v", name, err)
		}
	}

	list, err := manager.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(list) != len(keys) {
		t.Fatalf("List() len = %d, want %d", len(list), len(keys))
	}

	want := make(map[string]struct{})
	for _, name := range keys {
		want[name] = struct{}{}
	}

	for _, name := range list {
		if _, ok := want[name]; !ok {
			t.Fatalf("unexpected key %q in list", name)
		}
		delete(want, name)
	}

	if len(want) != 0 {
		t.Fatalf("missing keys in list: %v", want)
	}
}

func TestGetDefaultKeyPath(t *testing.T) {
	t.Parallel()

	manager, err := NewManager(t.TempDir())
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}

	got := manager.GetDefaultKeyPath()
	want := filepath.Join(manager.keysDir, DefaultKeyName)

	if got != want {
		t.Fatalf("GetDefaultKeyPath() = %s, want %s", got, want)
	}
}
