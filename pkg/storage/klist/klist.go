package klist

import (
	"sync"

	batchapi "goqueue/pkg/storage/batch"

	"golang.org/x/net/context"
)

type DB interface {
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
	Delete(key []byte) error
	Has(key []byte) (bool, error)
	Write(batch batchapi.List) error
	IsNotFoundErr(err error) bool
}

func New(ctx context.Context, name string, db DB) *KList {
	list := KList{
		ctx:  ctx,
		db:   db,
		name: []byte(name),
	}
	return &list
}

type KList struct {
	ctx  context.Context
	db   DB
	name []byte
	rw   sync.RWMutex
}

func (l *KList) Add(item []byte) error {
	l.rw.Lock()
	defer l.rw.Unlock()

	state, err := l.loadState()
	if err != nil {
		return err
	}

	return l.add(state, item)
}

func (l *KList) SetToBegin(item []byte) error {
	l.rw.Lock()
	defer l.rw.Unlock()

	state, err := l.loadState()
	if err != nil {
		return err
	}

	return l.setToBegin(state, item)
}

func (l *KList) SetToEnd(item []byte) error {
	l.rw.Lock()
	defer l.rw.Unlock()

	state, err := l.loadState()
	if err != nil {
		return err
	}

	return l.setToEnd(state, item)
}

func (l *KList) SetAfter(item, root []byte) error {
	l.rw.Lock()
	defer l.rw.Unlock()

	state, err := l.loadState()
	if err != nil {
		return err
	}

	return l.setAfter(state, item, root)
}

func (l *KList) SetBefore(item, root []byte) error {
	l.rw.Lock()
	defer l.rw.Unlock()

	state, err := l.loadState()
	if err != nil {
		return err
	}

	return l.setBefore(state, item, root)
}

func (l *KList) Pop() ([]byte, error) {
	l.rw.Lock()
	defer l.rw.Unlock()

	state, err := l.loadState()
	if err != nil {
		return nil, err
	}

	return l.pop(state)
}

func (l *KList) Delete(item []byte) error {
	l.rw.Lock()
	defer l.rw.Unlock()

	state, err := l.loadState()
	if err != nil {
		return err
	}

	return l.delete(state, item)
}

func (l *KList) GetFirst() ([]byte, error) {
	l.rw.RLock()
	defer l.rw.RUnlock()

	state, err := l.loadState()
	if err != nil {
		return nil, err
	}

	return state.GetFirstItem(), nil
}

func (l *KList) GetLast() ([]byte, error) {
	l.rw.RLock()
	defer l.rw.RUnlock()

	state, err := l.loadState()
	if err != nil {
		return nil, err
	}

	return state.GetLastItem(), nil
}

func (l *KList) GetNext(item []byte) ([]byte, error) {
	l.rw.RLock()
	defer l.rw.RUnlock()

	rec, err := l.readRecord(item)
	if err != nil {
		return nil, err
	}

	return rec.Next, nil
}

func (l *KList) GetPrev(item []byte) ([]byte, error) {
	l.rw.RLock()
	defer l.rw.RUnlock()

	rec, err := l.readRecord(item)
	if err != nil {
		return nil, err
	}

	return rec.Prev, nil
}

func (l *KList) GetCount() (int64, error) {
	l.rw.RLock()
	defer l.rw.RUnlock()

	state, err := l.loadState()
	if err != nil {
		return -1, err
	}

	return state.GetCount(), nil
}

func (l *KList) IsItemExists(item []byte) (bool, error) {
	l.rw.RLock()
	defer l.rw.RUnlock()

	return l.isItemExists(item)
}

func (l *KList) IsItemFirst(item []byte) (bool, error) {
	l.rw.RLock()
	defer l.rw.RUnlock()

	state, err := l.loadState()
	if err != nil {
		return false, err
	}

	return isItemFirst(state, item), nil
}

func (l *KList) IsItemLast(item []byte) (bool, error) {
	l.rw.RLock()
	defer l.rw.RUnlock()

	state, err := l.loadState()
	if err != nil {
		return false, err
	}

	return isItemLast(state, item), nil
}
