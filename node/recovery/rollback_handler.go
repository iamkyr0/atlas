package recovery

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/atlas/node/executor"
	"github.com/atlas/storage/pubsub"
)

func HandleRollback(exec *executor.Executor, taskID string, workDir string, ipfsAPIURL string) error {
	if err := exec.StopTask(taskID); err != nil {
		return fmt.Errorf("failed to stop task: %w", err)
	}
	
	if err := CleanupTaskState(taskID, workDir); err != nil {
		return fmt.Errorf("failed to cleanup: %w", err)
	}

	task, found := exec.GetTask(taskID)
	if !found {
		return fmt.Errorf("task not found: %s", taskID)
	}

	event := map[string]interface{}{
		"event_type": "task_rollback",
		"task_id":    taskID,
		"job_id":     task.JobID,
		"shard_id":   task.ShardID,
		"node_id":    task.NodeID,
		"timestamp":  time.Now().Unix(),
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	topic := fmt.Sprintf("/atlas/recovery/rollback/%s", task.JobID)
	if err := pubsub.PublishEvent(ipfsAPIURL, topic, eventData); err != nil {
		return fmt.Errorf("failed to publish rollback event: %w", err)
	}
	
	return nil
}

func CleanupTaskState(taskID string, workDir string) error {
	taskDir := filepath.Join(workDir, taskID)
	
	if err := os.RemoveAll(taskDir); err != nil {
		return fmt.Errorf("failed to remove task directory: %w", err)
	}
	
	return nil
}

