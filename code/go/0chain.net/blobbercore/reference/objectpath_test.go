package reference_test

import (
	"testing"

	"0chain.net/blobbercore/reference"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetObjectPath(t *testing.T) {
	tests := []struct {
		name         string
		allocationID string
		blockNum     int64
		want         *reference.ObjectPath
		wantErr      string
	}{
		{
			name:         "invalid_dir_struct",
			allocationID: "4",
			blockNum:     int64(0),
			wantErr:      "invalid_dir_struct",
		},
		{
			name:         "invalid_block_num",
			allocationID: "1",
			blockNum:     int64(11),
			wantErr:      "invalid_block_num",
		},
		{
			name:         "root numblocks is zero",
			allocationID: "3",
			blockNum:     int64(0),
			want: &reference.ObjectPath{
				Path: map[string]interface{}{
					"path":          "/",
					"num_of_blocks": int64(0),
					"type":          "d",
				},
			},
		},
		{
			name:         "failed_object_path",
			allocationID: "6",
			blockNum:     int64(1),
			wantErr:      "failed_object_path",
		},
		{
			name:         "ok",
			allocationID: "7",
			blockNum:     int64(0),
			want: &reference.ObjectPath{
				Meta: map[string]interface{}{
					"path":          "/1",
					"num_of_blocks": int64(10),
					"type":          "f",
				},
				Path: map[string]interface{}{
					"path":          "/",
					"num_of_blocks": int64(10),
					"type":          "d",
				},
				RefID: int64(18),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := reference.GetObjectPath(connCtx, tt.allocationID, tt.blockNum)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want.RootHash, got.RootHash)
			assert.Equal(t, tt.want.Meta["path"], got.Meta["path"])
			assert.Equal(t, tt.want.Meta["num_of_blocks"], got.Meta["num_of_blocks"])
			assert.Equal(t, tt.want.Meta["type"], got.Meta["type"])
			assert.Equal(t, tt.want.Path["path"], got.Path["path"])
			assert.Equal(t, tt.want.Path["num_of_blocks"], got.Path["num_of_blocks"])
			assert.Equal(t, tt.want.Path["type"], got.Path["type"])
			assert.Equal(t, tt.want.FileBlockNum, got.FileBlockNum)
			assert.Equal(t, tt.want.RefID, got.RefID)
		})
	}
}
