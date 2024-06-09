package group

import "goqueue/internal/db"

type DB interface {
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	IsNotFoundErr(err error) bool
}

type Group struct {
	db    DB
	topic string
	name  string
}

func New(db DB, topic, name string) *Group {
	return &Group{
		db:    db,
		topic: topic,
		name:  name,
	}
}

func Exists(db db.DB, topic, name string) (bool, error) {
	return db.Has(buildStateKey(topic, name))
}
