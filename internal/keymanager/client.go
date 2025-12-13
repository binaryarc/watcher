package keymanager

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
)

const (
	KeyPrefix      = "watcher_"
	KeyLength      = 32
	DefaultKeyName = "default"
)

// Manager handles client-side API key management
type Manager struct {
	keysDir string
}

// NewManager creates a new key manager
func NewManager(keysDir string) (*Manager, error) {
	if err := os.MkdirAll(keysDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create keys directory: %w", err)
	}

	return &Manager{
		keysDir: keysDir,
	}, nil
}

// Generate creates a new API key
func (m *Manager) Generate() (string, error) {
	randomBytes := make([]byte, KeyLength)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	encoded := base64.URLEncoding.EncodeToString(randomBytes)
	apiKey := KeyPrefix + encoded

	return apiKey, nil
}

// Save stores an API key with the given name
func (m *Manager) Save(name, apiKey string) error {
	keyPath := filepath.Join(m.keysDir, name)

	if err := os.WriteFile(keyPath, []byte(apiKey), 0600); err != nil {
		return fmt.Errorf("failed to save key: %w", err)
	}

	return nil
}

// Load reads an API key by name
func (m *Manager) Load(name string) (string, error) {
	keyPath := filepath.Join(m.keysDir, name)

	data, err := os.ReadFile(keyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("key %q not found", name)
		}
		return "", fmt.Errorf("failed to load key: %w", err)
	}

	return string(data), nil
}

// List returns all saved key names
func (m *Manager) List() ([]string, error) {
	entries, err := os.ReadDir(m.keysDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to list keys: %w", err)
	}

	var keys []string
	for _, entry := range entries {
		if !entry.IsDir() {
			keys = append(keys, entry.Name())
		}
	}

	return keys, nil
}

// Delete removes a saved key
func (m *Manager) Delete(name string) error {
	keyPath := filepath.Join(m.keysDir, name)

	if err := os.Remove(keyPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("key %q not found", name)
		}
		return fmt.Errorf("failed to delete key: %w", err)
	}

	return nil
}

// GetDefaultKeyPath returns the path to the default key file
func (m *Manager) GetDefaultKeyPath() string {
	return filepath.Join(m.keysDir, DefaultKeyName)
}
