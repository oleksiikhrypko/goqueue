package klist

import (
	"fmt"

	models "goqueue/pkg/proto/models"
	"goqueue/pkg/storage/batch"
	"goqueue/pkg/storage/db"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

func (l *KList) readRecord(item []byte) (*Record, error) {
	if item == nil {
		return nil, nil
	}
	var data models.Item
	err := db.ReadStruct(l.db, buildItemKey(l.name, item), &data)
	if err != nil {
		return nil, err
	}

	return &Record{
		Id:   item,
		Next: data.GetNext(),
		Prev: data.GetPrev(),
	}, nil
}

func (l *KList) loadState() (*models.KList, error) {
	var state models.KList
	err := db.ReadStruct(l.db, buildStateKey(l.name), &state)
	if err != nil {
		return nil, err
	}

	return &state, nil
}

func (l *KList) writeRootItem(actions batch.List, state *models.KList, item []byte) error {
	if !isEmpty(state) {
		return errors.New("failed on call 'writeRootItem': list is not empty")
	}
	var err error

	// update state
	state.FirstItem = item
	state.LastItem = item
	state.Count = 1
	if err = l.saveState(actions, state); err != nil {
		return err
	}
	// save item
	rec := Record{
		Id:   item,
		Next: nil,
		Prev: nil,
	}
	if err = l.saveRecord(actions, &rec); err != nil {
		return err
	}
	return nil
}

func (l *KList) add(actions batch.List, state *models.KList, item []byte) (err error) {
	if isEmpty(state) {
		return l.writeRootItem(actions, state, item)
	}

	// if list already has item -> skip adding
	exists, err := l.isItemExists(item)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	last, err := l.readRecord(state.LastItem)
	if err != nil {
		return err
	}

	// update state
	state.LastItem = item
	state.Count += 1
	if err = l.saveState(actions, state); err != nil {
		return err
	}

	// update records
	rec := Record{
		Id: item,
	}
	if err = l.insertAfterRecord(actions, &rec, last); err != nil {
		return err
	}

	return nil
}

func (l *KList) setToBegin(actions batch.List, state *models.KList, item []byte) (err error) {
	if isEmpty(state) {
		return l.writeRootItem(actions, state, item)
	}

	if isItemFirst(state, item) {
		return nil
	}

	// init record
	rec := &Record{
		Id: item,
	}
	exists, err := l.isItemExists(item)
	if err != nil {
		return err
	}

	first, err := l.readRecord(state.FirstItem)
	if err != nil {
		return err
	}

	// cut item
	if exists {
		rec, err = l.readRecord(item)
		if err != nil {
			return err
		}
		if err = l.cutRecord(actions, rec); err != nil {
			return err
		}
		if isItemFirst(state, item) {
			state.FirstItem = rec.Next
		}
		if isItemLast(state, item) {
			state.LastItem = rec.Prev
		}
	}

	// update state
	if !exists {
		state.Count += 1
	}
	state.FirstItem = item
	if err = l.saveState(actions, state); err != nil {
		return err
	}

	// insert before first
	if err = l.insertBeforeRecord(actions, rec, first); err != nil {
		return err
	}

	return nil
}

func (l *KList) setToEnd(actions batch.List, state *models.KList, item []byte) error {
	if isEmpty(state) {
		return l.writeRootItem(actions, state, item)
	}

	if isItemLast(state, item) {
		return nil
	}

	// init record
	rec := &Record{
		Id: item,
	}
	exists, err := l.isItemExists(item)
	if err != nil {
		return err
	}

	last, err := l.readRecord(state.LastItem)
	if err != nil {
		return err
	}

	// cut item
	if exists {
		rec, err = l.readRecord(item)
		if err != nil {
			return err
		}
		if err = l.cutRecord(actions, rec); err != nil {
			return err
		}
		if isItemFirst(state, item) {
			state.FirstItem = rec.Next
		}
		if isItemLast(state, item) {
			state.LastItem = rec.Prev
		}
	}

	// update state
	if !exists {
		state.Count += 1
	}
	state.LastItem = item
	if err = l.saveState(actions, state); err != nil {
		return err
	}

	// insert after last
	if err = l.insertAfterRecord(actions, rec, last); err != nil {
		return err
	}

	return nil
}

func (l *KList) setAfter(actions batch.List, state *models.KList, item, root []byte) error {
	if isEqual(item, root) {
		return errors.New("input items are equals")
	}

	// init root record
	rootExists, err := l.isItemExists(root)
	if err != nil {
		return err
	}
	if !rootExists {
		return fmt.Errorf("root item '%s' is not exists", string(root))
	}
	rootRec, err := l.readRecord(root)
	if err != nil {
		return err
	}
	if isEqual(item, rootRec.Next) {
		return nil
	}

	// init record
	rec := &Record{
		Id: item,
	}
	exists, err := l.isItemExists(item)
	if err != nil {
		return err
	}

	// cut item
	if exists {
		rec, err = l.readRecord(item)
		if err != nil {
			return err
		}
		if err = l.cutRecord(actions, rec); err != nil {
			return err
		}
		if isItemFirst(state, item) {
			state.FirstItem = rec.Next
		}
		if isItemLast(state, item) {
			state.LastItem = rec.Prev
		}
	}

	// update state
	if !exists {
		state.Count += 1
	}
	if isItemLast(state, root) {
		state.LastItem = item
	}
	if err = l.saveState(actions, state); err != nil {
		return err
	}

	// insert after root
	if err = l.insertAfterRecord(actions, rec, rootRec); err != nil {
		return err
	}

	return nil
}

func (l *KList) setBefore(actions batch.List, state *models.KList, item, root []byte) error {
	if isEqual(item, root) {
		return errors.New("input items are equals")
	}

	// init root record
	rootExists, err := l.isItemExists(root)
	if err != nil {
		return err
	}
	if !rootExists {
		return fmt.Errorf("root item '%s' is not exists", string(root))
	}
	rootRec, err := l.readRecord(root)
	if err != nil {
		return err
	}
	if isEqual(item, rootRec.Prev) {
		return nil
	}

	// init record
	rec := &Record{
		Id: item,
	}
	exists, err := l.isItemExists(item)
	if err != nil {
		return err
	}

	// cut item
	if exists {
		rec, err = l.readRecord(item)
		if err != nil {
			return err
		}
		if err = l.cutRecord(actions, rec); err != nil {
			return err
		}
		if isItemFirst(state, item) {
			state.FirstItem = rec.Next
		}
		if isItemLast(state, item) {
			state.LastItem = rec.Prev
		}
	}
	// update state
	if !exists {
		state.Count += 1
	}
	if isItemFirst(state, root) {
		state.FirstItem = item
	}
	if err = l.saveState(actions, state); err != nil {
		return err
	}

	// insert after root
	if err = l.insertBeforeRecord(actions, rec, rootRec); err != nil {
		return err
	}

	return nil
}

func (l *KList) pop(actions batch.List, state *models.KList) ([]byte, error) {
	if isEmpty(state) {
		return nil, nil
	}

	first, err := l.readRecord(state.FirstItem)
	if err != nil {
		return nil, err
	}

	// update state
	state.FirstItem = first.Next
	if isItemLast(state, first.Id) {
		state.LastItem = nil
	}
	state.Count -= 1
	if err = l.saveState(actions, state); err != nil {
		return nil, err
	}

	// cut from queue
	if err = l.cutRecord(actions, first); err != nil {
		return nil, err
	}

	// delete
	actions.Delete(buildItemKey(l.name, first.Id))

	return first.Id, nil
}

func (l *KList) delete(actions batch.List, state *models.KList, item []byte) error {
	exists, err := l.isItemExists(item)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	rec, err := l.readRecord(item)
	if err != nil {
		return err
	}

	// update state
	if isItemFirst(state, item) {
		state.FirstItem = rec.Next
	}
	if isItemLast(state, item) {
		state.LastItem = rec.Prev
	}
	state.Count -= 1
	if err = l.saveState(actions, state); err != nil {
		return err
	}

	// cut from queue
	if err = l.cutRecord(actions, rec); err != nil {
		return err
	}

	// delete
	actions.Delete(buildItemKey(l.name, rec.Id))

	return nil
}

func (l *KList) saveState(actions batch.List, state *models.KList) error {
	v, err := proto.Marshal(state)
	if err != nil {
		return errors.Wrap(err, "failed to build state data model")
	}
	actions.Put(buildStateKey(l.name), v)
	return nil
}

func (l *KList) saveRecord(actions batch.List, rec *Record) error {
	v, err := proto.Marshal(&models.Item{
		Next: rec.Next,
		Prev: rec.Prev,
	})
	if err != nil {
		return errors.Wrap(err, "failed to build state data model")
	}
	actions.Put(buildItemKey(l.name, rec.Id), v)
	return nil
}

func (l *KList) cutRecord(actions batch.List, rec *Record) error {
	// update prev if exists
	if rec.Prev != nil {
		prev, err := l.readRecord(rec.Prev)
		if err != nil {
			return err
		}
		prev.Next = rec.Next
		if err = l.saveRecord(actions, prev); err != nil {
			return err
		}
	}

	// update next if exists
	if rec.Next != nil {
		next, err := l.readRecord(rec.Next)
		if err != nil {
			return err
		}
		next.Prev = rec.Prev
		if err = l.saveRecord(actions, next); err != nil {
			return err
		}
	}
	return nil
}

func (l *KList) insertAfterRecord(actions batch.List, rec *Record, root *Record) error {
	// in case if it's sequence rec <- root
	if isEqual(rec.Id, root.Prev) {
		root.Prev = rec.Prev
		if root.Prev != nil {
			rprev, err := l.readRecord(root.Prev)
			if err != nil {
				return err
			}
			rprev.Next = root.Id
			if err = l.saveRecord(actions, rprev); err != nil {
				return err
			}
		}
	}

	// update rec
	rec.Prev = root.Id
	rec.Next = root.Next
	if err := l.saveRecord(actions, rec); err != nil {
		return err
	}

	// update next if exists
	if root.Next != nil {
		next, err := l.readRecord(root.Next)
		if err != nil {
			return err
		}
		next.Prev = rec.Id
		if err = l.saveRecord(actions, next); err != nil {
			return err
		}
	}

	root.Next = rec.Id
	if err := l.saveRecord(actions, root); err != nil {
		return err
	}

	return nil
}

func (l *KList) insertBeforeRecord(actions batch.List, rec *Record, root *Record) error {
	// in case if it's sequence root -> rec
	if isEqual(rec.Id, root.Next) {
		root.Next = rec.Next
		if root.Next != nil {
			rnext, err := l.readRecord(root.Next)
			if err != nil {
				return err
			}
			rnext.Prev = root.Id
			if err = l.saveRecord(actions, rnext); err != nil {
				return err
			}
		}
	}

	// update rec
	rec.Prev = root.Prev
	rec.Next = root.Id
	if err := l.saveRecord(actions, rec); err != nil {
		return err
	}

	// update prev if exists
	if root.Prev != nil {
		prev, err := l.readRecord(root.Prev)
		if err != nil {
			return err
		}
		prev.Next = rec.Id
		if err = l.saveRecord(actions, prev); err != nil {
			return err
		}
	}

	// update root
	root.Prev = rec.Id
	if err := l.saveRecord(actions, root); err != nil {
		return err
	}

	return nil
}
