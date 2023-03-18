package group

import (
	"goqueue/pkg/proto/models"
	"goqueue/pkg/storage/batch"
	"goqueue/pkg/storage/db"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

func (g *Group) loadState() (*models.Group, error) {
	var state models.Group
	err := db.ReadStruct(g.db, buildStateKey(g.name), &state)
	if err != nil {
		return nil, err
	}

	return &state, nil
}

func (g *Group) saveState(actions batch.List, state *models.Group) error {
	v, err := proto.Marshal(state)
	if err != nil {
		return errors.Wrap(err, "failed to build state data model")
	}
	actions.Put(buildStateKey(g.name), v)
	return nil
}
