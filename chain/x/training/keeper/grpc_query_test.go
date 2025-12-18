package keeper

import (
	"context"
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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	computekeeper "github.com/atlas/chain/x/compute/keeper"
	computetypes "github.com/atlas/chain/x/compute/types"
	storagekeeper "github.com/atlas/chain/x/storage/keeper"
	"github.com/atlas/chain/x/training/types"
)

func setupQueryServer(t *testing.T) (QueryServer, sdk.Context) {
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

	keeper := NewKeeper(cdc, storeKey, memStoreKey, computeKeeper, storageKeeper, bankKeeper)
	qs := NewQueryServer(keeper)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	return qs, ctx
}

func TestGetJob(t *testing.T) {
	qs, ctx := setupQueryServer(t)

	job := types.Job{
		ID:         "job-1",
		ModelID:    "model-1",
		DatasetCID: "QmABC123",
		Config:     map[string]interface{}{"epochs": "10"},
		Status:     types.TaskStatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Progress:   0.0,
		Tasks:      []string{},
	}

	qs.Keeper.SetJob(ctx, job)

	req := &types.QueryGetJobRequest{JobId: "job-1"}
	resp, err := qs.GetJob(sdk.WrapSDKContext(ctx), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, job.ID, resp.Job.ID)

	req.JobId = ""
	_, err = qs.GetJob(sdk.WrapSDKContext(ctx), req)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))

	req.JobId = "nonexistent"
	_, err = qs.GetJob(sdk.WrapSDKContext(ctx), req)
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))

	_, err = qs.GetJob(context.Background(), nil)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestListJobs(t *testing.T) {
	qs, ctx := setupQueryServer(t)

	job1 := types.Job{
		ID:         "job-1",
		ModelID:    "model-1",
		DatasetCID: "QmABC123",
		Config:     map[string]interface{}{"epochs": "10"},
		Status:     types.TaskStatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Progress:   0.0,
		Tasks:      []string{},
	}

	job2 := types.Job{
		ID:         "job-2",
		ModelID:    "model-2",
		DatasetCID: "QmDEF456",
		Config:     map[string]interface{}{"epochs": "20"},
		Status:     types.TaskStatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Progress:   0.0,
		Tasks:      []string{},
	}

	qs.Keeper.SetJob(ctx, job1)
	qs.Keeper.SetJob(ctx, job2)

	req := &types.QueryListJobsRequest{}
	resp, err := qs.ListJobs(sdk.WrapSDKContext(ctx), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Jobs, 2)

	_, err = qs.ListJobs(context.Background(), nil)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestGetTask(t *testing.T) {
	qs, ctx := setupQueryServer(t)

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

	qs.Keeper.SetTask(ctx, task)

	req := &types.QueryGetTaskRequest{TaskId: "task-1"}
	resp, err := qs.GetTask(sdk.WrapSDKContext(ctx), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, task.ID, resp.Task.ID)

	req.TaskId = ""
	_, err = qs.GetTask(sdk.WrapSDKContext(ctx), req)
	require.Error(t, err)

	req.TaskId = "nonexistent"
	_, err = qs.GetTask(sdk.WrapSDKContext(ctx), req)
	require.Error(t, err)
}

func TestGetTasksByJob(t *testing.T) {
	qs, ctx := setupQueryServer(t)

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

	task3 := types.Task{
		ID:            "task-3",
		JobID:         "job-2",
		ShardID:       "shard-3",
		NodeID:        "node-1",
		Status:        types.TaskStatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Progress:      0.0,
		CheckpointCID: "",
	}

	qs.Keeper.SetTask(ctx, task1)
	qs.Keeper.SetTask(ctx, task2)
	qs.Keeper.SetTask(ctx, task3)

	req := &types.QueryGetTasksByJobRequest{JobId: "job-1"}
	resp, err := qs.GetTasksByJob(sdk.WrapSDKContext(ctx), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Tasks, 2)

	req.JobId = "job-2"
	resp, err = qs.GetTasksByJob(sdk.WrapSDKContext(ctx), req)
	require.NoError(t, err)
	require.Len(t, resp.Tasks, 1)

	req.JobId = ""
	_, err = qs.GetTasksByJob(sdk.WrapSDKContext(ctx), req)
	require.Error(t, err)
}

