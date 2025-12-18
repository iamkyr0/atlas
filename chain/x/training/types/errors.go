package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrJobNotFound  = sdkerrors.Register(ModuleName, 1, "job not found")
	ErrTaskNotFound = sdkerrors.Register(ModuleName, 2, "task not found")
	ErrInvalidJob   = sdkerrors.Register(ModuleName, 3, "invalid job")
	ErrInvalidTask  = sdkerrors.Register(ModuleName, 4, "invalid task")
)

const (
	EventTypeJobCreated        = "job_created"
	EventTypeTaskCreated       = "task_created"
	EventTypeTaskStatusUpdated = "task_status_updated"
	
	AttributeKeyJobID    = "job_id"
	AttributeKeyTaskID   = "task_id"
	AttributeKeyModelID  = "model_id"
	AttributeKeyShardID  = "shard_id"
	AttributeKeyStatus   = "status"
	AttributeKeyDatasetCID = "dataset_cid"
)

