package storage

import (
	"context"
	"sync"
)

// MemoryStore is an in-memory implementation for development/testing.
type MemoryStore struct {
	mu    sync.RWMutex
	notes map[int64][]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		notes: make(map[int64][]string),
	}
}

func (s *MemoryStore) SaveNote(_ context.Context, userID int64, text string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notes[userID] = append(s.notes[userID], text)
	return nil
}

func (s *MemoryStore) GetNotes(_ context.Context, userID int64) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.notes[userID], nil
}

func (s *MemoryStore) CountNotes(_ context.Context, userID int64) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.notes[userID]), nil
}

func (s *MemoryStore) Close() error {
	return nil
}
