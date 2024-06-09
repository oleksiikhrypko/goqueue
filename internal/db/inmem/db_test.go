package inmem

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CRUD(t *testing.T) {
	db := NewDB()
	key := []byte("key1")

	v, err := db.Get(key)
	require.ErrorIs(t, err, ErrNotFound)
	require.True(t, db.IsNotFoundErr(err))
	require.Empty(t, v)
}
