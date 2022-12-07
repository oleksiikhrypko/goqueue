package klist

import (
	"goqueue/pkg/storage/batch"

	"github.com/pkg/errors"
)

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

func (l *KList) writeValue(key, value []byte) error {
	err := l.db.Put(key, value)
	if err != nil {
		return errors.Wrap(err, "failed to write data")
	}
	return nil
}

func (l *KList) writeBatch(actions *batch.Batch) error {
	err := l.db.Write(actions)
	if err != nil {
		return errors.Wrap(err, "failed to write batch")
	}
	return nil
}

func (l *KList) deleteValue(key []byte) error {
	err := l.db.Delete(key)
	if err != nil {
		return errors.Wrap(err, "failed to delete data")
	}
	return nil
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
