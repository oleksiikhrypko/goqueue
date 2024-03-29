package db

import "goqueue/pkg/storage/batch"

type DB interface {
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
	Delete(key []byte) error
	Has(key []byte) (bool, error)
	Write(batch batch.List) error
	IsNotFoundErr(err error) bool
}
