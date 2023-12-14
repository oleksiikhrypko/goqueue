package topic

import (
	"goqueue/pkg/storage/batch"
	"goqueue/pkg/storage/klist"

	"github.com/pkg/errors"
)

type DB interface {
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	IsNotFoundErr(err error) bool
}

type Topic struct {
	db   DB
	name string
}

func New(db DB, name string) *Topic {
	return &Topic{
		db:   db,
		name: name,
	}
}

func (t *Topic) AddItem(actions batch.List, subKey, item []byte) error {
	// add key to topic list
	kl := klist.New(t.name, t.db)
	err := kl.Add(actions, subKey)
	if err != nil {
		return errors.Wrap(err, "failed to add sub key to topic list")
	}

	// add item to sequence list
	il := klist.New(string(BuildSubListName(t.name, subKey)), t.db)
	err = il.Add(actions, item)
	if err != nil {
		return errors.Wrap(err, "failed to add message to sub list")
	}

	return nil
}

func (t *Topic) AddGroup(actions batch.List, group string) error {
	state, err := t.loadState()
	if err != nil {
		return err
	}

	for _, g := range state.Groups {
		if g == group {
			return nil
		}
	}

	state.Groups = append(state.Groups, group)

	return t.saveState(actions, state)
}
