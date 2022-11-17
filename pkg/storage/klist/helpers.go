package klist

import (
	"fmt"

	models "goqueue/pkg/proto/klist"
	batchapi "goqueue/pkg/storage/batch"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

func (l *KList) writeRootItem(item []byte) (err error) {
	if !l.isEmpty() {
		return errors.New("failed on call 'writeRootItem': list is not empty")
	}
	// update state
	batch := batchapi.New(2)
	l.state.FirstItem = item
	l.state.LastItem = item
	l.state.Count = 1
	if err = l.appendBatchSaveState(batch); err != nil {
		return err
	}
	// save item
	rec := Record{
		Id:   item,
		Next: nil,
		Prev: nil,
	}
	if err = l.appendBatchSaveRecord(batch, &rec); err != nil {
		return err
	}
	return l.writeBatch(batch)
}

func (l *KList) readRecord(item []byte) (*Record, error) {
	v, err := l.readValue(buildItemKey(l.name, item))
	if err != nil {
		return nil, err
	}
	if v == nil {
		return &Record{
			Id:   item,
			Next: nil,
			Prev: nil,
		}, nil
	}

	var data models.Item
	err = proto.Unmarshal(v, &data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal item `%s`", string(item))
	}
	return &Record{
		Id:   item,
		Next: data.GetNext(),
		Prev: data.GetPrev(),
	}, nil
}

// panic if failed
func (l *KList) mustLoadState() {
	err := l.loadState()
	if err != nil {
		panic(errors.Wrap(ErrCritical, err.Error()))
	}
}

func (l *KList) loadState() error {
	v, err := l.readValue(buildStateKey(l.name))
	if err != nil {
		if err == ErrNotFound {
			l.state = &models.State{}
			return nil
		}
		return errors.Wrap(err, "failed to read state")
	}
	if v == nil {
		l.state = &models.State{}
		return nil
	}
	err = proto.Unmarshal(v, l.state)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal state")
	}
	return nil
}

func (l *KList) saveState() error {
	v, err := proto.Marshal(l.state)
	if err != nil {
		return errors.Wrap(err, "failed to build state data model")
	}
	return l.writeValue(buildStateKey(l.name), v)
}

func (l *KList) appendBatchSaveState(batch *batchapi.Batch) error {
	v, err := proto.Marshal(l.state)
	if err != nil {
		return errors.Wrap(err, "failed to build state data model")
	}
	batch.Put(buildStateKey(l.name), v)
	return nil
}

func (l *KList) appendBatchSaveRecord(batch *batchapi.Batch, rec *Record) error {
	v, err := proto.Marshal(&models.Item{
		Next: rec.Next,
		Prev: rec.Prev,
	})
	if err != nil {
		return errors.Wrap(err, "failed to build state data model")
	}
	batch.Put(buildItemKey(l.name, rec.Id), v)
	return nil
}

