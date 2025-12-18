package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/atlas/chain/x/training/types"
)

type MsgServer struct {
	Keeper
}

func NewMsgServer(keeper Keeper) MsgServer {
	return MsgServer{Keeper: keeper}
}

func (ms MsgServer) SubmitJob(ctx context.Context, msg *types.MsgSubmitJob) (*types.MsgSubmitJobResponse, error) {
	if msg == nil {
		return nil, fmt.Errorf("invalid message")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	jobID := fmt.Sprintf("job-%d", sdkCtx.BlockTime().UnixNano())

	job := types.Job{
		ID:         jobID,
		ModelID:    msg.ModelId,
		DatasetCID: msg.DatasetCid,
		Config:     msg.Config,
		Status:     types.TaskStatusPending,
		CreatedAt:  sdkCtx.BlockTime(),
		UpdatedAt:  sdkCtx.BlockTime(),
		Progress:   0.0,
		Tasks:      []string{},
	}

	ms.Keeper.SetJob(sdkCtx, job)

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeJobCreated,
			sdk.NewAttribute(types.AttributeKeyJobID, jobID),
			sdk.NewAttribute(types.AttributeKeyModelID, msg.ModelId),
			sdk.NewAttribute(types.AttributeKeyDatasetCID, msg.DatasetCid),
		),
	)

	return &types.MsgSubmitJobResponse{JobId: jobID}, nil
}

func (ms MsgServer) CreateTask(ctx context.Context, msg *types.MsgCreateTask) (*types.MsgCreateTaskResponse, error) {
	if msg == nil {
		return nil, fmt.Errorf("invalid message")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	_, found := ms.Keeper.GetJob(sdkCtx, msg.JobId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrJobNotFound, "job %s not found", msg.JobId)
	}

	taskID := fmt.Sprintf("task-%d", sdkCtx.BlockTime().UnixNano())

	task := types.Task{
		ID:            taskID,
		JobID:         msg.JobId,
		ShardID:       msg.ShardId,
		NodeID:        msg.NodeId,
		Status:        types.TaskStatusPending,
		CreatedAt:     sdkCtx.BlockTime(),
		UpdatedAt:     sdkCtx.BlockTime(),
		Progress:      0.0,
		CheckpointCID: "",
	}

	ms.Keeper.SetTask(sdkCtx, task)

	job, _ := ms.Keeper.GetJob(sdkCtx, msg.JobId)
	job.Tasks = append(job.Tasks, taskID)
	ms.Keeper.SetJob(sdkCtx, job)

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTaskCreated,
			sdk.NewAttribute(types.AttributeKeyTaskID, taskID),
			sdk.NewAttribute(types.AttributeKeyJobID, msg.JobId),
			sdk.NewAttribute(types.AttributeKeyShardID, msg.ShardId),
		),
	)

	return &types.MsgCreateTaskResponse{TaskId: taskID}, nil
}

func (ms MsgServer) UpdateTaskStatus(ctx context.Context, msg *types.MsgUpdateTaskStatus) (*types.MsgUpdateTaskStatusResponse, error) {
	if msg == nil {
		return nil, fmt.Errorf("invalid message")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	task, found := ms.Keeper.GetTask(sdkCtx, msg.TaskId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrTaskNotFound, "task %s not found", msg.TaskId)
	}

	task.Status = types.TaskStatus(msg.Status)
	task.UpdatedAt = sdkCtx.BlockTime()
	if msg.Progress >= 0 {
		task.Progress = msg.Progress
	}
	if msg.CheckpointCid != "" {
		task.CheckpointCID = msg.CheckpointCid
	}

	ms.Keeper.SetTask(sdkCtx, task)

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTaskStatusUpdated,
			sdk.NewAttribute(types.AttributeKeyTaskID, msg.TaskId),
			sdk.NewAttribute(types.AttributeKeyStatus, msg.Status),
		),
	)

	return &types.MsgUpdateTaskStatusResponse{}, nil
}

