package klist

import (
	"goqueue/pkg/storage/batch"

	"golang.org/x/net/context"
)

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
}

func (l *KList) Add(actions batch.List, item []byte) error {
	state, err := l.loadState()
	if err != nil {
		return err
	}

	return l.add(actions, state, item)
}

func (l *KList) SetToBegin(actions batch.List, item []byte) error {
	state, err := l.loadState()
	if err != nil {
		return err
	}

	return l.setToBegin(actions, state, item)
}

func (l *KList) SetToEnd(actions batch.List, item []byte) error {
	state, err := l.loadState()
	if err != nil {
		return err
	}

	return l.setToEnd(actions, state, item)
}

func (l *KList) SetAfter(actions batch.List, item, root []byte) error {
	state, err := l.loadState()
	if err != nil {
		return err
	}

	return l.setAfter(actions, state, item, root)
}

func (l *KList) SetBefore(actions batch.List, item, root []byte) error {
	state, err := l.loadState()
	if err != nil {
		return err
	}

	return l.setBefore(actions, state, item, root)
}

func (l *KList) Pop(actions batch.List) ([]byte, error) {
	state, err := l.loadState()
	if err != nil {
		return nil, err
	}

	return l.pop(actions, state)
}

func (l *KList) Delete(actions batch.List, item []byte) error {
	state, err := l.loadState()
	if err != nil {
		return err
	}

	return l.delete(actions, state, item)
}

func (l *KList) GetFirst() ([]byte, error) {
	state, err := l.loadState()
	if err != nil {
		return nil, err
	}

	return state.GetFirstItem(), nil
}

func (l *KList) GetLast() ([]byte, error) {
	state, err := l.loadState()
	if err != nil {
		return nil, err
	}

	return state.GetLastItem(), nil
}

func (l *KList) GetNext(item []byte) ([]byte, error) {
	rec, err := l.readRecord(item)
	if err != nil {
		return nil, err
	}

	return rec.Next, nil
}

func (l *KList) GetPrev(item []byte) ([]byte, error) {
	rec, err := l.readRecord(item)
	if err != nil {
		return nil, err
	}

	return rec.Prev, nil
}

func (l *KList) GetCount() (int64, error) {
	state, err := l.loadState()
	if err != nil {
		return -1, err
	}

	return state.GetCount(), nil
}

func (l *KList) IsItemExists(item []byte) (bool, error) {
	return l.isItemExists(item)
}

func (l *KList) IsItemFirst(item []byte) (bool, error) {
	state, err := l.loadState()
	if err != nil {
		return false, err
	}

	return isItemFirst(state, item), nil
}

func (l *KList) IsItemLast(item []byte) (bool, error) {
	state, err := l.loadState()
	if err != nil {
		return false, err
	}

	return isItemLast(state, item), nil
}
