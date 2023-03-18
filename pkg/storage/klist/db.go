package klist

import (
	"goqueue/pkg/storage/db"

	"github.com/pkg/errors"
)

type DB interface {
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	IsNotFoundErr(err error) bool
}

func (l *KList) readValue(key []byte) ([]byte, error) {
	return db.ReadValue(l.db, key)
}

func (l *KList) hasKey(key []byte) (bool, error) {
	has, err := l.db.Has(key)
	if err != nil {
		return false, errors.Wrap(err, "failed to check if key exists")
	}
	return has, nil
}

func (l *KList) isItemExists(item []byte) (bool, error) {
	return l.hasKey(buildItemKey(l.name, item))
}
