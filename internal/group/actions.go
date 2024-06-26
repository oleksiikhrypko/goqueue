package group

import (
	"goqueue/internal/batch"
	"goqueue/internal/db"
	"goqueue/pkg/proto/models"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

func (g *Group) loadState() (*models.Group, error) {
	var state models.Group
	err := db.ReadStruct(g.db, buildStateKey(g.topic, g.name), &state)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read group state")
	}

	return &state, nil
}

func (g *Group) saveState(actions batch.ActionsList, state *models.Group) error {
	v, err := proto.Marshal(state)
	if err != nil {
		return errors.Wrap(err, "failed to build state data model")
	}
	actions.AddActionSet(buildStateKey(g.topic, g.name), v)
	return nil
}
