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
	storagekeeper "github.com/atlas/chain/x/storage/keeper"
	"github.com/atlas/chain/x/training/types"
)

func setupKeeper(t *testing.T) (*Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	computeStoreKey := sdk.NewKVStoreKey(computetypes.StoreKey)
	storageStoreKey := sdk.NewKVStoreKey("storage")

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(computeStoreKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(storageStoreKey, storetypes.StoreTypeIAVL, db)
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
	storageKeeper := storagekeeper.NewKeeper(cdc, storageStoreKey, storetypes.NewMemoryStoreKey("mem_storage"), bankKeeper)

	k := NewKeeper(cdc, storeKey, memStoreKey, computeKeeper, storageKeeper, bankKeeper)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	return k, ctx
}

func TestGetJob(t *testing.T) {
	k, ctx := setupKeeper(t)

	job := types.Job{
		ID:         "job-1",
		ModelID:    "model-1",
		DatasetCID: "QmABC123",
		Config:     map[string]string{"epochs": "10"},
		Status:     types.TaskStatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Progress:   0.0,
		Tasks:      []string{},
	}

	k.SetJob(ctx, job)

	retrievedJob, found := k.GetJob(ctx, "job-1")
	require.True(t, found)
	require.Equal(t, job.ID, retrievedJob.ID)
	require.Equal(t, job.ModelID, retrievedJob.ModelID)

	_, found = k.GetJob(ctx, "nonexistent")
	require.False(t, found)
}

func TestSetJob(t *testing.T) {
	k, ctx := setupKeeper(t)

	job := types.Job{
		ID:         "job-1",
		ModelID:    "model-1",
		DatasetCID: "QmABC123",
		Config:     map[string]string{"epochs": "10"},
		Status:     types.TaskStatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Progress:   0.0,
		Tasks:      []string{},
	}

	k.SetJob(ctx, job)

	retrievedJob, found := k.GetJob(ctx, "job-1")
	require.True(t, found)
	require.Equal(t, job.ID, retrievedJob.ID)
}

func TestGetTask(t *testing.T) {
	k, ctx := setupKeeper(t)

	task := types.Task{
		ID:            "task-1",
		JobID:         "job-1",
		ShardID:       "shard-1",
		NodeID:        "node-1",
		Status:        types.TaskStatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Progress:      0.0,
		CheckpointCID: "",
	}

	k.SetTask(ctx, task)

	retrievedTask, found := k.GetTask(ctx, "task-1")
	require.True(t, found)
	require.Equal(t, task.ID, retrievedTask.ID)
	require.Equal(t, task.JobID, retrievedTask.JobID)

	_, found = k.GetTask(ctx, "nonexistent")
	require.False(t, found)
}

func TestSetTask(t *testing.T) {
	k, ctx := setupKeeper(t)

	task := types.Task{
		ID:            "task-1",
		JobID:         "job-1",
		ShardID:       "shard-1",
		NodeID:        "node-1",
		Status:        types.TaskStatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Progress:      0.0,
		CheckpointCID: "",
	}

	k.SetTask(ctx, task)

	retrievedTask, found := k.GetTask(ctx, "task-1")
	require.True(t, found)
	require.Equal(t, task.ID, retrievedTask.ID)
}

func TestIterateTasks(t *testing.T) {
	k, ctx := setupKeeper(t)

	task1 := types.Task{
		ID:            "task-1",
		JobID:         "job-1",
		ShardID:       "shard-1",
		NodeID:        "node-1",
		Status:        types.TaskStatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Progress:      0.0,
		CheckpointCID: "",
	}

	task2 := types.Task{
		ID:            "task-2",
		JobID:         "job-1",
		ShardID:       "shard-2",
		NodeID:        "node-2",
		Status:        types.TaskStatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Progress:      0.0,
		CheckpointCID: "",
	}

	k.SetTask(ctx, task1)
	k.SetTask(ctx, task2)

	count := 0
	k.IterateTasks(ctx, func(task types.Task) (stop bool) {
		count++
		return false
	})

	require.Equal(t, 2, count)

	count = 0
	k.IterateTasks(ctx, func(task types.Task) (stop bool) {
		count++
		return true
	})

	require.Equal(t, 1, count)
}

