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

	computekeeper "github.com/atlas/chain/x/compute/keeper"
	computetypes "github.com/atlas/chain/x/compute/types"
	storagekeeper "github.com/atlas/chain/x/storage/keeper"
	"github.com/atlas/chain/x/training/types"
)

func setupMsgServer(t *testing.T) (MsgServer, sdk.Context) {
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
	ms := NewMsgServer(keeper)

	ctx := sdk.NewContext(stateStore, tmproto.Header{Time: time.Now()}, false, log.NewNopLogger())

	return ms, ctx
}

func TestSubmitJob(t *testing.T) {
	ms, ctx := setupMsgServer(t)

	msg := &types.MsgSubmitJob{
		Creator:    "cosmos1abc123",
		ModelId:    "model-1",
		DatasetCid: "QmABC123",
		Config:     map[string]string{"epochs": "10"},
	}

	resp, err := ms.SubmitJob(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.JobId)

	job, found := ms.Keeper.GetJob(ctx, resp.JobId)
	require.True(t, found)
	require.Equal(t, msg.ModelId, job.ModelID)
	require.Equal(t, msg.DatasetCid, job.DatasetCID)

	_, err = ms.SubmitJob(context.Background(), nil)
	require.Error(t, err)
}

func TestCreateTask(t *testing.T) {
	ms, ctx := setupMsgServer(t)

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

	ms.Keeper.SetJob(ctx, job)

	msg := &types.MsgCreateTask{
		Creator: "cosmos1abc123",
		JobId:   "job-1",
		ShardId: "shard-1",
		NodeId:  "node-1",
	}

	resp, err := ms.CreateTask(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.TaskId)

	task, found := ms.Keeper.GetTask(ctx, resp.TaskId)
	require.True(t, found)
	require.Equal(t, msg.JobId, task.JobID)
	require.Equal(t, msg.ShardId, task.ShardID)

	updatedJob, _ := ms.Keeper.GetJob(ctx, "job-1")
	require.Contains(t, updatedJob.Tasks, resp.TaskId)

	msg.JobId = "nonexistent"
	_, err = ms.CreateTask(sdk.WrapSDKContext(ctx), msg)
	require.Error(t, err)

	_, err = ms.CreateTask(context.Background(), nil)
	require.Error(t, err)
}

func TestUpdateTaskStatus(t *testing.T) {
	ms, ctx := setupMsgServer(t)

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

	ms.Keeper.SetTask(ctx, task)

	msg := &types.MsgUpdateTaskStatus{
		Creator:       "cosmos1abc123",
		TaskId:        "task-1",
		Status:        string(types.TaskStatusInProgress),
		Progress:      0.5,
		CheckpointCid: "QmCheckpoint123",
	}

	resp, err := ms.UpdateTaskStatus(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
	require.NotNil(t, resp)

	updatedTask, found := ms.Keeper.GetTask(ctx, "task-1")
	require.True(t, found)
	require.Equal(t, types.TaskStatusInProgress, updatedTask.Status)
	require.Equal(t, 0.5, updatedTask.Progress)
	require.Equal(t, "QmCheckpoint123", updatedTask.CheckpointCID)

	msg.TaskId = "nonexistent"
	_, err = ms.UpdateTaskStatus(sdk.WrapSDKContext(ctx), msg)
	require.Error(t, err)

	_, err = ms.UpdateTaskStatus(context.Background(), nil)
	require.Error(t, err)
}

