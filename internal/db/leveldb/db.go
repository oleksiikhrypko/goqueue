package leveldb

import (
	"goqueue/internal/batch"

	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/net/context"
)

func NewDB(ctx context.Context, lvl *leveldb.DB) *DB {
	return &DB{
		ctx: ctx,
		lvl: lvl,
	}
}

type DB struct {
	ctx context.Context
	lvl *leveldb.DB
}

func (db *DB) Get(key []byte) ([]byte, error) {
	return db.lvl.Get(key, nil)
}

func (db *DB) Put(key, value []byte) error {
	return db.lvl.Put(key, value, nil)
}

func (db *DB) Delete(key []byte) error {
	return db.lvl.Delete(key, nil)
}

func (db *DB) Has(key []byte) (bool, error) {
	return db.lvl.Has(key, nil)
}

func (db *DB) Write(batch batch.List) error {
	return db.lvl.Write(fromAPIBatch(batch), nil)
}

func (db *DB) IsNotFoundErr(err error) bool {
	return errors.Is(err, leveldb.ErrNotFound)
}

func fromAPIBatch(in batch.List) *leveldb.Batch {
	out := leveldb.Batch{}
	in.ForEach(func(action batch.ActionType, key, value []byte) {
		switch action {
		case batch.ActionTypePut:
			out.Put(key, value)
		case batch.ActionTypeDel:
			out.Delete(key)
		}
	})
	return &out
}
