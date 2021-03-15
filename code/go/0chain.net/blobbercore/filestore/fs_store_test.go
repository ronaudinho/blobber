package filestore_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"0chain.net/blobbercore/filestore"
	"0chain.net/core/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	store *filestore.FileFSStore
	input = map[string]*filestore.FileInputData{
		"write": &filestore.FileInputData{
			Name: "input_1",
			Path: "path_1",
			Hash: "1234567890",
		},
		"no-write": &filestore.FileInputData{
			Name: "input_2",
			Path: "path_2",
			Hash: "0987654321",
		},
	}
	output = map[string]*filestore.FileOutputData{
		"write": &filestore.FileOutputData{
			Name:        "input_1",
			Path:        "path_1",
			MerkleRoot:  "064420446ca28f56d79d0abef22339474ab036406b47b72d0756434f786e7aef",
			ContentHash: "e5fa44f2b31c1fb553b6021e7360d07d5d91ff5e",
			Size:        int64(2),
		},
	}
)

func TestMain(m *testing.M) {
	dir, err := ioutil.TempDir("", "TestFileFSStore")
	if err != nil {
		os.Exit(1)
	}
	fmt.Printf("created tmp dir at: %s\n", dir)

	store = &filestore.FileFSStore{
		RootDirectory: dir,
	}

	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

func TestGetFilePathFromHash(t *testing.T) {
	hash := "ed79cae70d439c11258236da1dfa6fc550f7cc569768304623e8fbd7d70efae4"
	wantDir := "ed7/9ca/e70"
	wantHash := "d439c11258236da1dfa6fc550f7cc569768304623e8fbd7d70efae4"
	gotDir, gotHash := filestore.GetFilePathFromHash(hash)
	assert.Equal(t, wantDir, gotDir)
	assert.Equal(t, wantHash, gotHash)
}

func TestFileFSStore_GetTotalDiskSizeUsed_empty(t *testing.T) {
	size, err := store.GetTotalDiskSizeUsed()
	require.NoError(t, err)
	assert.Equal(t, size, int64(0))
}

func TestFileFSStore_GetlDiskSizeUsed_isNotExist(t *testing.T) {
	_, err := store.GetlDiskSizeUsed("testfilefsstore_getldisksizeused_isnotexist")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestFileFSStore_GetTempPathSize_isNotExist(t *testing.T) {
	_, err := store.GetTempPathSize("testfilefsstore_gettemppathsize_isnotexist")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestFileFSStore_IterateObjects_isNotExist(t *testing.T) {
	err := store.IterateObjects("testfilefsstore_iterateobjects_isnotexist", func(hash string, size int64) {})
	require.NoError(t, err)
}

// TODO edit to match pattern lol
func TestFileFSStore_SetupAllocation_skipCreate(t *testing.T) {
	dir := store.RootDirectory
	allocID := "testfilefsstore_setupallocation_skipcreate"
	want := &filestore.StoreAllocation{
		ID:              allocID,
		Path:            fmt.Sprintf("%s//tes/tfi/lef/sstore_setupallocation_skipcreate", dir),
		ObjectsPath:     fmt.Sprintf("%s//tes/tfi/lef/sstore_setupallocation_skipcreate/objects", dir),
		TempObjectsPath: fmt.Sprintf("%s/tes/tfi/lef/sstore_setupallocation_skipcreate/objects/tmp", dir),
	}

	got, err := store.SetupAllocation(allocID, true)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestFileFSStore_File_writeDelete(t *testing.T) {
	allocID := "testfilefsstore_file_writedelete"
	in := input["write"]
	out := output["write"]

	f, err := os.Open("testdata/input_1")
	require.NoError(t, err)
	defer f.Close()
	got, err := store.WriteFile(allocID, in, f, allocID) // connection ID probably does not matter here
	require.NoError(t, err)
	require.Equal(t, out, got)

	// check if size is equal to write
	size, err := store.GetTotalDiskSizeUsed()
	require.NoError(t, err)
	assert.Equal(t, size, int64(2))
	size, err = store.GetlDiskSizeUsed(allocID)
	require.NoError(t, err)
	assert.Equal(t, size, int64(2))
	size, err = store.GetTempPathSize(allocID)
	require.NoError(t, err)
	assert.Equal(t, size, int64(2))

	// delete temp
	err = store.DeleteTempFile(allocID, in, allocID)
	require.NoError(t, err)
	err = store.DeleteFile(allocID, out.ContentHash)
	require.Error(t, err)

	// check if size is 0
	size, err = store.GetTotalDiskSizeUsed()
	require.NoError(t, err)
	assert.Equal(t, size, int64(0))
	size, err = store.GetlDiskSizeUsed(allocID)
	require.NoError(t, err)
	assert.Equal(t, size, int64(0))
	size, err = store.GetTempPathSize(allocID)
	require.NoError(t, err)
	assert.Equal(t, size, int64(0))
}

func TestFileFSStore_File_writeCommitDelete(t *testing.T) {
	allocID := "testfilefsstore_file_writecommitdelete"
	in := input["write"]
	out := output["write"]

	// TODO helper/actual testdata file
	f, err := os.Open("testdata/input_1")
	require.NoError(t, err)
	defer f.Close()
	got, err := store.WriteFile(allocID, in, f, allocID) // connection ID probably does not matter here
	require.NoError(t, err)
	require.Equal(t, out, got)

	// no need to recheck if size is equal to write
	// commit
	ok, err := store.CommitWrite(allocID, in, allocID)
	require.NoError(t, err)
	require.True(t, ok)
	size, err := store.GetTempPathSize(allocID)
	require.NoError(t, err)
	assert.Equal(t, size, int64(0))

	// delete
	err = store.DeleteFile(allocID, in.Hash)
	require.NoError(t, err)
	err = store.DeleteTempFile(allocID, in, allocID)
	require.Error(t, err)

	// check if size is 0
	size, err = store.GetTotalDiskSizeUsed()
	require.NoError(t, err)
	assert.Equal(t, size, int64(0))
	size, err = store.GetlDiskSizeUsed(allocID)
	require.NoError(t, err)
	assert.Equal(t, size, int64(0))
}

func TestFileFSStore_Block_IsNotExist(t *testing.T) {
	allocID := "testfilefsstore_block_isnotexist"
	in := input["write"]

	_, _, err := store.GetFileBlockForChallenge(allocID, in, 0)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
	_, err = store.GetFileBlock(allocID, in, int64(0), int64(0))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
	_, err = store.GetMerkleTreeForFile(allocID, in)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestFileFSStore_GetFileBlock(t *testing.T) {
	allocID := "testfilefsstore_getfileblock"
	in := input["write"]
	out := output["write"]

	f, err := os.Open("testdata/input_1")
	require.NoError(t, err)
	defer f.Close()
	got, err := store.WriteFile(allocID, in, f, allocID) // connection ID probably does not matter here
	require.NoError(t, err)
	require.Equal(t, out, got)
	ok, err := store.CommitWrite(allocID, in, allocID)
	require.NoError(t, err)
	require.True(t, ok)

	tests := []struct {
		name       string
		num        int64
		blocks     int64
		want       []byte
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "invalid block num < 1",
			num:        0,
			blocks:     1,
			wantErr:    true,
			wantErrMsg: "invalid_block_number",
		},
		{
			name:       "invalid block num > max",
			num:        2,
			blocks:     1,
			wantErr:    true,
			wantErrMsg: "invalid_block_number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.GetFileBlock(allocID, in, tt.num, tt.blocks)
			if !tt.wantErr {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				require.Contains(t, err.Error(), tt.wantErrMsg)
			}
		})
	}
}

func TestFileFSStore_GetFileBlockForChallenge(t *testing.T) {
	allocID := "testfilefsstore_getfileblock"
	in := input["write"]
	out := output["write"]

	f, err := os.Open("testdata/input_1")
	require.NoError(t, err)
	defer f.Close()
	got, err := store.WriteFile(allocID, in, f, allocID) // connection ID probably does not matter here
	require.NoError(t, err)
	require.Equal(t, out, got)
	ok, err := store.CommitWrite(allocID, in, allocID)
	require.NoError(t, err)
	require.True(t, ok)

	tests := []struct {
		name       string
		offset     int
		want       json.RawMessage
		tree       util.MerkleTreeI
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "invalid block num offset < 0",
			offset:     -1,
			wantErr:    true,
			wantErrMsg: "invalid_block_number",
		},
		{
			name:       "invalid block num offset > 1024",
			offset:     1025,
			wantErr:    true,
			wantErrMsg: "invalid_block_number",
		},
		{
			name:   "valid block num ",
			offset: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := store.GetFileBlockForChallenge(allocID, in, tt.offset)
			if !tt.wantErr {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				require.Contains(t, err.Error(), tt.wantErrMsg)
			}
		})
	}
}

func TestFileFSStore_UploadToCloud(t *testing.T) {
	t.Skip("require changing fs_store.go")
}

func TestFileFSStore_DownloadFromCloud(t *testing.T) {
	t.Skip("require changing fs_store.go")
}

func TestFileFSStore_RemoveFromCloud(t *testing.T) {
	t.Skip("require changing fs_store.go")
}
