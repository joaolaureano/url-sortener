// Package storage holds the code -> original URL references. It is keyed by
// short code alone, so the link format of package url can change without
// invalidating stored data.
package storage

import (
	"errors"
	"fmt"
	"sync"
)

// ErrNotFound reports that a key holds no reference.
var ErrNotFound = errors.New("storage: reference not found")

// Store persists short code references.
//
// Implementations must be safe for concurrent use and must never let a create
// silently replace an existing, different value.
type Store interface {
	// CreateReference binds key to value.
	//
	// created is true when key ends up bound to value, either by a fresh
	// write or by an idempotent repeat of the same value. It is false when
	// key already holds a different value — a collision — in which case
	// current carries the value that was kept.
	CreateReference(key string, value string) (created bool, current string)

	// RecoverReference returns the value bound to key, or ErrNotFound.
	RecoverReference(key string) (string, error)
}

// MemoryStore is an in-process Store backed by a map. The zero value is not
// usable; call NewMemoryStore.
type MemoryStore struct {
	mu    sync.RWMutex
	table map[string]string
}

// NewMemoryStore returns an empty MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{table: make(map[string]string)}
}

// Default is the Store backing the package-level functions. Replace it to swap
// in a durable backend.
var Default Store = NewMemoryStore()

// CreateReference implements Store.
func (s *MemoryStore) CreateReference(key string, value string) (created bool, current string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if existing, ok := s.table[key]; ok {
		return existing == value, existing
	}

	s.table[key] = value

	return true, value
}

// RecoverReference implements Store.
func (s *MemoryStore) RecoverReference(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.table[key]
	if !ok {
		return "", fmt.Errorf("%w: %q", ErrNotFound, key)
	}

	return value, nil
}

// CreateReference binds key to value in Default. See Store.CreateReference.
func CreateReference(key string, value string) (created bool, current string) {
	return Default.CreateReference(key, value)
}

// RecoverReference returns the value bound to key in Default, or ErrNotFound.
func RecoverReference(key string) (string, error) {
	return Default.RecoverReference(key)
}
