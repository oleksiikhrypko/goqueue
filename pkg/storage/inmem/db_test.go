package inmem

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func Test_CRUD(t *testing.T) {
	ctx := context.Background()

	db := NewDB(ctx)
	key := []byte("key1")

	v, err := db.Get(key)
	require.ErrorIs(t, err, ErrNotFound)
	require.Empty(t, v)
}
