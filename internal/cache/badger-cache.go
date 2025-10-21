package cache

import "github.com/dgraph-io/badger/v4"

func NewBadgerDB() (*badger.DB, error) {

	options := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(options)

	return db, err
}
