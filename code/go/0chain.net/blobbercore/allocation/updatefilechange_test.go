package allocation_test

import (
	"testing"

	"0chain.net/blobbercore/allocation"
	"0chain.net/blobbercore/reference"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateFileChange_ProcessChange(t *testing.T) {
	t.Skip()
	tests := []struct {
		name       string
		ufc        *allocation.UpdateFileChange
		ac         *allocation.AllocationChange
		allocRoot  string
		want       *reference.Ref
		wantErr    bool
		wantErrMsg string
	}{
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
			name:       "failed setting file attrs",
			wantErr:    true,
			wantErrMsg: "setting file attributes",
		},
		{
			name: "ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ufc.ProcessChange(connCtx, tt.ac, tt.allocRoot)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUpdateFileChange_Marshal(t *testing.T) {
	t.Skip("tested in allocationchange_test.go")
}

func TestUpdateFileChange_Unmarshal(t *testing.T) {
	t.Skip("tested in allocationchange_test.go")
}

func TestUpdateFileChange_DeleteTempFile(t *testing.T) {
	t.Skip()
	tests := []struct {
		name       string
		ufc        *allocation.UpdateFileChange
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "failed committing to filestore",
			wantErr:    true,
			wantErrMsg: "placeholder",
		},
		{
			name:       "thumbnail/failed committing to filestore",
			wantErr:    true,
			wantErrMsg: "setting file attributes",
		},
		{
			name: "ok",
		},
		{
			name: "thumbnail/ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ufc.DeleteTempFile()
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			}
		})
	}
}

func TestUpdateFileChange_CommitToFileStore(t *testing.T) {
	t.Skip()
	tests := []struct {
		name       string
		ufc        *allocation.UpdateFileChange
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "failed committing to filestore",
			wantErr:    true,
			wantErrMsg: "placeholder",
		},
		{
			name:       "thumbnail/failed committing to filestore",
			wantErr:    true,
			wantErrMsg: "setting file attributes",
		},
		{
			name: "ok",
		},
		{
			name: "thumbnail/ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ufc.CommitToFileStore(connCtx)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			}
		})
	}
}
