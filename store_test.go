package store

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStore(t *testing.T) {
	const inMemory = true
	store, err := NewStore(inMemory)
	assert.NoError(t, err, "failed to create store")

	key, value := []byte("first"), []byte("first_value")

	t.Run("Set", func(t *testing.T) {
		assert.NoError(t, store.Set(key, value))

		result, err := store.Get(key)
		assert.NoError(t, err, "failed to Get from store")

		assert.Equal(t, result, value)
	})

	newKey := []byte("second")
	t.Run("Rename", func(t *testing.T) {
		assert.NoError(t, store.Rename(key, newKey))

		_, err := store.Get(key)
		assert.ErrorIs(t, err, badger.ErrKeyNotFound)

		result, err := store.Get(newKey)
		assert.NoError(t, err)

		assert.Equal(t, result, value)
	})

	t.Run("Delete", func(t *testing.T) {
		assert.NoError(t, store.Delete(newKey))

		_, err := store.Get(newKey)
		assert.ErrorIs(t, err, badger.ErrKeyNotFound)
	})

	assert.NoError(t, store.Close(), "failed to close store")
}
