package topic

import (
	"goqueue/pkg/proto/models"
	"goqueue/pkg/storage/batch"
	"goqueue/pkg/storage/db"
	"goqueue/pkg/storage/klist"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

func (t *Topic) loadState() (*models.Topic, error) {
	var state models.Topic
	err := db.ReadStruct(t.db, buildKey(t.name), &state)
	if err != nil {
		return nil, err
	}

	return &state, nil
}

func (t *Topic) saveState(actions batch.List, state *models.Topic) error {
	v, err := proto.Marshal(state)
	if err != nil {
		return errors.Wrap(err, "failed to build state data model")
	}
	actions.Put(buildKey(t.name), v)
	return nil
}

func (t *Topic) AddMessage(actions batch.List, key []byte, item []byte) error {
	// add key to source list
	kl := klist.New(string(buildKey(t.name)), t.db)
	err := kl.Add(actions, key)
	if err != nil {
		return errors.Wrap(err, "failed to add message to source list")
	}

	// add item to sequence list
	il := klist.New(string(buildSequenceListName(t.name, key)), t.db)
	err = il.Add(actions, item)
	if err != nil {
		return errors.Wrap(err, "failed to add message to sequence list")
	}

	return nil
}

func (t *Topic) AddGroup(actions batch.List, group string) error {
	state, err := t.loadState()
	if err != nil {
		return err
	}

	state.Groups = append(state.Groups, group)

	return t.saveState(actions, state)
}
