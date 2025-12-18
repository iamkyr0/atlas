package recovery

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
	"github.com/atlas/storage/manager"
)

type Checkpoint struct {
	TaskID      string
	Epoch       int
	Iteration   int
	CID         string
	Timestamp   time.Time
	Signature   string
}

func SaveCheckpoint(ipfsManager *manager.IPFSManager, taskID string, epoch int, iteration int, modelPath string) (*Checkpoint, error) {
	checkpointDir := filepath.Join("/tmp/checkpoints", taskID, fmt.Sprintf("epoch_%d_iter_%d", epoch, iteration))
	if err := os.MkdirAll(checkpointDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create checkpoint directory: %w", err)
	}

	checkpointModelPath := filepath.Join(checkpointDir, "model.pt")
	if err := copyFile(modelPath, checkpointModelPath); err != nil {
		return nil, fmt.Errorf("failed to copy model: %w", err)
	}

	metadata := map[string]interface{}{
		"task_id":   taskID,
		"epoch":     epoch,
		"iteration": iteration,
		"timestamp": time.Now().Unix(),
	}
	metadataPath := filepath.Join(checkpointDir, "metadata.json")
	metadataJSON, _ := json.Marshal(metadata)
	if err := os.WriteFile(metadataPath, metadataJSON, 0644); err != nil {
		return nil, fmt.Errorf("failed to write metadata: %w", err)
	}

	cid, err := ipfsManager.AddFile(checkpointDir)
	if err != nil {
		return nil, fmt.Errorf("failed to upload checkpoint: %w", err)
	}

	signature := calculateSignature(taskID, epoch, iteration, cid)

	checkpoint := &Checkpoint{
		TaskID:    taskID,
		Epoch:     epoch,
		Iteration: iteration,
		CID:       cid,
		Timestamp: time.Now(),
		Signature: signature,
	}

	return checkpoint, nil
}

func LoadCheckpoint(ipfsManager *manager.IPFSManager, cid string, outputPath string) error {
	if err := ipfsManager.GetFile(cid, outputPath); err != nil {
		return fmt.Errorf("failed to download checkpoint: %w", err)
	}
	return nil
}

func ValidateCheckpoint(cp *Checkpoint) error {
	expectedSignature := calculateSignature(cp.TaskID, cp.Epoch, cp.Iteration, cp.CID)
	if cp.Signature != expectedSignature {
		return fmt.Errorf("invalid checkpoint signature")
	}

	age := time.Since(cp.Timestamp)
	if age > 7*24*time.Hour {
		return fmt.Errorf("checkpoint too old")
	}

	return nil
}

func calculateSignature(taskID string, epoch int, iteration int, cid string) string {
	data := fmt.Sprintf("%s:%d:%d:%s", taskID, epoch, iteration, cid)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
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

