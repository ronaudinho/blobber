package reference_test

import (
	"testing"

	"0chain.net/blobbercore/reference"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetReferencePath(t *testing.T) {
	tests := []struct {
		name         string
		allocationID string
		path         string
		want         *reference.Ref
		wantErr      string
	}{
		{
			name:         "invalid dir tree/root",
			allocationID: "3",
			path:         "/1",
			wantErr:      "invalid_dir_tree",
		},
		{
			name:         "invalid dir tree/child",
			allocationID: "2",
			path:         "/1/1",
			wantErr:      "invalid_dir_tree",
		},
		{
			name:         "ok",
			allocationID: "1",
			path:         "/2/1",
			want: &reference.Ref{
				AllocationID: "1",
				Path:         "/",
				PathLevel:    1,
				LookupHash:   "d1091edc9f167fc0453d8b3a054c5ed9e0b89952630236a89367b183afec65e3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := reference.GetReferencePath(connCtx, tt.allocationID, tt.path)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want.AllocationID, got.AllocationID)
			assert.Equal(t, tt.want.Path, got.Path)
			assert.Equal(t, tt.want.PathLevel, got.PathLevel)
			assert.Equal(t, tt.want.ParentPath, got.ParentPath)
			assert.Equal(t, tt.want.LookupHash, got.LookupHash)
		})
	}
}

func TestGetReferencePathFromPaths(t *testing.T) {
	tests := []struct {
		name         string
		allocationID string
		paths        []string
		want         *reference.Ref
		wantErr      string
	}{
		{
			name:         "not found",
			allocationID: "1",
			paths:        []string{"/4", "/5"},
			want: &reference.Ref{
				AllocationID: "1",
				Path:         "/",
				PathLevel:    1,
				LookupHash:   "d1091edc9f167fc0453d8b3a054c5ed9e0b89952630236a89367b183afec65e3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := reference.GetReferencePathFromPaths(connCtx, tt.allocationID, tt.paths)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want.AllocationID, got.AllocationID)
			assert.Equal(t, tt.want.Path, got.Path)
			assert.Equal(t, tt.want.PathLevel, got.PathLevel)
			assert.Equal(t, tt.want.ParentPath, got.ParentPath)
			assert.Equal(t, tt.want.LookupHash, got.LookupHash)
		})
	}
}

func TestGetObjectTree(t *testing.T) {
	tests := []struct {
		name         string
		allocationID string
		path         string
		want         *reference.Ref
		wantErr      string
	}{
		{
			name:         "invalid path",
			allocationID: "1",
			path:         "/99",
			wantErr:      "invalid_parameters",
		},
		{
			name:         "invalid object tree",
			allocationID: "2",
			path:         "/1",
			wantErr:      "invalid_object_tree",
		},
		{
			name:         "ok",
			allocationID: "1",
			path:         "/2/1",
			want: &reference.Ref{
				AllocationID: "1",
				Path:         "/2/1",
				PathLevel:    3,
				ParentPath:   "/2",
				LookupHash:   "1e48e4bfdb1abcd2e4f94526749b861ed045cb3e1580a557c587898cdab0a3ae",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := reference.GetObjectTree(connCtx, tt.allocationID, tt.path)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want.AllocationID, got.AllocationID)
			assert.Equal(t, tt.want.Path, got.Path)
			assert.Equal(t, tt.want.PathLevel, got.PathLevel)
			assert.Equal(t, tt.want.ParentPath, got.ParentPath)
			assert.Equal(t, tt.want.LookupHash, got.LookupHash)
		})
	}
}