// cut records from queue
// batch cap=2
func (l *KList) appendBatchCutRecord(batch *batchapi.Batch, rec *Record) error {
	// update prev if exists
	if rec.Prev != nil {
		prev, err := l.readRecord(rec.Prev)
		if err != nil {
			return err
		}
		prev.Next = rec.Next
		if err = l.appendBatchSaveRecord(batch, prev); err != nil {
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
		if err = l.appendBatchSaveRecord(batch, next); err != nil {
			return err
		}
	}
	return nil
}

// insert record after item
// batch cap=4
func (l *KList) appendBatchInsertAfterRecord(batch *batchapi.Batch, rec *Record, root *Record) error {
	// in case if it's sequence rec <- root
	if isEqual(rec.Id, root.Prev) {
		root.Prev = rec.Prev
		if root.Prev != nil {
			rprev, err := l.readRecord(root.Prev)
			if err != nil {
				return err
			}
			rprev.Next = root.Id
			if err = l.appendBatchSaveRecord(batch, rprev); err != nil {
				return err
			}
		}
	}

	// update rec
	rec.Prev = root.Id
	rec.Next = root.Next
	if err := l.appendBatchSaveRecord(batch, rec); err != nil {
		return err
	}

	// update next if exists
	if root.Next != nil {
		next, err := l.readRecord(root.Next)
		if err != nil {
			return err
		}
		next.Prev = rec.Id
		if err = l.appendBatchSaveRecord(batch, next); err != nil {
			return err
		}
	}

	root.Next = rec.Id
	if err := l.appendBatchSaveRecord(batch, root); err != nil {
		return err
	}

	return nil
}

// insert record before item
// batch cap=4
func (l *KList) appendBatchInsertBeforeRecord(batch *batchapi.Batch, rec *Record, root *Record) error {
	// in case if it's sequence root -> rec
	if isEqual(rec.Id, root.Next) {
		root.Next = rec.Next
		if root.Next != nil {
			rnext, err := l.readRecord(root.Next)
			if err != nil {
				return err
			}
			rnext.Prev = root.Id
			if err = l.appendBatchSaveRecord(batch, rnext); err != nil {
				return err
			}
		}
	}

	// update rec
	rec.Prev = root.Prev
	rec.Next = root.Id
	if err := l.appendBatchSaveRecord(batch, rec); err != nil {
		return err
	}

	// update prev if exists
	if root.Prev != nil {
		prev, err := l.readRecord(root.Prev)
		if err != nil {
			return err
		}
		prev.Next = rec.Id
		if err = l.appendBatchSaveRecord(batch, prev); err != nil {
			return err
		}
	}

	// update root
	root.Prev = rec.Id
	if err := l.appendBatchSaveRecord(batch, root); err != nil {
		return err
	}

	return nil
}

func (l *KList) isItemFirst(item []byte) bool {
	return isEqual(l.state.FirstItem, item)
}

func (l *KList) isItemLast(item []byte) bool {
	return isEqual(l.state.LastItem, item)
}

func (l *KList) add(item []byte) (err error) {
	if l.isEmpty() {
		return l.writeRootItem(item)
	}

	// if list already has item -> skip adding
	exists, err := l.isItemExists(item)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	last, err := l.readRecord(l.state.LastItem)
	if err != nil {
		return err
	}

	batch := batchapi.New(5)

	// update state
	l.state.LastItem = item
	l.state.Count += 1
	if err = l.appendBatchSaveState(batch); err != nil {
		return err
	}

	// update records
	rec := Record{
		Id: item,
	}
	if err = l.appendBatchInsertAfterRecord(batch, &rec, last); err != nil {
		return err
	}

	return l.writeBatch(batch)
}

func (l *KList) setToBegin(item []byte) (err error) {
	if l.isEmpty() {
		return l.writeRootItem(item)
	}

	if l.isItemFirst(item) {
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

	first, err := l.readRecord(l.state.FirstItem)
	if err != nil {
		return err
	}

	batch := batchapi.New(7)

	// cut item
	if exists {
		rec, err = l.readRecord(item)
		if err != nil {
			return err
		}
		if err = l.appendBatchCutRecord(batch, rec); err != nil {
			return err
		}
		if l.isItemFirst(item) {
			l.state.FirstItem = rec.Next
		}
		if l.isItemLast(item) {
			l.state.LastItem = rec.Prev
		}
	}

	// update state
	if !exists {
		l.state.Count += 1
	}
	l.state.FirstItem = item
	if err = l.appendBatchSaveState(batch); err != nil {
		return err
	}

	// insert before first
	if err = l.appendBatchInsertBeforeRecord(batch, rec, first); err != nil {
		return err
	}

	return l.writeBatch(batch)
}

func (l *KList) setToEnd(item []byte) error {
	if l.isEmpty() {
		return l.writeRootItem(item)
	}

	if l.isItemLast(item) {
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

	last, err := l.readRecord(l.state.LastItem)
	if err != nil {
		return err
	}

	batch := batchapi.New(7)

	// cut item
	if exists {
		rec, err = l.readRecord(item)
		if err != nil {
			return err
		}
		if err = l.appendBatchCutRecord(batch, rec); err != nil {
			return err
		}
		if l.isItemFirst(item) {
			l.state.FirstItem = rec.Next
		}
		if l.isItemLast(item) {
			l.state.LastItem = rec.Prev
		}
	}

	// update state
	if !exists {
		l.state.Count += 1
	}
	l.state.LastItem = item
	if err = l.appendBatchSaveState(batch); err != nil {
		return err
	}

	// insert after last
	if err = l.appendBatchInsertAfterRecord(batch, rec, last); err != nil {
		return err
	}

	return l.writeBatch(batch)
}

func (l *KList) setAfter(item, root []byte) error {
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

	batch := batchapi.New(7)

	// cut item
	if exists {
		rec, err = l.readRecord(item)
		if err != nil {
			return err
		}
		if err = l.appendBatchCutRecord(batch, rec); err != nil {
			return err
		}
		if l.isItemFirst(item) {
			l.state.FirstItem = rec.Next
		}
		if l.isItemLast(item) {
			l.state.LastItem = rec.Prev
		}
	}

	// update state
	if !exists {
		l.state.Count += 1
	}
	if l.isItemLast(root) {
		l.state.LastItem = item
	}
	if err = l.appendBatchSaveState(batch); err != nil {
		return err
	}

	// insert after root
	if err = l.appendBatchInsertAfterRecord(batch, rec, rootRec); err != nil {
		return err
	}

	return l.writeBatch(batch)
}

func (l *KList) setBefore(item, root []byte) error {
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

	batch := batchapi.New(7)

	// cut item
	if exists {
		rec, err = l.readRecord(item)
		if err != nil {
			return err
		}
		if err = l.appendBatchCutRecord(batch, rec); err != nil {
			return err
		}
		if l.isItemFirst(item) {
			l.state.FirstItem = rec.Next
		}
		if l.isItemLast(item) {
			l.state.LastItem = rec.Prev
		}
	}
	// update state
	if !exists {
		l.state.Count += 1
	}
	if l.isItemFirst(root) {
		l.state.FirstItem = item
	}
	if err = l.appendBatchSaveState(batch); err != nil {
		return err
	}

	// insert after root
	if err = l.appendBatchInsertBeforeRecord(batch, rec, rootRec); err != nil {
		return err
	}

	return l.writeBatch(batch)
}

func (l *KList) pop() ([]byte, error) {
	if l.isEmpty() {
		return nil, nil
	}

	first, err := l.readRecord(l.state.FirstItem)
	if err != nil {
		return nil, err
	}

	batch := batchapi.New(4)

	// update state
	l.state.FirstItem = first.Next
	if l.isItemLast(first.Id) {
		l.state.LastItem = nil
	}
	l.state.Count -= 1
	if err = l.appendBatchSaveState(batch); err != nil {
		return nil, err
	}

	// cut from queue
	if err = l.appendBatchCutRecord(batch, first); err != nil {
		return nil, err
	}

	// delete
	batch.Delete(buildItemKey(l.name, first.Id))

	if err = l.db.Write(batch); err != nil {
		return nil, errors.Wrap(err, "failed to update data")
	}

	return first.Id, nil
}

func (l *KList) delete(item []byte) error {
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

	batch := batchapi.New(4)

	// update state
	if l.isItemFirst(item) {
		l.state.FirstItem = rec.Next
	}
	if l.isItemLast(item) {
		l.state.LastItem = rec.Prev
	}
	l.state.Count -= 1
	if err = l.appendBatchSaveState(batch); err != nil {
		return err
	}

	// cut from queue
	if err = l.appendBatchCutRecord(batch, rec); err != nil {
		return err
	}

	// delete
	batch.Delete(buildItemKey(l.name, rec.Id))

	if err = l.db.Write(batch); err != nil {
		return errors.Wrap(err, "failed to update data")
	}

	return nil
}
