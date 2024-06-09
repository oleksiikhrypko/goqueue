package leveldb

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/net/context"
)

func Test_CRUD(t *testing.T) {
	ctx := context.Background()

	lvl, err := leveldb.OpenFile("tmp/test_klist", nil)
	if err != nil {
		require.NoError(t, err)
		return
	}
	defer lvl.Close()

	db := NewDB(ctx, lvl)
	key := []byte("key1")

	v, err := db.Get(key)
	require.Equal(t, err, leveldb.ErrNotFound)
	require.True(t, db.IsNotFoundErr(err))
	require.Empty(t, v)

	// TODO: add more...
}
