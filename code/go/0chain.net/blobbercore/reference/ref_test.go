package reference_test

import (
	"encoding/json"
	"testing"
	"time"

	"0chain.net/blobbercore/reference"
	"0chain.net/core/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

func TestAttributes_IsZero(t *testing.T) {
	tests := []struct {
		name  string
		attrs *reference.Attributes
		want  bool
	}{
		{
			name:  "empty",
			attrs: &reference.Attributes{},
			want:  true,
		},
		{
			name: "owner",
			attrs: &reference.Attributes{
				WhoPaysForReads: common.WhoPaysOwner,
			},
			want: true, // NOTE probably unintended
		},
		{
			name: "3rd party",
			attrs: &reference.Attributes{
				WhoPaysForReads: common.WhoPays3rdParty,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.attrs.IsZero())
		})
	}
}

func TestAttributes_Validate(t *testing.T) {
	tests := []struct {
		name    string
		attrs   *reference.Attributes
		wantErr bool
	}{
		{
			name:    "empty",
			attrs:   &reference.Attributes{},
			wantErr: false, // NOTE probably unintended
		},
		{
			name: "owner",
			attrs: &reference.Attributes{
				WhoPaysForReads: common.WhoPaysOwner,
			},
			wantErr: false,
		},
		{
			name: "3rd party",
			attrs: &reference.Attributes{
				WhoPaysForReads: common.WhoPays3rdParty,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.attrs.Validate()
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestGetReferenceLookup(t *testing.T) {
	want := "2cf13675e7f85a2b6ad6c42fd66b306b6445d0f0b30420a4b4b24e4ec333d6f0"
	got := reference.GetReferenceLookup("allocation_id", "path")
	assert.Equal(t, want, got)
}

func TestNewDirectoryRef(t *testing.T) {
	want := &reference.Ref{Type: reference.DIRECTORY, Attributes: datatypes.JSON("{}")}
	got := reference.NewDirectoryRef()
	assert.Equal(t, want, got)
}

func TestNewFileRef(t *testing.T) {
	want := &reference.Ref{Type: reference.FILE, Attributes: datatypes.JSON("{}")}
	got := reference.NewFileRef()
	assert.Equal(t, want, got)
}

func TestRef_GetAttributes(t *testing.T) {
	tests := []struct {
		name    string
		ref     *reference.Ref
		want    *reference.Attributes
		wantErr string
	}{
		{
			name: "invalid",
			ref: &reference.Ref{
				Attributes: datatypes.JSON(json.RawMessage(`{"who_pays_for_reads": }`)),
			},
			wantErr: "decoding file attributes",
		},
		{
			name: "empty",
			ref: &reference.Ref{
				Attributes: datatypes.JSON(json.RawMessage(``)),
			},
			want: &reference.Attributes{},
		},
		{
			name: "empty/who pays owner?",
			ref: &reference.Ref{
				Attributes: datatypes.JSON(json.RawMessage(``)),
			},
			want: &reference.Attributes{
				WhoPaysForReads: 0,
			},
		},
		{
			name: "valid",
			ref: &reference.Ref{
				Attributes: datatypes.JSON(json.RawMessage(`{"who_pays_for_reads": 0}`)),
			},
			want: &reference.Attributes{
				WhoPaysForReads: 0,
			},
		},
		{
			name: "valid/like empty?",
			ref: &reference.Ref{
				Attributes: datatypes.JSON(json.RawMessage(`{"who_pays_for_reads": 0}`)),
			},
			want: &reference.Attributes{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ref.GetAttributes()
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRef_SetAttributes(t *testing.T) {
	tests := []struct {
		name  string
		attrs *reference.Attributes
		want  datatypes.JSON
	}{
		{
			name:  "nil",
			attrs: nil,
			want:  datatypes.JSON("{}"),
		},
		{
			name:  "empty",
			attrs: &reference.Attributes{},
			want:  datatypes.JSON("{}"),
		},
		{
			name: "owner/like empty?",
			attrs: &reference.Attributes{
				WhoPaysForReads: common.WhoPaysOwner,
			},
			want: datatypes.JSON("{}"),
		},
		{
			name: "3rd party",
			attrs: &reference.Attributes{
				WhoPaysForReads: common.WhoPays3rdParty,
			},
			want: datatypes.JSON(json.RawMessage(`{"who_pays_for_reads":1}`)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref := &reference.Ref{}
			err := ref.SetAttributes(tt.attrs)
			require.NoError(t, err)
			assert.Equal(t, tt.want, ref.Attributes)
		})
	}
}

func TestGetReference(t *testing.T) {
	want := &reference.Ref{
		AllocationID: "1",
		Path:         "/1",
		PathLevel:    2,
		ParentPath:   "/",
		LookupHash:   "27116ad042d0ff98060681fd5c94ad9ad1a317668726c4d282fc88affbced275",
	}
	got, err := reference.GetReference(connCtx, "1", "/1")
	require.NoError(t, err)
	assert.Equal(t, want.AllocationID, got.AllocationID)
	assert.Equal(t, want.Path, got.Path)
	assert.Equal(t, want.PathLevel, got.PathLevel)
	assert.Equal(t, want.ParentPath, got.ParentPath)
	assert.Equal(t, want.LookupHash, got.LookupHash)

	_, err = reference.GetReference(connCtx, "1", "/99")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
}

func TestGetReferenceFromLookupHash(t *testing.T) {
	want := &reference.Ref{
		AllocationID: "1",
		Path:         "/2/1",
		PathLevel:    3,
		ParentPath:   "/2",
		LookupHash:   "1e48e4bfdb1abcd2e4f94526749b861ed045cb3e1580a557c587898cdab0a3ae",
	}
	got, err := reference.GetReferenceFromLookupHash(connCtx, "1", "1e48e4bfdb1abcd2e4f94526749b861ed045cb3e1580a557c587898cdab0a3ae")
	require.NoError(t, err)
	assert.Equal(t, want.AllocationID, got.AllocationID)
	assert.Equal(t, want.Path, got.Path)
	assert.Equal(t, want.PathLevel, got.PathLevel)
	assert.Equal(t, want.ParentPath, got.ParentPath)
	assert.Equal(t, want.LookupHash, got.LookupHash)

	_, err = reference.GetReferenceFromLookupHash(connCtx, "1", "1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
}

func TestGetSubDirsFromPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want []string
	}{
		{
			name: "root",
			path: "/",
			want: []string{},
		},
		{
			name: "one",
			path: "/1",
			want: []string{"1"},
		},
		{
			name: "two",
			path: "/1/2",
			want: []string{"1", "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reference.GetSubDirsFromPath(tt.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetRefWithChildren(t *testing.T) {
	tests := []struct {
		name         string
		allocationID string
		path         string
		want         *reference.Ref
		wantErr      string
	}{
		{
			name:         "not found",
			allocationID: "99",
			path:         "/1",
			want: &reference.Ref{
				Type:         reference.DIRECTORY,
				AllocationID: "99",
				Path:         "/1",
			},
		},
		{
			name:         "invalid dir tree/root",
			allocationID: "5",
			path:         "/1",
			wantErr:      "invalid_dir_tree",
		},
		{
			name:         "invalid dir tree/child",
			allocationID: "5",
			path:         "/1",
			wantErr:      "invalid_dir_tree",
		},
		{
			name:         "ok",
			allocationID: "1",
			path:         "/2",
			want: &reference.Ref{
				AllocationID: "1",
				Path:         "/2",
				ParentPath:   "/",
				PathLevel:    2,
				LookupHash:   "f7621491220df9fb02d17d764a3935fb2e886496c772e1075408191f81173ea0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := reference.GetRefWithChildren(connCtx, tt.allocationID, tt.path)
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

func TestGetRefWithSortedChildren(t *testing.T) {
	tests := []struct {
		name         string
		allocationID string
		path         string
		want         *reference.Ref
		wantErr      string
	}{
		{
			name:         "not found",
			allocationID: "99",
			path:         "/1",
			want: &reference.Ref{
				Type:         reference.DIRECTORY,
				AllocationID: "99",
				Path:         "/1",
			},
		},
		{
			name:         "invalid dir tree/root",
			allocationID: "5",
			path:         "/1",
			wantErr:      "invalid_dir_tree",
		},
		{
			name:         "invalid dir tree/child",
			allocationID: "5",
			path:         "/1",
			wantErr:      "invalid_dir_tree",
		},
		{
			name:         "ok",
			allocationID: "1",
			path:         "/2",
			want: &reference.Ref{
				AllocationID: "1",
				Path:         "/2",
				ParentPath:   "/",
				PathLevel:    2,
				LookupHash:   "f7621491220df9fb02d17d764a3935fb2e886496c772e1075408191f81173ea0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := reference.GetRefWithSortedChildren(connCtx, tt.allocationID, tt.path)
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

func TestRef_GetFileHashData(t *testing.T) {
}

func TestRef_CalculateHash(t *testing.T) {
	tests := []struct {
		name     string
		save     bool
		ref      *reference.Ref
		children []*reference.Ref
		want     string
	}{
		{
			name: "no children file",
			ref: &reference.Ref{
				AllocationID: "5",
				Path:         "/1",
				PathLevel:    1,
				Type:         reference.FILE,
			},
			want: "db44285b06f02ef3242fa953942e6278ba73d919a291666a74a03384e7bc20f6",
		},
		{
			name: "no children dir",
			ref: &reference.Ref{
				AllocationID: "5",
				Path:         "/1",
				PathLevel:    1,
				Type:         reference.DIRECTORY,
			},
			want: "",
		},
		{
			name: "with children",
			ref: &reference.Ref{
				AllocationID: "5",
				Path:         "/",
				PathLevel:    1,
				Type:         reference.DIRECTORY,
				Children: []*reference.Ref{
					&reference.Ref{
						AllocationID: "5",
						Path:         "/1",
						PathLevel:    2,
						ParentPath:   "/",
						Type:         reference.FILE,
					},
				},
			},
			want: "81dc4eedd52487ffbc7c15ffa3b8e84445438f2774f9bcac0cda40b0d5205943",
		},
		{
			name: "with children loaded",
			ref: &reference.Ref{
				AllocationID: "5",
				Path:         "/",
				PathLevel:    1,
				Type:         reference.DIRECTORY,
			},
			children: []*reference.Ref{
				&reference.Ref{
					AllocationID: "5",
					Path:         "/1",
					PathLevel:    2,
					ParentPath:   "/",
					Type:         reference.FILE,
				},
			},
			want: "81dc4eedd52487ffbc7c15ffa3b8e84445438f2774f9bcac0cda40b0d5205943",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, child := range tt.children {
				tt.ref.AddChild(child)
			}
			got, err := tt.ref.CalculateHash(connCtx, tt.save)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRef_AddChild(t *testing.T) {
	want_1 := []*reference.Ref{
		&reference.Ref{
			Path:       "/1",
			PathLevel:  2,
			LookupHash: "2",
			Type:       reference.FILE,
		},
	}
	want_2 := []*reference.Ref{
		&reference.Ref{
			Path:       "/2",
			PathLevel:  2,
			LookupHash: "1",
			Type:       reference.FILE,
		},
		&reference.Ref{
			Path:       "/1",
			PathLevel:  2,
			LookupHash: "2",
			Type:       reference.FILE,
		},
	}
	ref := &reference.Ref{
		Path:      "/",
		PathLevel: 1,
		Type:      reference.DIRECTORY,
	}

	ref.AddChild(&reference.Ref{
		Path:       "/1",
		PathLevel:  2,
		LookupHash: "2",
		Type:       reference.FILE,
	})
	assert.Equal(t, want_1, ref.Children)

	ref.AddChild(&reference.Ref{
		Path:       "/2",
		PathLevel:  2,
		LookupHash: "1",
		Type:       reference.FILE,
	})
	assert.Equal(t, want_2, ref.Children)
}

func TestRef_RemoveChild(t *testing.T) {
	ref := &reference.Ref{
		Children: []*reference.Ref{
			&reference.Ref{
				Path:      "/1",
				PathLevel: 2,
				Type:      reference.FILE,
			},
			&reference.Ref{
				Path:      "/2",
				PathLevel: 2,
				Type:      reference.FILE,
			},
		},
	}

	ref.RemoveChild(-1)
	assert.Equal(t, 2, len(ref.Children))

	ref.RemoveChild(1)
	assert.Equal(t, 1, len(ref.Children))

	// NOTE cannot remove all children
	// ref.RemoveChild(1)
	// assert.Equal(t, 0, len(ref.Children))
}

func TestRef_UpdatePath(t *testing.T) {
	want := &reference.Ref{
		Path:       "/path",
		ParentPath: "/parent",
		PathLevel:  2,
		LookupHash: "b2be7d55b69c1084866e2ddc0b50ee3eeac45bf6e9c0cf67b9c4c22992cb17f1",
	}
	ref := &reference.Ref{}
	ref.UpdatePath("/path", "/parent")
	assert.Equal(t, want.Path, ref.Path)
	assert.Equal(t, want.ParentPath, ref.ParentPath)
	assert.Equal(t, want.PathLevel, ref.PathLevel)
	assert.Equal(t, want.LookupHash, ref.LookupHash)
}

func TestDeleteReference(t *testing.T) {
	tests := []struct {
		name    string
		id      int64
		hash    string
		wantErr string
	}{
		{
			name:    "invalid ref id",
			id:      int64(-1),
			wantErr: "invalid_ref_id",
		},
		{
			name: "not found",
			id:   int64(99),
		},
		{
			name: "found",
			id:   int64(10),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := reference.DeleteReference(connCtx, tt.id, tt.hash)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRef_Save(t *testing.T) {
	ref := &reference.Ref{}
	require.NoError(t, ref.Save(connCtx))
}

func TestRef_GetListingData(t *testing.T) {
	tests := []struct {
		name string
		ref  *reference.Ref
		want map[string]interface{}
	}{
		{
			name: "dir",
			ref: &reference.Ref{
				AllocationID: "1",
				Path:         "/",
				PathLevel:    1,
				Type:         reference.DIRECTORY,
			},
			want: map[string]interface{}{
				"created_at":    time.Time{},
				"hash":          "",
				"lookup_hash":   "",
				"name":          "",
				"num_of_blocks": int64(0),
				"path":          "/",
				"path_hash":     "",
				"size":          int64(0),
				"type":          "d",
				"updated_at":    time.Time{},
			},
		},
		{
			name: "file",
			ref: &reference.Ref{
				AllocationID: "1",
				Path:         "/",
				PathLevel:    1,
				Type:         reference.FILE,
			},
			want: map[string]interface{}{
				"actual_file_hash":      "",
				"actual_file_size":      int64(0),
				"actual_thumbnail_hash": "",
				"actual_thumbnail_size": int64(0),
				"attributes":            datatypes.JSON(nil),
				"commit_meta_txns":      []reference.CommitMetaTxn(nil),
				"content_hash":          "",
				"created_at":            time.Time{},
				"custom_meta":           "",
				"encrypted_key":         "",
				"hash":                  "",
				"lookup_hash":           "",
				"merkle_root":           "",
				"mimetype":              "",
				"name":                  "",
				"num_of_blocks":         int64(0),
				"on_cloud":              false,
				"path":                  "/",
				"path_hash":             "",
				"size":                  int64(0),
				"thumbnail_hash":        "",
				"thumbnail_size":        int64(0),
				"type":                  "f",
				"updated_at":            time.Time{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ref.GetListingData(connCtx)
			assert.Equal(t, tt.want, got)
		})
	}
}
