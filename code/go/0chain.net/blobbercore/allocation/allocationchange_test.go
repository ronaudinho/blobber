package allocation_test

import (
	"testing"

	"0chain.net/blobbercore/allocation"
	"0chain.net/blobbercore/datastore"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllocationChanges(t *testing.T) {
	tests := []struct {
		name       string
		connID     string
		allocID    string
		clientID   string
		want       *allocation.AllocationChangeCollector
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:     "not found",
			connID:   "99",
			allocID:  "1",
			clientID: "1",
			want: &allocation.AllocationChangeCollector{
				ConnectionID: "99",
				AllocationID: "1",
				ClientID:     "1",
				Status:       allocation.CommittedConnection,
			},
		},
		{
			name:     "found",
			connID:   "1",
			allocID:  "1",
			clientID: "1",
			want: &allocation.AllocationChangeCollector{
				ConnectionID: "1",
				AllocationID: "1",
				ClientID:     "1",
			},
		},
		{
			name:     "found/status deleted",
			connID:   "2",
			allocID:  "1",
			clientID: "1",
			want: &allocation.AllocationChangeCollector{
				ConnectionID: "2",
				AllocationID: "1",
				ClientID:     "1",
				Status:       allocation.NewConnection,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := allocation.GetAllocationChanges(connCtx, tt.connID, tt.allocID, tt.clientID)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			}
			// not comparing created time
			assert.Equal(t, tt.want.ConnectionID, got.ConnectionID)
			assert.Equal(t, tt.want.AllocationID, got.AllocationID)
			assert.Equal(t, tt.want.ClientID, got.ClientID)
		})
	}
}

