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
	trainingkeeper "github.com/atlas/chain/x/training/keeper"
	trainingtypes "github.com/atlas/chain/x/training/types"
)

func setupKeeper(t *testing.T) (*Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey("recovery")
	memStoreKey := storetypes.NewMemoryStoreKey("mem_recovery")
	computeStoreKey := sdk.NewKVStoreKey(computetypes.StoreKey)
	trainingStoreKey := sdk.NewKVStoreKey("training")
	healthStoreKey := sdk.NewKVStoreKey("health")

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(computeStoreKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(trainingStoreKey, storetypes.StoreTypeIAVL, db)
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
	trainingKeeper := trainingkeeper.NewKeeper(cdc, trainingStoreKey, storetypes.NewMemoryStoreKey("mem_training"), computeKeeper, nil, bankKeeper)

	k := NewKeeper(cdc, storeKey, memStoreKey, trainingKeeper, computeKeeper, healthKeeper)

	ctx := sdk.NewContext(stateStore, tmproto.Header{Time: time.Now()}, false, log.NewNopLogger())

	return k, ctx
}

func TestRollbackTasksForNode(t *testing.T) {
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

	task1 := trainingtypes.Task{
		ID:            "task-1",
		JobID:         "job-1",
		ShardID:       "shard-1",
		NodeID:        "node-1",
		Status:        trainingtypes.TaskStatus_IN_PROGRESS,
		CreatedAt:     ctx.BlockTime(),
		UpdatedAt:     ctx.BlockTime(),
		Progress:      0.5,
		CheckpointCID: "",
	}

	task2 := trainingtypes.Task{
		ID:            "task-2",
		JobID:         "job-1",
		ShardID:       "shard-2",
		NodeID:        "node-1",
		Status:        trainingtypes.TaskStatus_ASSIGNED,
		CreatedAt:     ctx.BlockTime(),
		UpdatedAt:     ctx.BlockTime(),
		Progress:      0.0,
		CheckpointCID: "",
	}

	task3 := trainingtypes.Task{
		ID:            "task-3",
		JobID:         "job-1",
		ShardID:       "shard-3",
		NodeID:        "node-2",
		Status:        trainingtypes.TaskStatus_IN_PROGRESS,
		CreatedAt:     ctx.BlockTime(),
		UpdatedAt:     ctx.BlockTime(),
		Progress:      0.3,
		CheckpointCID: "",
	}

	k.trainingKeeper.SetTask(ctx, task1)
	k.trainingKeeper.SetTask(ctx, task2)
	k.trainingKeeper.SetTask(ctx, task3)

	err := k.RollbackTasksForNode(ctx, "node-1")
	require.NoError(t, err)

	rolledBackTask1, _ := k.trainingKeeper.GetTask(ctx, "task-1")
	require.Equal(t, trainingtypes.TaskStatus_PENDING, rolledBackTask1.Status)
	require.Equal(t, "", rolledBackTask1.NodeID)

	rolledBackTask2, _ := k.trainingKeeper.GetTask(ctx, "task-2")
	require.Equal(t, trainingtypes.TaskStatus_PENDING, rolledBackTask2.Status)

	task3After, _ := k.trainingKeeper.GetTask(ctx, "task-3")
	require.Equal(t, trainingtypes.TaskStatus_IN_PROGRESS, task3After.Status)
	require.Equal(t, "node-2", task3After.NodeID)
}

func TestReassignTask(t *testing.T) {
	k, ctx := setupKeeper(t)

	node1 := computetypes.Node{
		ID:            "node-1",
		Address:       "cosmos1abc123",
		Status:        "online",
		Resources:     make(map[string]string),
		Reputation:    100.0,
		UptimePercent: 99.5,
		LastHeartbeat: ctx.BlockTime(),
	}

	node2 := computetypes.Node{
		ID:            "node-2",
		Address:       "cosmos1def456",
		Status:        "online",
		Resources:     make(map[string]string),
		Reputation:    95.0,
		UptimePercent: 98.0,
		LastHeartbeat: ctx.BlockTime(),
	}

	k.computeKeeper.SetNode(ctx, node1)
	k.computeKeeper.SetNode(ctx, node2)

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

	err := k.ReassignTask(ctx, "task-1", "node-2")
	require.NoError(t, err)

	reassignedTask, _ := k.trainingKeeper.GetTask(ctx, "task-1")
	require.Equal(t, "node-2", reassignedTask.NodeID)
	require.Equal(t, trainingtypes.TaskStatus_ASSIGNED, reassignedTask.Status)

	err = k.ReassignTask(ctx, "nonexistent", "node-2")
	require.Error(t, err)

	err = k.ReassignTask(ctx, "task-1", "nonexistent")
	require.Error(t, err)
}

func TestHandleNodeOffline(t *testing.T) {
	k, ctx := setupKeeper(t)

	node1 := computetypes.Node{
		ID:            "node-1",
		Address:       "cosmos1abc123",
		Status:        "online",
		Resources:     make(map[string]string),
		Reputation:    100.0,
		UptimePercent: 99.5,
		LastHeartbeat: ctx.BlockTime(),
	}

	node2 := computetypes.Node{
		ID:            "node-2",
		Address:       "cosmos1def456",
		Status:        "online",
		Resources:     make(map[string]string),
		Reputation:    95.0,
		UptimePercent: 98.0,
		LastHeartbeat: ctx.BlockTime(),
	}

	k.computeKeeper.SetNode(ctx, node1)
	k.computeKeeper.SetNode(ctx, node2)

	task1 := trainingtypes.Task{
		ID:            "task-1",
		JobID:         "job-1",
		ShardID:       "shard-1",
		NodeID:        "node-1",
		Status:        trainingtypes.TaskStatus_IN_PROGRESS,
		CreatedAt:     ctx.BlockTime(),
		UpdatedAt:     ctx.BlockTime(),
		Progress:      0.5,
		CheckpointCID: "",
	}

	task2 := trainingtypes.Task{
		ID:            "task-2",
		JobID:         "job-1",
		ShardID:       "shard-2",
		NodeID:        "",
		Status:        trainingtypes.TaskStatus_PENDING,
		CreatedAt:     ctx.BlockTime(),
		UpdatedAt:     ctx.BlockTime(),
		Progress:      0.0,
		CheckpointCID: "",
	}

	k.trainingKeeper.SetTask(ctx, task1)
	k.trainingKeeper.SetTask(ctx, task2)

	err := k.HandleNodeOffline(ctx, "node-1")
	require.NoError(t, err)

	rolledBackTask, _ := k.trainingKeeper.GetTask(ctx, "task-1")
	require.Equal(t, trainingtypes.TaskStatus_PENDING, rolledBackTask.Status)
}

