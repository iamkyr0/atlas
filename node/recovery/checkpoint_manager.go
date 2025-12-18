package recovery

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
	"github.com/atlas/storage/manager"
)

type CheckpointManager struct {
	ipfsManager *manager.IPFSManager
	checkpointDir string
}

func NewCheckpointManager(ipfsAPIURL string, checkpointDir string) *CheckpointManager {
	return &CheckpointManager{
		ipfsManager: manager.NewIPFSManager(ipfsAPIURL),
		checkpointDir: checkpointDir,
	}
}

func (cm *CheckpointManager) SaveCheckpoint(ctx context.Context, taskID string, epoch int, iteration int, modelPath string) (*Checkpoint, error) {
	checkpointPath := filepath.Join(cm.checkpointDir, taskID, fmt.Sprintf("epoch_%d_iter_%d", epoch, iteration))
	if err := os.MkdirAll(checkpointPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create checkpoint directory: %w", err)
	}

	checkpointModelPath := filepath.Join(checkpointPath, "model.pt")
	if err := copyFile(modelPath, checkpointModelPath); err != nil {
		return nil, fmt.Errorf("failed to copy model: %w", err)
	}

	metadata := map[string]interface{}{
		"task_id":   taskID,
		"epoch":     epoch,
		"iteration": iteration,
		"timestamp": time.Now().Unix(),
	}
	metadataPath := filepath.Join(checkpointPath, "metadata.json")
	metadataJSON, _ := json.Marshal(metadata)
	if err := os.WriteFile(metadataPath, metadataJSON, 0644); err != nil {
		return nil, fmt.Errorf("failed to write metadata: %w", err)
	}

	cid, err := cm.ipfsManager.AddFile(checkpointPath)
	if err != nil {
		return nil, fmt.Errorf("failed to upload checkpoint: %w", err)
	}

	checkpoint := &Checkpoint{
		TaskID:    taskID,
		Epoch:     epoch,
		Iteration: iteration,
		CID:       cid,
		Timestamp: time.Now(),
	}

	return checkpoint, nil
}

func (cm *CheckpointManager) LoadCheckpoint(ctx context.Context, checkpoint *Checkpoint, outputPath string) error {
	if err := cm.ipfsManager.GetFile(checkpoint.CID, outputPath); err != nil {
		return fmt.Errorf("failed to download checkpoint: %w", err)
	}

	if err := ValidateCheckpoint(checkpoint); err != nil {
		return fmt.Errorf("checkpoint validation failed: %w", err)
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

