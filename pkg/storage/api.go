package storage

import "goqueue/internal/batch"

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
	Add(actions batch.ActionsList, item []byte) error
	SetToBegin(actions batch.ActionsList, item []byte) error
	SetToEnd(actions batch.ActionsList, item []byte) error
	SetAfter(actions batch.ActionsList, item, root []byte) error
	SetBefore(actions batch.ActionsList, item, root []byte) error
	GetFirst() ([]byte, error)
	GetLast() ([]byte, error)
	GetNext(item []byte) ([]byte, error)
	GetPrev(item []byte) ([]byte, error)
	GetCount() (int64, error)
	Pop(actions batch.ActionsList) ([]byte, error)
	Delete(actions batch.ActionsList, item []byte) error
	IsItemExists(item []byte) (bool, error)
	IsItemFirst(item []byte) (bool, error)
	IsItemLast(item []byte) (bool, error)
}

type DB interface {
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
	Delete(key []byte) error
	Has(key []byte) (bool, error)
	Write(batch batch.ActionsList) error
	IsNotFoundErr(err error) bool
}
