package keeper

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	computekeeper "github.com/atlas/chain/x/compute/keeper"
	computetypes "github.com/atlas/chain/x/compute/types"
	healthkeeper "github.com/atlas/chain/x/health/keeper"
	shardingkeeper "github.com/atlas/chain/x/sharding/keeper"
	trainingkeeper "github.com/atlas/chain/x/training/keeper"
	trainingtypes "github.com/atlas/chain/x/training/types"
)

func setupKeeper(t *testing.T) (*Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey("validation")
	memStoreKey := storetypes.NewMemoryStoreKey("mem_validation")
	computeStoreKey := sdk.NewKVStoreKey(computetypes.StoreKey)
	trainingStoreKey := sdk.NewKVStoreKey("training")
	shardingStoreKey := sdk.NewKVStoreKey("sharding")
	healthStoreKey := sdk.NewKVStoreKey("health")

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(computeStoreKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(trainingStoreKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(shardingStoreKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(healthStoreKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	bankKeeper := bankkeeper.NewBaseKeeper(
		cdc,
		storeKey,
		nil,
		nil,
		banktypes.DefaultGenesisState().DenomMetadata,
	)

	computeKeeper := computekeeper.NewKeeper(cdc, computeStoreKey, storetypes.NewMemoryStoreKey("mem_compute"), bankKeeper)
	healthKeeper := healthkeeper.NewKeeper(cdc, healthStoreKey, storetypes.NewMemoryStoreKey("mem_health"), computeKeeper)
	shardingKeeper := shardingkeeper.NewKeeper(cdc, shardingStoreKey, storetypes.NewMemoryStoreKey("mem_sharding"))
	trainingKeeper := trainingkeeper.NewKeeper(cdc, trainingStoreKey, storetypes.NewMemoryStoreKey("mem_training"), computeKeeper, nil, bankKeeper)

	k := NewKeeper(cdc, storeKey, memStoreKey, shardingKeeper, trainingKeeper, computeKeeper, healthKeeper)

	ctx := sdk.NewContext(stateStore, tmproto.Header{Time: time.Now()}, false, log.NewNopLogger())

	return k, ctx
}

func TestValidateShardAssignment(t *testing.T) {
	k, ctx := setupKeeper(t)

	shard := &shardingkeeper.Shard{
		ID:     "shard-1",
		JobID:  "job-1",
		CID:    "QmShard123",
		Hash:   "hash123",
		NodeID: "",
		Status: "pending",
		Size:   1000,
	}

	k.shardingKeeper.RegisterShard(ctx, shard)

	err := k.ValidateShardAssignment(ctx, "shard-1", "node-1")
	require.NoError(t, err)

	err = k.ValidateShardAssignment(ctx, "nonexistent", "node-1")
	require.Error(t, err)

	shard.NodeID = "node-2"
	k.shardingKeeper.RegisterShard(ctx, shard)

	err = k.ValidateShardAssignment(ctx, "shard-1", "node-1")
	require.Error(t, err)

	shard2 := &shardingkeeper.Shard{
		ID:     "shard-2",
		JobID:  "job-1",
		CID:    "QmShard456",
		Hash:   "hash123",
		NodeID: "node-1",
		Status: "assigned",
		Size:   2000,
	}

	k.shardingKeeper.RegisterShard(ctx, shard2)

	err = k.ValidateShardAssignment(ctx, "shard-1", "node-1")
	require.Error(t, err)
}

func TestCheckDuplicateShard(t *testing.T) {
	k, ctx := setupKeeper(t)

	shard1 := &shardingkeeper.Shard{
		ID:     "shard-1",
		JobID:  "job-1",
		CID:    "QmShard123",
		Hash:   "hash123",
		NodeID: "node-1",
		Status: "assigned",
		Size:   1000,
	}

	k.shardingKeeper.RegisterShard(ctx, shard1)

	duplicate, err := k.CheckDuplicateShard(ctx, "hash123")
	require.NoError(t, err)
	require.True(t, duplicate)

	duplicate, err = k.CheckDuplicateShard(ctx, "hash456")
	require.NoError(t, err)
	require.False(t, duplicate)

	duplicate, err = k.CheckDuplicateShard(ctx, "")
	require.NoError(t, err)
	require.False(t, duplicate)
}

func TestValidateTaskAssignment(t *testing.T) {
	k, ctx := setupKeeper(t)

	node := computetypes.Node{
		ID:            "node-1",
		Address:       "cosmos1abc123",
		Status:        "online",
		Resources:     make(map[string]string),
		Reputation:    100.0,
		UptimePercent: 99.5,
		LastHeartbeat: ctx.BlockTime(),
	}
	k.computeKeeper.SetNode(ctx, node)

	task := trainingtypes.Task{
		ID:            "task-1",
		JobID:         "job-1",
		ShardID:       "shard-1",
		NodeID:        "",
		Status:        trainingtypes.TaskStatus_PENDING,
		CreatedAt:     ctx.BlockTime(),
		UpdatedAt:     ctx.BlockTime(),
		Progress:      0.0,
		CheckpointCID: "",
	}

	k.trainingKeeper.SetTask(ctx, task)

	err := k.ValidateTaskAssignment(ctx, "task-1", "node-1")
	require.NoError(t, err)

	err = k.ValidateTaskAssignment(ctx, "nonexistent", "node-1")
	require.Error(t, err)

	err = k.ValidateTaskAssignment(ctx, "task-1", "nonexistent")
	require.Error(t, err)

	node.Status = "offline"
	k.computeKeeper.SetNode(ctx, node)

	err = k.ValidateTaskAssignment(ctx, "task-1", "node-1")
	require.Error(t, err)
}

