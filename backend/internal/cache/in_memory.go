package cache

import (
	"errors"
	"sync"
	"time"
)

type CacheItem struct {
	Value      string
	Expiration time.Time
}

type InMemoryStore struct {
	data map[string]CacheItem
	mu   sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string]CacheItem),
	}
}

func (s *InMemoryStore) Set(key string, value string, duration time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = CacheItem{
		Value:      value,
		Expiration: time.Now().Add(duration),
	}
	return nil
}

func (s *InMemoryStore) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.data[key]
	if !exists {
		return "", errors.New("key not found")
	}

	if time.Now().After(item.Expiration) {
		delete(s.data, key)
		return "", errors.New("key expired")
	}

	return item.Value, nil
}
