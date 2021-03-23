package reference_test

import (
	"testing"

	"0chain.net/blobbercore/reference"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollaborator(t *testing.T) {
	got := reference.IsACollaborator(connCtx, int64(1), "1")
	assert.False(t, got)

	err := reference.RemoveCollaborator(connCtx, int64(1), "1")
	require.NoError(t, err)

	err = reference.AddCollaborator(connCtx, int64(1), "1")
	require.NoError(t, err)

	colls, err := reference.GetCollaborators(connCtx, int64(1))
	require.NoError(t, err)
	assert.Equal(t, 1, len(colls))

	got = reference.IsACollaborator(connCtx, int64(1), "1")
	require.NoError(t, err)
	assert.True(t, got)

	err = reference.AddCollaborator(connCtx, int64(1), "1")
	require.NoError(t, err)

	err = reference.RemoveCollaborator(connCtx, int64(1), "1")
	require.NoError(t, err)

	colls, err = reference.GetCollaborators(connCtx, int64(1))
	require.NoError(t, err)
	assert.Equal(t, 0, len(colls))
}
