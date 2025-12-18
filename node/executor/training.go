package executor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	
	"github.com/atlas/storage/manager"
)

type TrainingExecutor struct {
	executor    *Executor
	workDir     string
	ipfsManager *manager.IPFSManager
}

func NewTrainingExecutor(e *Executor, workDir string, ipfsAPIURL string) *TrainingExecutor {
	return &TrainingExecutor{
		executor:    e,
		workDir:     workDir,
		ipfsManager: manager.NewIPFSManager(ipfsAPIURL),
	}
}

func (te *TrainingExecutor) ExecuteTraining(ctx context.Context, task *Task, modelPath string, datasetPath string) error {
	// Create working directory for task
	taskDir := filepath.Join(te.workDir, task.ID)
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		return fmt.Errorf("failed to create task directory: %w", err)
	}

	// Download model and dataset if needed
	if err := te.downloadIfNeeded(ctx, modelPath, taskDir); err != nil {
		return err
	}
	if err := te.downloadIfNeeded(ctx, datasetPath, taskDir); err != nil {
		return err
	}

	// Execute training script
	scriptPath := filepath.Join(taskDir, "train.py")
	if err := te.createTrainingScript(scriptPath, modelPath, datasetPath); err != nil {
		return err
	}

	// Try python3 first, fallback to python
	cmd := exec.CommandContext(ctx, "python3", scriptPath)
	cmd.Dir = taskDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Fallback to python
		cmd = exec.CommandContext(ctx, "python", scriptPath)
		cmd.Dir = taskDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("training failed: %w", err)
		}
	}

	return nil
}

func (te *TrainingExecutor) downloadIfNeeded(ctx context.Context, path string, destDir string) error {
	// Check if path is IPFS CID (starts with Qm for v0 or ba for v1)
	if len(path) >= 2 && (path[:2] == "Qm" || path[:2] == "ba") {
		// Download from IPFS
		destPath := filepath.Join(destDir, "downloaded_file")
		if err := te.ipfsManager.GetFile(path, destPath); err != nil {
			return fmt.Errorf("IPFS download failed: %w", err)
		}
		return nil
	}
	
	// Assume it's a local path - verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("local file not found: %s", path)
	}
	
	return nil
}

func (te *TrainingExecutor) createTrainingScript(scriptPath string, modelPath string, datasetPath string) error {
	script := fmt.Sprintf(`
import torch
import torch.nn as nn
from torch.utils.data import DataLoader

# Load model
model = torch.load('%s')

# Load dataset
# dataset = load_dataset('%s')

# Training loop
optimizer = torch.optim.Adam(model.parameters())
criterion = nn.CrossEntropyLoss()

for epoch in range(10):
    # Training code here
    pass

# Save checkpoint
torch.save(model.state_dict(), 'checkpoint.pt')
`, modelPath, datasetPath)

	return os.WriteFile(scriptPath, []byte(script), 0755)
}

