package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/atlas/chain/x/training/types"
)

type QueryServer struct {
	Keeper
}

func NewQueryServer(keeper Keeper) QueryServer {
	return QueryServer{Keeper: keeper}
}

func (qs QueryServer) GetJob(ctx context.Context, req *types.QueryGetJobRequest) (*types.QueryGetJobResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.JobId == "" {
		return nil, status.Error(codes.InvalidArgument, "job_id cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	job, found := qs.Keeper.GetJob(sdkCtx, req.JobId)
	if !found {
		return nil, status.Error(codes.NotFound, "job not found")
	}

	return &types.QueryGetJobResponse{Job: &job}, nil
}

func (qs QueryServer) ListJobs(ctx context.Context, req *types.QueryListJobsRequest) (*types.QueryListJobsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	var jobs []types.Job

	store := sdkCtx.KVStore(qs.Keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte("job:"))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var job types.Job
		qs.Keeper.cdc.MustUnmarshal(iterator.Value(), &job)
		jobs = append(jobs, job)
	}

	return &types.QueryListJobsResponse{Jobs: jobs}, nil
}

func (qs QueryServer) GetTask(ctx context.Context, req *types.QueryGetTaskRequest) (*types.QueryGetTaskResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.TaskId == "" {
		return nil, status.Error(codes.InvalidArgument, "task_id cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	task, found := qs.Keeper.GetTask(sdkCtx, req.TaskId)
	if !found {
		return nil, status.Error(codes.NotFound, "task not found")
	}

	return &types.QueryGetTaskResponse{Task: &task}, nil
}

func (qs QueryServer) GetTasksByJob(ctx context.Context, req *types.QueryGetTasksByJobRequest) (*types.QueryGetTasksByJobResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.JobId == "" {
		return nil, status.Error(codes.InvalidArgument, "job_id cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	var tasks []types.Task

	qs.Keeper.IterateTasks(sdkCtx, func(task types.Task) (stop bool) {
		if task.JobID == req.JobId {
			tasks = append(tasks, task)
		}
		return false
	})

	return &types.QueryGetTasksByJobResponse{Tasks: tasks}, nil
}

