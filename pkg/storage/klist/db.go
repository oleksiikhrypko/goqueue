package klist

import (
	"github.com/pkg/errors"
)

type DB interface {
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	IsNotFoundErr(err error) bool
}

func (l *KList) readValue(key []byte) ([]byte, error) {
	value, err := l.db.Get(key)
	if err != nil {
		if l.db.IsNotFoundErr(err) {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to read data")
	}
	if len(value) == 0 {
		return nil, nil
	}
	return value, nil
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
