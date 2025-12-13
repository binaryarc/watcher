package common

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/binaryarc/watcher/internal/keystore"
)

func KeyStore() (*keystore.Store, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	keysDir := filepath.Join(homeDir, ".watcher", "server")
	if err := os.MkdirAll(keysDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create keys directory: %w", err)
	}

	keystorePath := filepath.Join(keysDir, "keys.json")
	return keystore.NewStore(keystorePath)
}

func MaskKey(key string) string {
	if len(key) <= 14 {
		return key
	}
	return key[:10] + "..." + key[len(key)-4:]
}
