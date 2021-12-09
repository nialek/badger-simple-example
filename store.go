package store

import (
	"github.com/dgraph-io/badger/v3"
)

type Store struct {
	badger    *badger.DB
	namespace string
}

func (s Store) Namespace(namespace string) Store {
	s.namespace = namespace
	return s
}

func NewStore(inMemory bool) (s Store, err error) {
	options := badger.DefaultOptions("/var/lib/badger/db").WithInMemory(false)
	if inMemory {
		options = badger.DefaultOptions("").WithInMemory(true)
	}

	s.badger, err = badger.Open(options)
	return s, err
}

func (s Store) Close() error {
	return s.badger.Close()
}

func (s Store) get(key []byte) (value []byte, err error) {
	return value, s.badger.View(func(txn *badger.Txn) error {
		item, err := txn.Get(s.concat(key))
		if err != nil {
			return err
		}

		value, err = item.ValueCopy(nil)
		return err
	})
}

func (s Store) set(key, value []byte) error {
	return s.badger.Update(func(txn *badger.Txn) error {
		return txn.Set(s.concat(key), value)
	})
}

func (s Store) delete(key []byte) error {
	return s.badger.Update(func(txn *badger.Txn) error {
		return txn.Delete(s.concat(key))
	})
}

func (s Store) rename(key []byte, newKey []byte) error {
	return s.badger.Update(func(txn *badger.Txn) error {
		oldKey := s.concat(key)
		item, err := txn.Get(oldKey)
		if err != nil {
			return err
		}

		partnerId, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		err = txn.Delete(oldKey)
		if err != nil {
			return err
		}

		return txn.Set(s.concat(newKey), partnerId)
	})
}

func (s Store) concat(key []byte) []byte {
	namespace := []byte(s.namespace + ":")
	return append(namespace, key...)
}
