package topic

import (
	"goqueue/internal/batch"
	"goqueue/internal/db"
	"goqueue/pkg/proto/models"

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
	actions.AppendPut(buildKey(t.name), v)
	return nil
}
