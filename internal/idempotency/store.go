package idempotency

import (
	"context"
	"sync"
	"time"
)

// Entry holds the cached response for a given idempotency key.
type Entry struct {
	StatusCode  int
	Body        []byte
	ContentType string
	CreatedAt   time.Time
}

// IStore is the storage contract for idempotency entries.
type IStore interface {
	Get(ctx context.Context, key string) (*Entry, error)
	Set(ctx context.Context, key string, entry Entry) error
}

// InMemoryStore is a thread-safe, TTL-aware in-memory implementation of IStore.
// Expiry is checked lazily on Get — no background goroutine needed.
type InMemoryStore struct {
	mu      sync.RWMutex
	entries map[string]Entry
	ttl     time.Duration

	// Programmable hooks for testing error paths.
	OnGet func(ctx context.Context, key string) (*Entry, error)
	OnSet func(ctx context.Context, key string, entry Entry) error
}

func NewInMemoryStore(ttl time.Duration) *InMemoryStore {
	return &InMemoryStore{
		entries: make(map[string]Entry),
		ttl:     ttl,
	}
}

func (s *InMemoryStore) Get(ctx context.Context, key string) (*Entry, error) {
	if s.OnGet != nil {
		return s.OnGet(ctx, key)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.entries[key]
	if !ok || time.Since(entry.CreatedAt) > s.ttl {
		return nil, nil
	}
	return &entry, nil
}

func (s *InMemoryStore) Set(ctx context.Context, key string, entry Entry) error {
	if s.OnSet != nil {
		return s.OnSet(ctx, key, entry)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries[key] = entry
	return nil
}

// Clear resets all entries and hooks — use between tests.
func (s *InMemoryStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = make(map[string]Entry)
	s.OnGet = nil
	s.OnSet = nil
}
