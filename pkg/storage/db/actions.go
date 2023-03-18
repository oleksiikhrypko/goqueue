package db

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

type Reader interface {
	Get(key []byte) ([]byte, error)
	IsNotFoundErr(err error) bool
}

func ReadValue(db Reader, key []byte) ([]byte, error) {
	value, err := db.Get(key)
	if err != nil {
		if db.IsNotFoundErr(err) {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to read data")
	}
	if len(value) == 0 {
		return nil, nil
	}
	return value, nil
}

func ReadStruct(db Reader, key []byte, dest proto.Message) error {
	v, err := ReadValue(db, key)
	if err != nil {
		if err == ErrNotFound {
			return nil
		}
		return errors.Wrap(err, "failed to read state")
	}
	if v == nil {
		return nil
	}

	err = proto.Unmarshal(v, dest)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal state")
	}
	return nil
}
