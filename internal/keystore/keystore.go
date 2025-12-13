package keystore

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// KeyInfo stores metadata about an API key
type KeyInfo struct {
	Key         string    `json:"key"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Store manages API keys on the server side
type Store struct {
	mu       sync.RWMutex
	keys     map[string]*KeyInfo
	filePath string
}

// NewStore creates a new key store
func NewStore(filePath string) (*Store, error) {
	store := &Store{
		keys:     make(map[string]*KeyInfo),
		filePath: filePath,
	}

	if filePath != "" {
		if err := store.load(); err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load keys: %w", err)
		}
	}

	return store, nil
}

// Add adds a new API key to the store
func (s *Store) Add(key, description string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.keys[key]; exists {
		return fmt.Errorf("key already exists")
	}

	s.keys[key] = &KeyInfo{
		Key:         key,
		Description: description,
		CreatedAt:   time.Now(),
	}

	return s.save()
}

// Remove removes an API key from the store
func (s *Store) Remove(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.keys[key]; !exists {
		return fmt.Errorf("key not found")
	}

	delete(s.keys, key)

	return s.save()
}

// Validate checks if an API key is valid using constant-time comparison
func (s *Store) Validate(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 키가 없으면 모든 요청 거부
	if len(s.keys) == 0 {
		return false // 변경: true → false
	}

	for storedKey := range s.keys {
		if subtle.ConstantTimeCompare([]byte(key), []byte(storedKey)) == 1 {
			return true
		}
	}

	return false
}

// List returns all stored keys
func (s *Store) List() []*KeyInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*KeyInfo, 0, len(s.keys))
	for _, info := range s.keys {
		result = append(result, info)
	}

	return result
}

// Clear removes all keys
func (s *Store) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.keys = make(map[string]*KeyInfo)

	return s.save()
}

// IsEmpty returns true if no keys are stored
func (s *Store) IsEmpty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.keys) == 0
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	var keys []*KeyInfo
	if err := json.Unmarshal(data, &keys); err != nil {
		return fmt.Errorf("failed to parse keys file: %w", err)
	}

	for _, info := range keys {
		s.keys[info.Key] = info
	}

	return nil
}

func (s *Store) save() error {
	if s.filePath == "" {
		return nil
	}

	keys := make([]*KeyInfo, 0, len(s.keys))
	for _, info := range s.keys {
		keys = append(keys, info)
	}

	data, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode keys: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write keys file: %w", err)
	}

	return nil
}
