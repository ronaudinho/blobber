package allocation_test

import (
	"math"
	"testing"
	"time"

	"0chain.net/blobbercore/allocation"
	"0chain.net/core/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllocation_RestDurationInTimeUnits(t *testing.T) {
	tests := []struct {
		name  string
		alloc *allocation.Allocation
		wmt   common.Timestamp
		want  float64
	}{
		{
			name: "wmt == expt",
			alloc: &allocation.Allocation{
				Expiration: common.Timestamp(10),
				TimeUnit:   time.Duration(1),
			},
			wmt: common.Timestamp(10),
		},
		{
			name: "wmt < expt",
			alloc: &allocation.Allocation{
				Expiration: common.Timestamp(10),
				TimeUnit:   time.Duration(1),
			},
			wmt:  common.Timestamp(5),
			want: float64(5000000000),
		},
		{
			name: "wmt > expt",
			alloc: &allocation.Allocation{
				Expiration: common.Timestamp(10),
				TimeUnit:   time.Duration(1),
			},
			wmt:  common.Timestamp(15),
			want: float64(-5000000000),
		},
		{
			name: "0 wmt and expt",
			alloc: &allocation.Allocation{
				TimeUnit: time.Duration(1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.alloc.RestDurationInTimeUnits(tt.wmt)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAllocation_RestDurationInTimeUnits_NaN(t *testing.T) {
	alloc := &allocation.Allocation{
		Expiration: common.Timestamp(10),
	}
	wmt := common.Timestamp(10)
	got := alloc.RestDurationInTimeUnits(wmt)
	assert.True(t, math.IsNaN(got))
}

func TestAllocation_WantRead(t *testing.T) {
	alloc := &allocation.Allocation{
		Terms: []*allocation.Terms{
			&allocation.Terms{},
			&allocation.Terms{
				BlobberID: "1",
				ReadPrice: int64(1),
			},
			&allocation.Terms{
				BlobberID: "2",
			},
			&allocation.Terms{
				BlobberID: "3",
				ReadPrice: int64(1),
			},
		},
	}
	tests := []struct {
		name   string
		id     string
		blocks int64
		want   int64
	}{
		{
			name: "empty terms",
		},
		{
			name: "found/0 numblocks",
			id:   "1",
		},
		{
			name:   "found/0 read price",
			id:     "2",
			blocks: int64(20000),
		},
		{
			name:   "found",
			id:     "3",
			blocks: int64(20000),
			want:   int64(1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := alloc.WantRead(tt.id, tt.blocks)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAllocation_WantWrite(t *testing.T) {
	alloc := &allocation.Allocation{
		Expiration: common.Timestamp(10),
		TimeUnit:   time.Duration(1),
		Terms: []*allocation.Terms{
			&allocation.Terms{},
			&allocation.Terms{
				BlobberID:  "1",
				WritePrice: int64(1),
			},
			&allocation.Terms{
				BlobberID:  "2",
				WritePrice: int64(1),
			},
			&allocation.Terms{
				BlobberID: "3",
			},
			&allocation.Terms{
				BlobberID:  "4",
				WritePrice: int64(1),
			},
		},
	}
	tests := []struct {
		name string
		id   string
		size int64
		wmt  common.Timestamp
		want int64
	}{
		{
			name: "ignored",
			size: int64(-1),
		},
		{
			name: "empty terms",
		},
		{
			name: "found/0 size",
			id:   "1",
			wmt:  common.Timestamp(1),
		},
		{
			name: "found/0 wmt",
			id:   "2",
			size: allocation.GB,
			want: int64(10000000000),
		},
		{
			name: "found/0 write price",
			id:   "3",
			size: allocation.GB,
			wmt:  common.Timestamp(1),
		},
		{
			name: "found",
			id:   "4",
			size: allocation.GB,
			wmt:  common.Timestamp(1),
			want: int64(9000000000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := alloc.WantWrite(tt.id, tt.size, tt.wmt)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAllocation_HaveRead(t *testing.T) {
	alloc := &allocation.Allocation{
		Expiration: common.Timestamp(10),
		TimeUnit:   time.Duration(1),
		Terms: []*allocation.Terms{
			&allocation.Terms{
				BlobberID: "1",
				ReadPrice: int64(1),
			},
		},
	}
	tests := []struct {
		name   string
		rps    []*allocation.ReadPool
		id     string
		blocks int64
		want   int64
	}{
		{
			name:   "empty pool",
			rps:    []*allocation.ReadPool{},
			id:     "1",
			blocks: int64(20000),
			want:   int64(-1),
		},
		{
			name: "no balance",
			rps: []*allocation.ReadPool{
				&allocation.ReadPool{},
			},
			id:     "1",
			blocks: int64(20000),
			want:   int64(-1),
		},
		{
			name: "balance == have",
			rps: []*allocation.ReadPool{
				&allocation.ReadPool{
					Balance: int64(1),
				},
			},
			id:     "1",
			blocks: int64(20000),
		},
		{
			name: "balance > have",
			rps: []*allocation.ReadPool{
				&allocation.ReadPool{
					Balance: int64(2),
				},
			},
			id:     "1",
			blocks: int64(20000),
			want:   int64(1),
		},
		{
			name: "multi pool/balance > have",
			rps: []*allocation.ReadPool{
				&allocation.ReadPool{
					Balance: int64(1),
				},
				&allocation.ReadPool{
					Balance: int64(1),
				},
			},
			id:     "1",
			blocks: int64(20000),
			want:   int64(1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := alloc.HaveRead(tt.rps, tt.id, tt.blocks)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPending_AddPendingWrite(t *testing.T) {
	tests := []struct {
		name string
		add  int64
		want int64
	}{
		{
			name: "zero",
			add:  int64(0),
			want: int64(0),
		},
		{
			name: "positive",
			add:  int64(1),
			want: int64(1),
		},
		{
			name: "negative",
			add:  int64(-1),
			want: int64(-1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pending := &allocation.Pending{}
			pending.AddPendingWrite(tt.add)
			assert.Equal(t, tt.want, pending.PendingWrite)
		})
	}
}

func TestPending_SubPendingWrite(t *testing.T) {
	tests := []struct {
		name string
		pw   int64
		sub  int64
		want int64
	}{
		{
			name: "from negative", // should not be possible?
			pw:   int64(-1),
			sub:  int64(0),
			want: int64(0),
		},
		{
			name: "to negative",
			pw:   int64(1),
			sub:  int64(2),
			want: int64(0),
		},
		{
			name: "ok",
			pw:   int64(2),
			sub:  int64(1),
			want: int64(1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pending := &allocation.Pending{PendingWrite: tt.pw}
			pending.SubPendingWrite(tt.sub)
			assert.Equal(t, tt.want, pending.PendingWrite)
		})
	}
}

func TestPending_WritePools(t *testing.T) {
	t.Skip()
}

func TestPending_HaveWrite(t *testing.T) {
	alloc := &allocation.Allocation{
		Expiration: common.Timestamp(1),
		TimeUnit:   time.Duration(1),
		Terms: []*allocation.Terms{
			&allocation.Terms{
				BlobberID:  "1",
				WritePrice: int64(1),
			},
		},
	}
	tests := []struct {
		name string
		wps  []*allocation.WritePool
		want int64
	}{
		{
			name: "no writepool",
			wps:  []*allocation.WritePool{},
		},
		{
			name: "single writepool",
			wps: []*allocation.WritePool{
				&allocation.WritePool{
					Balance: allocation.MB,
				},
			},
			want: int64(1048576),
		},
		{
			name: "multiple writepool",
			wps: []*allocation.WritePool{
				&allocation.WritePool{
					Balance: allocation.MB,
				},
				&allocation.WritePool{
					Balance: allocation.MB,
				},
			},
			want: int64(2097152),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pending := &allocation.Pending{
				BlobberID:    "1",
				PendingWrite: allocation.KB,
			}
			got := pending.HaveWrite(tt.wps, alloc, common.Timestamp(1))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPending_Save(t *testing.T) {
	tests := []struct {
		name string
		id   int64
	}{
		{
			name: "create",
			id:   int64(0),
		},
		{
			name: "update",
			id:   int64(99),
		},
		{
			name: "repeat update",
			id:   int64(99),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pending := &allocation.Pending{
				ID: tt.id,
			}
			err := pending.Save(db)
			require.NoError(t, err)
		})
	}
}

func TestReadPools(t *testing.T) {
	db.Create(&allocation.ReadPool{
		ClientID:     "1",
		BlobberID:    "1",
		AllocationID: "1",
		ExpireAt:     common.Timestamp(10),
	})
	tests := []struct {
		name      string
		clientID  string
		allocID   string
		blobberID string
		until     common.Timestamp
		want      []*allocation.ReadPool
	}{
		{
			name:      "ok",
			clientID:  "1",
			blobberID: "1",
			allocID:   "1",
			until:     common.Timestamp(1),
			want: []*allocation.ReadPool{
				&allocation.ReadPool{
					ClientID:     "1",
					BlobberID:    "1",
					AllocationID: "1",
					ExpireAt:     common.Timestamp(10),
				},
			},
		},
		{
			name:      "client id not found",
			blobberID: "1",
			allocID:   "1",
			until:     common.Timestamp(1),
			want:      []*allocation.ReadPool{},
		},
		{
			name:     "blobber id not found",
			clientID: "1",
			allocID:  "1",
			until:    common.Timestamp(1),
			want:     []*allocation.ReadPool{},
		},
		{
			name:      "allocation id not found",
			clientID:  "1",
			blobberID: "1",
			until:     common.Timestamp(1),
			want:      []*allocation.ReadPool{},
		},
		{
			name:      "expired",
			clientID:  "1",
			blobberID: "1",
			allocID:   "1",
			until:     common.Timestamp(11),
			want:      []*allocation.ReadPool{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := allocation.ReadPools(db, tt.clientID, tt.allocID, tt.blobberID, tt.until)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetPending(t *testing.T) {
	db.Create(&allocation.Pending{
		ClientID:     "1",
		BlobberID:    "1",
		AllocationID: "1",
	})
	tests := []struct {
		name      string
		clientID  string
		allocID   string
		blobberID string
		want      *allocation.Pending
	}{
		{
			name:      "ok",
			clientID:  "1",
			blobberID: "1",
			allocID:   "1",
			want: &allocation.Pending{
				ClientID:     "1",
				AllocationID: "1",
				BlobberID:    "1",
			},
		},
		{
			name:      "client id not found",
			blobberID: "1",
			allocID:   "1",
			want: &allocation.Pending{
				AllocationID: "1",
				BlobberID:    "1",
			},
		},
		{
			name:      "allocation id not found",
			clientID:  "1",
			blobberID: "1",
			want: &allocation.Pending{
				ClientID:  "1",
				BlobberID: "1",
			},
		},
		{
			name:     "blobber id not found",
			clientID: "1",
			allocID:  "1",
			want: &allocation.Pending{
				ClientID:     "1",
				AllocationID: "1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := allocation.GetPending(db, tt.clientID, tt.allocID, tt.blobberID)
			require.NoError(t, err)
			assert.Equal(t, tt.want.ClientID, got.ClientID)
			assert.Equal(t, tt.want.AllocationID, got.AllocationID)
			assert.Equal(t, tt.want.BlobberID, got.BlobberID)
		})
	}
}

func TestSetReadPools(t *testing.T) {
	t.Skip()
}

func TestSetWritePools(t *testing.T) {
	t.Skip()
}

func TestSubReadRedeemed(t *testing.T) {
	t.Skip()
}
