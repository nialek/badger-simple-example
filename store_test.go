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

	t.Run("set", func(t *testing.T) {
		assert.NoError(t, store.set(key, value))

		result, err := store.get(key)
		assert.NoError(t, err, "failed to get from store")

		assert.Equal(t, result, value)
	})

	newKey := []byte("second")
	t.Run("rename", func(t *testing.T) {
		assert.NoError(t, store.rename(key, newKey))

		_, err := store.get(key)
		assert.ErrorIs(t, err, badger.ErrKeyNotFound)

		result, err := store.get(newKey)
		assert.NoError(t, err)

		assert.Equal(t, result, value)
	})

	t.Run("delete", func(t *testing.T) {
		assert.NoError(t, store.delete(newKey))

		_, err := store.get(newKey)
		assert.ErrorIs(t, err, badger.ErrKeyNotFound)
	})

	t.Run("namespaces", func(t *testing.T) {
		assert.NoError(t, store.Namespace("A").set(key, value))

		_, err := store.Namespace("B").get(key)
		assert.ErrorIs(t, err, badger.ErrKeyNotFound)
	})

	assert.NoError(t, store.Close(), "failed to close store")
}
