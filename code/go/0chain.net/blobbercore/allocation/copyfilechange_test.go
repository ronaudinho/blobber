package allocation_test

import (
	"testing"

	"0chain.net/blobbercore/allocation"
	"0chain.net/blobbercore/reference"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyFileChange_ProcessChange(t *testing.T) {
	t.Skip()
	tests := []struct {
		name       string
		cfc        *allocation.CopyFileChange
		ac         *allocation.AllocationChange
		allocRoot  string
		want       *reference.Ref
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "invalid tree",
			wantErr:    true,
			wantErrMsg: "placeholder",
		},
		{
			name:       "invalid destination path",
			wantErr:    true,
			wantErrMsg: "Invalid destination path",
		},
		{
			name:       "invalid parent path",
			wantErr:    true,
			wantErrMsg: "placeholder",
		},
		{
			name:       "failed getting ref path",
			wantErr:    true,
			wantErrMsg: "placeholder",
		},
		{
			name:       "invalid ref path",
			wantErr:    true,
			wantErrMsg: "placeholder",
		},
		{
			name:       "file not found",
			wantErr:    true,
			wantErrMsg: "file_not_found",
		},
		{
			name:       "invalid root ref",
			wantErr:    true,
			wantErrMsg: "placeholder",
		},
		{
			name: "ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cfc.ProcessChange(connCtx, tt.ac, tt.allocRoot)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCopyFileChange_Marshal(t *testing.T) {
	t.Skip("tested in allocationchange_test.go")
}

func TestCopyFileChange_Unmarshal(t *testing.T) {
	t.Skip("tested in allocationchange_test.go")
}

func TestCopyFileChange_DeleteTempFile(t *testing.T) {
	alloc := &allocation.CopyFileChange{}
	err := alloc.DeleteTempFile()
	require.Error(t, err)
	require.Equal(t, err, allocation.OperationNotApplicable)
}

func TestCopyFileChange_CommitToFileStore(t *testing.T) {
	alloc := &allocation.CopyFileChange{}
	err := alloc.CommitToFileStore(connCtx)
	require.NoError(t, err)
}
