package leveldb

import (
	"log"

	"goqueue/pkg/storage"

	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/net/context"
)

type Config struct {
	Path string
}

func New(ctx context.Context, conf Config) (*Storage, error) {
	db, err := leveldb.OpenFile(conf.Path, nil)
	if err != nil {
		log.Println("init lvl failed", err)
	}
	st := &Storage{
		ctx: ctx,
		db:  db,
	}
	return st, nil
}

type Storage struct {
	ctx context.Context
	db  *leveldb.DB
}

func (s Storage) Close() error {
	return s.db.Close()
}

func (s Storage) Push(topic, sequenceKey, msgID string, payload []byte) error {
	return nil
}

func (s Storage) Get(topic string) (*storage.Message, error) {
	return &storage.Message{}, nil
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
