package allocation_test

import (
	"testing"

	"0chain.net/blobbercore/allocation"
	"0chain.net/blobbercore/reference"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteFileChange_ProcessChange(t *testing.T) {
	t.Skip()
	tests := []struct {
		name       string
		dfc        *allocation.DeleteFileChange
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
			name:       "failed getting ref path",
			wantErr:    true,
			wantErrMsg: "placeholder",
		},
		{
			name:       "invalid ref path",
			wantErr:    true,
			wantErrMsg: "invalid_reference_path",
		},
		{
			name:       "file_not_found",
			wantErr:    true,
			wantErrMsg: "file_not_found",
		},
		{
			name: "ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.dfc.ProcessChange(connCtx, tt.ac, tt.allocRoot)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDeleteFileChange_Marshal(t *testing.T) {
	t.Skip("tested in allocationchange_test.go")
}

func TestDeleteFileChange_Unmarshal(t *testing.T) {
	t.Skip("tested in allocationchange_test.go")
}

func TestDeleteFileChange_DeleteTempFile(t *testing.T) {
	alloc := &allocation.DeleteFileChange{}
	err := alloc.DeleteTempFile()
	require.Error(t, err)
	require.Equal(t, err, allocation.OperationNotApplicable)
}

func TestDeleteFileChange_CommitToFileStore(t *testing.T) {
	t.Skip()
	tests := []struct {
		name       string
		dfc        *allocation.DeleteFileChange
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "failed committing to filestore",
			wantErr:    true,
			wantErrMsg: "placeholder",
		},
		{
			name: "file not found",
		},
		{
			name:       "ok",
			wantErr:    true,
			wantErrMsg: "placeholder",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dfc.CommitToFileStore(connCtx)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			}
		})
	}
}
