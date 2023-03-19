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

	state.Groups = append(state.Groups, group)

	return t.saveState(actions, state)
}
