package storage

import (
	"golang.org/x/net/context"
)

type Config struct {
	Path string
}

func New(ctx context.Context, conf Config, db DB) (*Storage, error) {
	st := &Storage{
		ctx: ctx,
		db:  db,
	}
	return st, nil
}

type Storage struct {
	ctx context.Context
	db  DB
}

func (s Storage) Push(topic, sequenceKey, msgID string, payload []byte) error {
	return nil
}

func (s Storage) Get(topic string) (*Message, error) {
	return &Message{}, nil
}

func (s Storage) Commit(msgID string) error {
	return nil
}

func (s Storage) Rollback(msgID string) error {
	return nil

}

func (s Storage) Touch(msgID string, visibilityTimeoutSec int) error {
	return nil

}
