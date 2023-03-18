package group

import "goqueue/pkg/storage/batch"

type DB interface {
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	IsNotFoundErr(err error) bool
	Write(batch batch.List) error
}

type List interface {
	GetFirst() ([]byte, error)
	GetNext(item []byte) ([]byte, error)
	GetCount() (int64, error)
	IsItemLast(item []byte) (bool, error)
}

type Group struct {
	db   DB
	name []byte
	src  List
}

func New(db DB, name []byte, src List) *Group {
	return &Group{
		db:   db,
		name: name,
		src:  src,
	}
}
