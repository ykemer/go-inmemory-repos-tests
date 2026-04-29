package idempotency

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryStore_Get(t *testing.T) {
	t.Run("MissingKey_ReturnsNil", func(t *testing.T) {
		store := NewInMemoryStore(24 * time.Hour)
		ctx := context.Background()

		entry, err := store.Get(ctx, "nonexistent")

		assert.NoError(t, err)
		assert.Nil(t, entry)
	})

	t.Run("ExistingKey_ReturnsEntry", func(t *testing.T) {
		store := NewInMemoryStore(24 * time.Hour)
		ctx := context.Background()

		store.Set(ctx, "k1", Entry{StatusCode: 201, Body: []byte(`{"id":1}`), CreatedAt: time.Now()})
		entry, err := store.Get(ctx, "k1")

		assert.NoError(t, err)
		assert.NotNil(t, entry)
		assert.Equal(t, 201, entry.StatusCode)
		assert.Equal(t, []byte(`{"id":1}`), entry.Body)
	})

	t.Run("ExpiredEntry_ReturnsNil", func(t *testing.T) {
		store := NewInMemoryStore(1 * time.Millisecond)
		ctx := context.Background()

		store.Set(ctx, "k1", Entry{StatusCode: 201, CreatedAt: time.Now()})
		time.Sleep(5 * time.Millisecond)

		entry, err := store.Get(ctx, "k1")

		assert.NoError(t, err)
		assert.Nil(t, entry)
	})

	t.Run("OnGetHook_PropagatesError", func(t *testing.T) {
		store := NewInMemoryStore(24 * time.Hour)
		store.OnGet = func(_ context.Context, _ string) (*Entry, error) {
			return nil, fmt.Errorf("store unavailable")
		}

		_, err := store.Get(context.Background(), "k1")

		assert.EqualError(t, err, "store unavailable")
	})

	t.Run("OnGetHook_CanReturnEntry", func(t *testing.T) {
		store := NewInMemoryStore(24 * time.Hour)
		fakeEntry := &Entry{StatusCode: 200}
		store.OnGet = func(_ context.Context, _ string) (*Entry, error) {
			return fakeEntry, nil
		}

		entry, err := store.Get(context.Background(), "any-key")

		assert.NoError(t, err)
		assert.Equal(t, fakeEntry, entry)
	})
}

func TestInMemoryStore_Set(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		store := NewInMemoryStore(24 * time.Hour)
		ctx := context.Background()

		err := store.Set(ctx, "k1", Entry{StatusCode: 201, CreatedAt: time.Now()})

		assert.NoError(t, err)
	})

	t.Run("OverwritesExistingKey", func(t *testing.T) {
		store := NewInMemoryStore(24 * time.Hour)
		ctx := context.Background()

		store.Set(ctx, "k1", Entry{StatusCode: 200, CreatedAt: time.Now()})
		store.Set(ctx, "k1", Entry{StatusCode: 201, CreatedAt: time.Now()})
		entry, _ := store.Get(ctx, "k1")

		assert.Equal(t, 201, entry.StatusCode)
	})

	t.Run("OnSetHook_PropagatesError", func(t *testing.T) {
		store := NewInMemoryStore(24 * time.Hour)
		store.OnSet = func(_ context.Context, _ string, _ Entry) error {
			return fmt.Errorf("write failed")
		}

		err := store.Set(context.Background(), "k1", Entry{})

		assert.EqualError(t, err, "write failed")
	})
}

func TestInMemoryStore_Clear(t *testing.T) {
	t.Run("ResetsEntriesAndHooks", func(t *testing.T) {
		store := NewInMemoryStore(24 * time.Hour)
		ctx := context.Background()

		store.Set(ctx, "k1", Entry{StatusCode: 201, CreatedAt: time.Now()})
		store.OnGet = func(_ context.Context, _ string) (*Entry, error) {
			return nil, fmt.Errorf("hook error")
		}
		store.OnSet = func(_ context.Context, _ string, _ Entry) error {
			return fmt.Errorf("hook error")
		}

		store.Clear()

		// Hooks cleared — real path executes again
		entry, err := store.Get(ctx, "k1")
		assert.NoError(t, err)
		assert.Nil(t, entry)
	})
}
