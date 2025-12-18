package types

import (
	"time"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusAssigned   TaskStatus = "assigned"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusPaused     TaskStatus = "paused"
	TaskStatusRollback   TaskStatus = "rollback"
	TaskStatusDelegated  TaskStatus = "delegated"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
)

type Task struct {
	ID           string     `json:"id"`
	JobID        string     `json:"job_id"`
	ShardID      string     `json:"shard_id"`
	NodeID       string     `json:"node_id"`
	Status       TaskStatus `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Progress     float64    `json:"progress"`
	CheckpointCID string    `json:"checkpoint_cid"`
}

type Job struct {
	ID          string    `json:"id"`
	ModelID     string    `json:"model_id"`
	DatasetCID  string    `json:"dataset_cid"`
	Config      map[string]interface{} `json:"config"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Progress    float64    `json:"progress"`
	Tasks       []string   `json:"tasks"`
}

