package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

func setupKeeper(t *testing.T) (*Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey("sharding")
	memStoreKey := storetypes.NewMemoryStoreKey("mem_sharding")

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	k := NewKeeper(cdc, storeKey, memStoreKey)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	return k, ctx
}

func TestRegisterShard(t *testing.T) {
	k, ctx := setupKeeper(t)

	shard := &Shard{
		ID:     "shard-1",
		JobID:  "job-1",
		CID:    "QmShard123",
		Hash:   "hash123",
		NodeID: "",
		Status: "pending",
		Size:   1000,
	}

	err := k.RegisterShard(ctx, shard)
	require.NoError(t, err)

	retrievedShard, found := k.GetShard(ctx, "shard-1")
	require.True(t, found)
	require.Equal(t, shard.ID, retrievedShard.ID)
	require.Equal(t, shard.JobID, retrievedShard.JobID)

	err = k.RegisterShard(ctx, shard)
	require.Error(t, err)
}

func TestGetShard(t *testing.T) {
	k, ctx := setupKeeper(t)

	shard := &Shard{
		ID:     "shard-1",
		JobID:  "job-1",
		CID:    "QmShard123",
		Hash:   "hash123",
		NodeID: "",
		Status: "pending",
		Size:   1000,
	}

	k.RegisterShard(ctx, shard)

	retrievedShard, found := k.GetShard(ctx, "shard-1")
	require.True(t, found)
	require.Equal(t, shard.ID, retrievedShard.ID)

	_, found = k.GetShard(ctx, "nonexistent")
	require.False(t, found)
}

func TestAssignShardToNode(t *testing.T) {
	k, ctx := setupKeeper(t)

	shard := &Shard{
		ID:     "shard-1",
		JobID:  "job-1",
		CID:    "QmShard123",
		Hash:   "hash123",
		NodeID: "",
		Status: "pending",
		Size:   1000,
	}

	k.RegisterShard(ctx, shard)

	err := k.AssignShardToNode(ctx, "shard-1", "node-1")
	require.NoError(t, err)

	assignedShard, _ := k.GetShard(ctx, "shard-1")
	require.Equal(t, "node-1", assignedShard.NodeID)
	require.Equal(t, "assigned", assignedShard.Status)

	err = k.AssignShardToNode(ctx, "shard-1", "node-2")
	require.Error(t, err)

	err = k.AssignShardToNode(ctx, "nonexistent", "node-1")
	require.Error(t, err)
}

func TestGetShardsForJob(t *testing.T) {
	k, ctx := setupKeeper(t)

	shard1 := &Shard{
		ID:     "shard-1",
		JobID:  "job-1",
		CID:    "QmShard123",
		Hash:   "hash123",
		NodeID: "node-1",
		Status: "assigned",
		Size:   1000,
	}

	shard2 := &Shard{
		ID:     "shard-2",
		JobID:  "job-1",
		CID:    "QmShard456",
		Hash:   "hash456",
		NodeID: "node-2",
		Status: "assigned",
		Size:   2000,
	}

	shard3 := &Shard{
		ID:     "shard-3",
		JobID:  "job-2",
		CID:    "QmShard789",
		Hash:   "hash789",
		NodeID: "node-1",
		Status: "assigned",
		Size:   3000,
	}

	k.RegisterShard(ctx, shard1)
	k.RegisterShard(ctx, shard2)
	k.RegisterShard(ctx, shard3)

	shards := k.GetShardsForJob(ctx, "job-1")
	require.Len(t, shards, 2)

	shards = k.GetShardsForJob(ctx, "job-2")
	require.Len(t, shards, 1)

	shards = k.GetShardsForJob(ctx, "nonexistent")
	require.Len(t, shards, 0)
}

func TestGetShardsByNode(t *testing.T) {
	k, ctx := setupKeeper(t)

	shard1 := &Shard{
		ID:     "shard-1",
		JobID:  "job-1",
		CID:    "QmShard123",
		Hash:   "hash123",
		NodeID: "node-1",
		Status: "assigned",
		Size:   1000,
	}

	shard2 := &Shard{
		ID:     "shard-2",
		JobID:  "job-1",
		CID:    "QmShard456",
		Hash:   "hash456",
		NodeID: "node-1",
		Status: "assigned",
		Size:   2000,
	}

	shard3 := &Shard{
		ID:     "shard-3",
		JobID:  "job-2",
		CID:    "QmShard789",
		Hash:   "hash789",
		NodeID: "node-2",
		Status: "assigned",
		Size:   3000,
	}

	k.RegisterShard(ctx, shard1)
	k.RegisterShard(ctx, shard2)
	k.RegisterShard(ctx, shard3)

	shards := k.GetShardsByNode(ctx, "node-1")
	require.Len(t, shards, 2)

	shards = k.GetShardsByNode(ctx, "node-2")
	require.Len(t, shards, 1)

	shards = k.GetShardsByNode(ctx, "nonexistent")
	require.Len(t, shards, 0)
}

func TestGetShardsByHash(t *testing.T) {
	k, ctx := setupKeeper(t)

	shard1 := &Shard{
		ID:     "shard-1",
		JobID:  "job-1",
		CID:    "QmShard123",
		Hash:   "hash123",
		NodeID: "node-1",
		Status: "assigned",
		Size:   1000,
	}

	shard2 := &Shard{
		ID:     "shard-2",
		JobID:  "job-2",
		CID:    "QmShard456",
		Hash:   "hash123",
		NodeID: "node-2",
		Status: "assigned",
		Size:   2000,
	}

	shard3 := &Shard{
		ID:     "shard-3",
		JobID:  "job-1",
		CID:    "QmShard789",
		Hash:   "hash456",
		NodeID: "node-1",
		Status: "assigned",
		Size:   3000,
	}

	k.RegisterShard(ctx, shard1)
	k.RegisterShard(ctx, shard2)
	k.RegisterShard(ctx, shard3)

	shards := k.GetShardsByHash(ctx, "hash123")
	require.Len(t, shards, 2)

	shards = k.GetShardsByHash(ctx, "hash456")
	require.Len(t, shards, 1)

	shards = k.GetShardsByHash(ctx, "nonexistent")
	require.Len(t, shards, 0)

	shards = k.GetShardsByHash(ctx, "")
	require.Len(t, shards, 0)
}

func TestGetAllShards(t *testing.T) {
	k, ctx := setupKeeper(t)

	shard1 := &Shard{
		ID:     "shard-1",
		JobID:  "job-1",
		CID:    "QmShard123",
		Hash:   "hash123",
		NodeID: "node-1",
		Status: "assigned",
		Size:   1000,
	}

	shard2 := &Shard{
		ID:     "shard-2",
		JobID:  "job-1",
		CID:    "QmShard456",
		Hash:   "hash456",
		NodeID: "node-2",
		Status: "assigned",
		Size:   2000,
	}

	k.RegisterShard(ctx, shard1)
	k.RegisterShard(ctx, shard2)

	allShards := k.GetAllShards(ctx)
	require.Len(t, allShards, 2)
}

