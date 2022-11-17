package storage

import "goqueue/pkg/storage/batch"

type Message struct {
	ID      string
	Payload []byte
}

type Producer interface {
	Push(topic, sequenceKey, msgID string, payload []byte) error
}

type Consumer interface {
	Get(topic, group string) (*Message, error)
	Commit(msgID string) error
	Rollback(msgID string) error
	Touch(msgID string, visibilityTimeoutSec int) error
}

type List interface {
	Add(item []byte) error
	SetToBegin(item []byte) error
	SetToEnd(item []byte) error
	SetAfter(item, root []byte) error
	SetBefore(item, root []byte) error
	GetFirst() []byte
	GetLast() []byte
	GetNext(item []byte) ([]byte, error)
	GetPrev(item []byte) ([]byte, error)
	GetCount() int64
	Pop() ([]byte, error)
	Delete(item []byte) error
	IsItemExists(item []byte) (bool, error)
	IsItemFirst(item []byte) bool
	IsItemLast(item []byte) bool
}

type DB interface {
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
	Delete(key []byte) error
	Has(key []byte) (bool, error)
	Write(batch batch.List) error
	IsNotFoundErr(err error) bool
}
