package db

import (
	"errors"

	"github.com/golang/protobuf/proto"
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
		return nil, ErrCritical.Consume(err)
	}
	if len(value) == 0 {
		return nil, nil
	}
	return value, nil
}

func ReadStruct(db Reader, key []byte, dest proto.Message) error {
	v, err := ReadValue(db, key)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil
		}
		return ErrCritical.Consume(err)
	}
	if v == nil {
		return nil
	}

	err = proto.Unmarshal(v, dest)
	if err != nil {
		return ErrCritical.Consume(err)
	}
	return nil
}
