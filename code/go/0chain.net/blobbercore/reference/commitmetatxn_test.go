package reference_test

import (
	"testing"

	"0chain.net/blobbercore/reference"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommitMetaTxn(t *testing.T) {
	err := reference.AddCommitMetaTxn(connCtx, int64(1), "1")
	require.NoError(t, err)

	txn, err := reference.GetCommitMetaTxns(connCtx, int64(1))
	require.NoError(t, err)
	assert.Equal(t, 1, len(txn))

	err = reference.AddCommitMetaTxn(connCtx, int64(1), "1")
	require.NoError(t, err)

	txn, err = reference.GetCommitMetaTxns(connCtx, int64(1))
	require.NoError(t, err)
	assert.Equal(t, 2, len(txn))
}
