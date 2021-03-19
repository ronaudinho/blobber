package allocation_test

import (
	"testing"

	"0chain.net/blobbercore/allocation"
	"0chain.net/blobbercore/reference"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileChange_ProcessChange(t *testing.T) {
	t.Skip()
	tests := []struct {
		name       string
		nfc        *allocation.NewFileChange
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
			got, err := tt.nfc.ProcessChange(connCtx, tt.ac, tt.allocRoot)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewFileChange_Marshal(t *testing.T) {
	t.Skip("tested in allocationchange_test.go")
}

func TestNewFileChange_Unmarshal(t *testing.T) {
	t.Skip("tested in allocationchange_test.go")
}

func TestNewFileChange_DeleteTempFile(t *testing.T) {
	t.Skip()
	tests := []struct {
		name       string
		nfc        *allocation.NewFileChange
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
			err := tt.nfc.DeleteTempFile()
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			}
		})
	}
}

func TestNewFileChange_CommitToFileStore(t *testing.T) {
	t.Skip()
	tests := []struct {
		name       string
		nfc        *allocation.NewFileChange
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
			err := tt.nfc.CommitToFileStore(connCtx)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			}
		})
	}
}