func TestAllocationChangeCollector_AddChange(t *testing.T) {
	tests := []struct {
		name string
		acp  allocation.AllocationChangeProcessor
		want *allocation.AllocationChangeCollector
	}{
		{
			name: "copy",
			acp:  &allocation.CopyFileChange{},
			want: &allocation.AllocationChangeCollector{
				AllocationChanges: []allocation.AllocationChangeProcessor{
					&allocation.CopyFileChange{},
				},
				Changes: []*allocation.AllocationChange{
					&allocation.AllocationChange{
						Input:       `{"connection_id":"","allocation_id":"","path":"","dest_path":""}`,
						ModelWithTS: datastore.ModelWithTS{},
					},
				},
			},
		},
		{
			name: "update",
			acp:  &allocation.UpdateFileChange{},
			want: &allocation.AllocationChangeCollector{
				AllocationChanges: []allocation.AllocationChangeProcessor{
					&allocation.UpdateFileChange{},
				},
				Changes: []*allocation.AllocationChange{
					&allocation.AllocationChange{
						Input:       `{"connection_id":"","allocation_id":"","filename":"","thumbnail_filename":"","filepath":"","size":0,"thumbnail_size":0,"actual_thumb_size":0,"actual_thumb_hash":"","attributes":{}}`,
						ModelWithTS: datastore.ModelWithTS{},
					},
				},
			},
		},
		{
			name: "delete",
			acp:  &allocation.DeleteFileChange{},
			want: &allocation.AllocationChangeCollector{
				AllocationChanges: []allocation.AllocationChangeProcessor{
					&allocation.DeleteFileChange{},
				},
				Changes: []*allocation.AllocationChange{
					&allocation.AllocationChange{
						Input:       `{"connection_id":"","allocation_id":"","name":"","path":"","size":0,"hash":"","ContentHash":null}`,
						ModelWithTS: datastore.ModelWithTS{},
					},
				},
			},
		},
		{
			name: "new",
			acp:  &allocation.NewFileChange{},
			want: &allocation.AllocationChangeCollector{
				AllocationChanges: []allocation.AllocationChangeProcessor{
					&allocation.NewFileChange{},
				},
				Changes: []*allocation.AllocationChange{
					&allocation.AllocationChange{
						Input:       `{"connection_id":"","allocation_id":"","filename":"","thumbnail_filename":"","filepath":"","size":0,"thumbnail_size":0,"actual_thumb_size":0,"actual_thumb_hash":"","attributes":{}}`,
						ModelWithTS: datastore.ModelWithTS{},
					},
				},
			},
		},
		{
			name: "attributes",
			acp:  &allocation.AttributesChange{},
			want: &allocation.AllocationChangeCollector{
				AllocationChanges: []allocation.AllocationChangeProcessor{
					&allocation.AttributesChange{},
				},
				Changes: []*allocation.AllocationChange{
					&allocation.AllocationChange{
						Input:       `{"connection_id":"","allocation_id":"","path":"","attributes":null}`,
						ModelWithTS: datastore.ModelWithTS{},
					},
				},
			},
		},
		{
			name: "rename",
			acp:  &allocation.RenameFileChange{},
			want: &allocation.AllocationChangeCollector{
				AllocationChanges: []allocation.AllocationChangeProcessor{
					&allocation.RenameFileChange{},
				},
				Changes: []*allocation.AllocationChange{
					&allocation.AllocationChange{
						Input:       `{"connection_id":"","allocation_id":"","path":"","new_name":""}`,
						ModelWithTS: datastore.ModelWithTS{},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc := &allocation.AllocationChangeCollector{}
			ac := &allocation.AllocationChange{}
			acc.AddChange(ac, tt.acp)
			assert.Equal(t, tt.want, acc)
		})
	}
}

func TestAllocationChangeCollector_Save(t *testing.T) {
	acc_1 := &allocation.AllocationChangeCollector{
		ConnectionID: "3",
		Status:       allocation.NewConnection,
	}
	acc_1a := &allocation.AllocationChangeCollector{
		ConnectionID: "3",
		Status:       allocation.NewConnection,
	}
	acc_1b := &allocation.AllocationChangeCollector{
		ConnectionID: "3",
		Status:       allocation.InProgressConnection,
	}
	want := allocation.InProgressConnection
	t.Run("first", func(t *testing.T) {
		err := acc_1.Save(connCtx)
		require.NoError(t, err)
		assert.Equal(t, want, acc_1.Status)
	})
	t.Run("repeat", func(t *testing.T) {
		err := acc_1a.Save(connCtx)
		require.Error(t, err)
	})
	t.Run("update", func(t *testing.T) {
		err := acc_1b.Save(connCtx)
		require.NoError(t, err)
		assert.Equal(t, want, acc_1b.Status)
	})
}

func TestAllocationChangeCollector_ComputeProperties(t *testing.T) {
	tests := []struct {
		name string
		acc  *allocation.AllocationChangeCollector
		want *allocation.AllocationChangeCollector
	}{
		{
			name: "insert and update",
			acc: &allocation.AllocationChangeCollector{
				Changes: []*allocation.AllocationChange{
					&allocation.AllocationChange{
						Operation: allocation.INSERT_OPERATION,
					},
					&allocation.AllocationChange{
						Operation: allocation.UPDATE_OPERATION,
					},
				},
			},
			want: &allocation.AllocationChangeCollector{
				Changes: []*allocation.AllocationChange{
					&allocation.AllocationChange{
						Operation: allocation.INSERT_OPERATION,
					},
					&allocation.AllocationChange{
						Operation: allocation.UPDATE_OPERATION,
					},
				},
				AllocationChanges: []allocation.AllocationChangeProcessor{
					&allocation.NewFileChange{},
					&allocation.UpdateFileChange{},
				},
			},
		},
		{
			name: "delete, nil, and rename",
			acc: &allocation.AllocationChangeCollector{
				Changes: []*allocation.AllocationChange{
					&allocation.AllocationChange{
						Operation: allocation.DELETE_OPERATION,
					},
					&allocation.AllocationChange{},
					&allocation.AllocationChange{
						Operation: allocation.RENAME_OPERATION,
					},
				},
			},
			want: &allocation.AllocationChangeCollector{
				Changes: []*allocation.AllocationChange{
					&allocation.AllocationChange{
						Operation: allocation.DELETE_OPERATION,
					},
					&allocation.AllocationChange{},
					&allocation.AllocationChange{
						Operation: allocation.RENAME_OPERATION,
					},
				},
				AllocationChanges: []allocation.AllocationChangeProcessor{
					&allocation.DeleteFileChange{},
					&allocation.RenameFileChange{},
				},
			},
		},
		{
			name: "copy, invalid, and attributes",
			acc: &allocation.AllocationChangeCollector{
				Changes: []*allocation.AllocationChange{
					&allocation.AllocationChange{
						Operation: allocation.COPY_OPERATION,
					},
					&allocation.AllocationChange{
						Operation: "invalid",
					},
					&allocation.AllocationChange{
						Operation: allocation.UPDATE_ATTRS_OPERATION,
					},
				},
			},
			want: &allocation.AllocationChangeCollector{
				Changes: []*allocation.AllocationChange{
					&allocation.AllocationChange{
						Operation: allocation.COPY_OPERATION,
					},
					&allocation.AllocationChange{
						Operation: "invalid",
					},
					&allocation.AllocationChange{
						Operation: allocation.UPDATE_ATTRS_OPERATION,
					},
				},
				AllocationChanges: []allocation.AllocationChangeProcessor{
					&allocation.CopyFileChange{},
					&allocation.AttributesChange{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.acc.ComputeProperties()
			assert.Equal(t, tt.want, tt.acc)
		})
	}
}

func TestAllocationChangeCollector_ApplyChanges(t *testing.T) {
	t.Skip("TODO: after tests on implementations?")
	tests := []struct {
		name       string
		acc        *allocation.AllocationChangeCollector
		allocRoot  string
		want       *allocation.AllocationChangeCollector
		wantErr    bool
		wantErrMsg string
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.acc.ApplyChanges(connCtx, tt.allocRoot)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			}
			assert.Equal(t, tt.want, tt.acc)
		})
	}
}

func TestAllocationChangeCollector_CommitToFileStore(t *testing.T) {
	t.Skip("TODO: after tests on implementations?")
	tests := []struct {
		name       string
		acc        *allocation.AllocationChangeCollector
		wantErr    bool
		wantErrMsg string
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.acc.CommitToFileStore(connCtx)
			if !tt.wantErr {
				require.NoError(t, err)
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			} else {
			}
		})
	}
}

func TestAllocationChangeCollector_DeleteChanges(t *testing.T) {
	t.Skip("TODO: after tests on implementations?")
	tests := []struct {
		name       string
		acc        *allocation.AllocationChangeCollector
		wantErr    bool
		wantErrMsg string
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.acc.DeleteChanges(connCtx)
			if !tt.wantErr {
				require.NoError(t, err)
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			} else {
			}
		})
	}
}
