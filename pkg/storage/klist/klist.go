package klist

import (
	"sync"

	models "goqueue/pkg/proto/klist"
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

func New(ctx context.Context, name string, db DB) (*KList, error) {
	list := KList{
		ctx:   ctx,
		db:    db,
		name:  []byte(name),
		state: &models.State{},
	}
	if err := list.loadState(); err != nil {
		return nil, err
	}
	return &list, nil
}

type KList struct {
	ctx  context.Context
	db   DB
	name []byte
	rw   sync.RWMutex

	state *models.State
}

func (l *KList) Add(item []byte) (err error) {
	l.rw.Lock()
	defer l.rw.Unlock()
	defer func() {
		if err != nil {
			l.mustLoadState()
		}
	}()

	err = l.add(item)

	return
}

func (l *KList) SetToBegin(item []byte) (err error) {
	l.rw.Lock()
	defer l.rw.Unlock()
	defer func() {
		if err != nil {
			l.mustLoadState()
		}
	}()

	err = l.setToBegin(item)

	return
}

func (l *KList) SetToEnd(item []byte) (err error) {
	l.rw.Lock()
	defer l.rw.Unlock()
	defer func() {
		if err != nil {
			l.mustLoadState()
		}
	}()

	err = l.setToEnd(item)

	return
}

func (l *KList) SetAfter(item, root []byte) (err error) {
	l.rw.Lock()
	defer l.rw.Unlock()
	defer func() {
		if err != nil {
			l.mustLoadState()
		}
	}()

	err = l.setAfter(item, root)

	return
}

func (l *KList) SetBefore(item, root []byte) (err error) {
	l.rw.Lock()
	defer l.rw.Unlock()
	defer func() {
		if err != nil {
			l.mustLoadState()
		}
	}()

	err = l.setBefore(item, root)

	return
}

func (l *KList) Pop() (item []byte, err error) {
	l.rw.Lock()
	defer l.rw.Unlock()
	defer func() {
		if err != nil {
			l.mustLoadState()
		}
	}()

	item, err = l.pop()

	return
}

func (l *KList) Delete(item []byte) (err error) {
	l.rw.Lock()
	defer l.rw.Unlock()
	defer func() {
		if err != nil {
			l.mustLoadState()
		}
	}()

	err = l.delete(item)

	return
}

func (l *KList) GetFirst() []byte {
	l.rw.RLock()
	defer l.rw.RUnlock()

	return l.state.GetFirstItem()
}

func (l *KList) GetLast() []byte {
	l.rw.RLock()
	defer l.rw.RUnlock()

	return l.state.GetLastItem()
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

func (l *KList) GetCount() int64 {
	l.rw.RLock()
	defer l.rw.RUnlock()

	return l.state.Count
}

func (l *KList) IsItemExists(item []byte) (bool, error) {
	l.rw.RLock()
	defer l.rw.RUnlock()

	return l.isItemExists(item)
}

func (l *KList) IsItemFirst(item []byte) bool {
	l.rw.RLock()
	defer l.rw.RUnlock()

	return l.isItemFirst(item)
}

func (l *KList) IsItemLast(item []byte) bool {
	l.rw.RLock()
	defer l.rw.RUnlock()

	return l.isItemLast(item)
}
