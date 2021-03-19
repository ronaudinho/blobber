package allocation_test

import (
	"testing"

	"0chain.net/blobbercore/allocation"
	"0chain.net/blobbercore/reference"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenameFileChange_ProcessChange(t *testing.T) {
	t.Skip()
	tests := []struct {
		name       string
		rfc        *allocation.RenameFileChange
		ac         *allocation.AllocationChange
		allocRoot  string
		want       *reference.Ref
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "failed getting obj tree",
			wantErr:    true,
			wantErrMsg: "placeholder",
		},
		{
			name:       "failed getting ref path",
			wantErr:    true,
			wantErrMsg: "placeholder",
		},
		{
			name:       "invalid reference path",
			wantErr:    true,
			wantErrMsg: "invalid_reference_path",
		},
		{
			name:       "file not found",
			wantErr:    true,
			wantErrMsg: "file_not_found",
		},
		{
			name:       "failed calc root hash",
			wantErr:    true,
			wantErrMsg: "placeholder",
		},
		{
			name: "ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.rfc.ProcessChange(connCtx, tt.ac, tt.allocRoot)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRenameFileChange_Marshal(t *testing.T) {
	t.Skip("tested in allocationchange_test.go")
}

func TestRenameFileChange_Unmarshal(t *testing.T) {
	t.Skip("tested in allocationchange_test.go")
}

func TestRenameFileChange_DeleteTempFile(t *testing.T) {
	alloc := &allocation.RenameFileChange{}
	err := alloc.DeleteTempFile()
	require.Error(t, err)
	require.Equal(t, err, allocation.OperationNotApplicable)
}

func TestRenameFileChange_CommitToFileStore(t *testing.T) {
	alloc := &allocation.RenameFileChange{}
	err := alloc.CommitToFileStore(connCtx)
	require.NoError(t, err)
}
