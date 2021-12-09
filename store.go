package store

import "github.com/dgraph-io/badger/v3"

type Store struct {
	badger *badger.DB
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

func (s Store) Get(key []byte) (value []byte, err error) {
	return value, s.badger.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		value, err = item.ValueCopy(value)
		return err
	})
}

func (s Store) Set(key, value []byte) error {
	return s.badger.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func (s Store) Delete(key []byte) error {
	return s.badger.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func (s Store) Rename(key []byte, newKey []byte) error {
	return s.badger.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		partnerId, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		err = txn.Delete(key)
		if err != nil {
			return err
		}

		return txn.Set(newKey, partnerId)
	})
}

func ConcatKey(namespace string, key []byte) []byte {
	return append([]byte(namespace+":"), key...)
}
