package recovery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/atlas/node/executor"
	"github.com/stretchr/testify/require"
)

func TestCleanupTaskState(t *testing.T) {
	tempDir := t.TempDir()
	taskID := "test-task"
	taskDir := filepath.Join(tempDir, taskID)
	
	err := os.MkdirAll(taskDir, 0755)
	require.NoError(t, err)
	
	testFile := filepath.Join(taskDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test"), 0644)
	require.NoError(t, err)
	
	err = CleanupTaskState(taskID, tempDir)
	require.NoError(t, err)
	
	_, err = os.Stat(taskDir)
	require.True(t, os.IsNotExist(err))
}

func TestHandleRollback(t *testing.T) {
	tempDir := t.TempDir()
	ipfsAPIURL := "/ip4/127.0.0.1/tcp/5001"
	
	exec := executor.NewExecutor(nil)
	exec.SetWorkDir(tempDir)
	
	task := &executor.Task{
		ID:      "test-task",
		JobID:   "test-job",
		ShardID: "test-shard",
		NodeID:  "test-node",
		Status:  "in_progress",
	}
	
	err := exec.AddTask(task)
	require.NoError(t, err)
	
	taskDir := filepath.Join(tempDir, task.ID)
	err = os.MkdirAll(taskDir, 0755)
	require.NoError(t, err)
	
	err = HandleRollback(exec, task.ID, tempDir, ipfsAPIURL)
	require.NoError(t, err)
}

