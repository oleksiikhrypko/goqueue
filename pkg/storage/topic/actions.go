package topic

import (
	"goqueue/pkg/proto/models"
	"goqueue/pkg/storage/batch"
	"goqueue/pkg/storage/db"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

func (t *Topic) loadState() (*models.Topic, error) {
	var state models.Topic
	err := db.ReadStruct(t.db, buildStateKey(t.name), &state)
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
	actions.Put(buildStateKey(t.name), v)
	return nil
}
