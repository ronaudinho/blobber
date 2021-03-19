package allocation_test

import (
	"testing"

	"0chain.net/blobbercore/allocation"
	"0chain.net/blobbercore/reference"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAttributesChange_ProcessChange(t *testing.T) {
	t.Skip()
	db.Create(&reference.Ref{})
	tests := []struct {
		name       string
		ac         *allocation.AttributesChange
		allocRoot  string
		want       *reference.Ref
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "root ref",
			wantErr:    true,
			wantErrMsg: "getting root reference path",
		},
		{
			name:       "invalid ref path",
			wantErr:    true,
			wantErrMsg: "invalid reference path from the blobber",
		},
		{
			name:       "file not found",
			wantErr:    true,
			wantErrMsg: "file to update not found in blobber",
		},
		{
			name:       "attrs not found",
			wantErr:    true,
			wantErrMsg: "setting new attributes",
		},
		{
			name:       "invalid hash",
			wantErr:    true,
			wantErrMsg: "saving updated reference",
		},
		{
			name: "ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ac.ProcessChange(connCtx, nil, tt.allocRoot)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAttributesChange_Marshal(t *testing.T) {
	t.Skip("tested in allocationchange_test.go")
}

func TestAttributesChange_Unmarshal(t *testing.T) {
	t.Skip("tested in allocationchange_test.go")
}

func TestAttributesChange_DeleteTempFile(t *testing.T) {
	alloc := &allocation.AttributesChange{}
	err := alloc.DeleteTempFile()
	require.Error(t, err)
	require.Equal(t, err, allocation.OperationNotApplicable)
}

func TestAttributesChange_CommitToFileStore(t *testing.T) {
	alloc := &allocation.AttributesChange{}
	err := alloc.CommitToFileStore(connCtx)
	require.NoError(t, err)
}
