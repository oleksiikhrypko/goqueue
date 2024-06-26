package inmem

import (
	"sync"

	batchapi "goqueue/internal/batch"

	"github.com/pkg/errors"
)

var (
	ErrNotFound = errors.New("inmemdb: not found")
)

func NewDB() *DB {
	return &DB{
		data: make(map[string][]byte),
		rw:   sync.RWMutex{},
	}
}

type DB struct {
	data map[string][]byte
	rw   sync.RWMutex
}

func (db *DB) Get(key []byte) ([]byte, error) {
	db.rw.RLock()
	defer db.rw.RUnlock()
	v, ok := db.data[string(key)]
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}

func (db *DB) Put(key, value []byte) error {
	db.rw.Lock()
	defer db.rw.Unlock()
	db.data[string(key)] = value
	return nil
}

func (db *DB) Delete(key []byte) error {
	db.rw.Lock()
	defer db.rw.Unlock()
	delete(db.data, string(key))
	return nil
}

func (db *DB) Has(key []byte) (bool, error) {
	db.rw.RLock()
	defer db.rw.RUnlock()
	_, ok := db.data[string(key)]
	return ok, nil
}

func (db *DB) Write(batch batchapi.ActionsList) error {
	db.rw.Lock()
	defer db.rw.Unlock()
	batch.ForEach(func(action batchapi.ActionType, key, value []byte) {
		db.data[string(key)] = value
	})
	return nil
}

func (db *DB) IsNotFoundErr(err error) bool {
	return errors.Is(err, ErrNotFound)
}
