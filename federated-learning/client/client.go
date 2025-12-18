package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	
	"github.com/atlas/federated-learning/protocols"
	"github.com/atlas/storage/manager"
)

type FLClient struct {
	nodeID      string
	protocol    *protocols.FLProtocol
	ipfsManager *manager.IPFSManager
	workDir     string
	keepWorkDir bool // If true, don't cleanup training directory after training
}

func NewFLClient(nodeID string, ipfsAPIURL string, workDir string) *FLClient {
	return &FLClient{
		nodeID:      nodeID,
		protocol:    protocols.NewFLProtocol(ipfsAPIURL, nodeID),
		ipfsManager: manager.NewIPFSManager(ipfsAPIURL),
		workDir:     workDir,
		keepWorkDir: false, // Default: cleanup after training
	}
}

func (c *FLClient) SetKeepWorkDir(keep bool) {
	c.keepWorkDir = keep
}

func (c *FLClient) GetTrainScriptPath(shardID string) string {
	shardPrefix := shardID
	if len(shardID) > 8 {
		shardPrefix = shardID[:8]
	}
	trainDir := filepath.Join(c.workDir, fmt.Sprintf("train_%s_%s_%d", c.nodeID, shardPrefix, time.Now().Unix()))
	return filepath.Join(trainDir, "train.py")
}

func (c *FLClient) Train(ctx context.Context, shardID string, modelPath string) ([]float64, error) {
	trainDir := filepath.Join(c.workDir, fmt.Sprintf("train_%s_%d", c.nodeID, time.Now().Unix()))
	if err := os.MkdirAll(trainDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create training directory: %w", err)
	}
	
	if !c.keepWorkDir {
		defer os.RemoveAll(trainDir)
	} else {
		fmt.Printf("Training directory preserved: %s\n", trainDir)
		fmt.Printf("train.py location: %s\n", filepath.Join(trainDir, "train.py"))
	}

	modelLocalPath := filepath.Join(trainDir, "model")
	if err := c.downloadFromIPFS(ctx, modelPath, modelLocalPath); err != nil {
		return nil, fmt.Errorf("failed to download model: %w", err)
	}

	shardLocalPath := filepath.Join(trainDir, "shard")
	if err := c.downloadFromIPFS(ctx, shardID, shardLocalPath); err != nil {
		return nil, fmt.Errorf("failed to download shard: %w", err)
	}

	gradients, err := c.performTraining(ctx, modelLocalPath, shardLocalPath, trainDir)
	if err != nil {
		return nil, fmt.Errorf("training failed: %w", err)
	}

	return gradients, nil
}

func (c *FLClient) downloadFromIPFS(ctx context.Context, cidOrPath string, destPath string) error {
	if len(cidOrPath) >= 2 && (cidOrPath[:2] == "Qm" || cidOrPath[:2] == "ba") {
		if err := c.ipfsManager.GetFile(cidOrPath, destPath); err != nil {
			return fmt.Errorf("IPFS download failed: %w", err)
		}
		return nil
	}
	
	// Assume it's a local path - copy file
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}
	
	srcFile, err := os.Open(cidOrPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	if _, err := destFile.ReadFrom(srcFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

func (c *FLClient) performTraining(ctx context.Context, modelPath string, shardPath string, workDir string) ([]float64, error) {
	scriptPath := filepath.Join(workDir, "train.py")
	if err := c.createTrainingScript(scriptPath, modelPath, shardPath); err != nil {
		return nil, fmt.Errorf("failed to create training script: %w", err)
	}
	
	gradientsPath := filepath.Join(workDir, "gradients.json")
	
	cmd := exec.CommandContext(ctx, "python3", scriptPath)
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		cmd = exec.CommandContext(ctx, "python", scriptPath)
		cmd.Dir = workDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("training script execution failed: %w", err)
		}
	}
	
	gradients, err := c.readGradientsFromFile(gradientsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read gradients: %w", err)
	}

	return gradients, nil
}

func (c *FLClient) readGradientsFromFile(filePath string) ([]float64, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read gradients file: %w", err)
	}
	
	var gradients []float64
	if err := json.Unmarshal(data, &gradients); err != nil {
		return nil, fmt.Errorf("failed to parse gradients JSON: %w", err)
	}
	
	if len(gradients) == 0 {
		return nil, fmt.Errorf("gradients array is empty")
	}
	
	return gradients, nil
}

func (te *TrainingExecutor) createTrainingScript(scriptPath string, modelPath string, datasetPath string) error {
	epochs := 10
	batchSize := 32
	learningRate := 0.001
	
	script := fmt.Sprintf(`#!/usr/bin/env python3
import json
import os
import sys
import torch
import torch.nn as nn
import torch.optim as optim
from torch.utils.data import Dataset, DataLoader, TensorDataset
import numpy as np

def load_model(model_path):
    """Load PyTorch model from file"""
    try:
        if os.path.isdir(model_path):
            model_files = [f for f in os.listdir(model_path) if f.endswith(('.pt', '.pth'))]
            if model_files:
                model_path = os.path.join(model_path, model_files[0])
            else:
                raise FileNotFoundError(f"No model file found in directory: {model_path}")
        
        if not os.path.exists(model_path):
            raise FileNotFoundError(f"Model file not found: {model_path}")
        
        device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
        print(f"Loading model from: {model_path}", file=sys.stderr)
        
        if model_path.endswith('.pt') or model_path.endswith('.pth'):
            model_data = torch.load(model_path, map_location=device)
            
            if isinstance(model_data, nn.Module):
                model = model_data.to(device)
                model.train()
                return model, device
            elif isinstance(model_data, dict):
                if 'model' in model_data:
                    model = model_data['model'].to(device)
                    model.train()
                    return model, device
                else:
                    raise ValueError("Invalid model format in dict")
            else:
                raise ValueError(f"Unsupported model type: {type(model_data)}")
        else:
            raise ValueError(f"Unsupported model format: {model_path}")
    except Exception as e:
        print(f"Error loading model: {e}", file=sys.stderr)
        import traceback
        traceback.print_exc(file=sys.stderr)
        raise

def load_dataset(dataset_path):
    """Load dataset from file"""
    try:
        if os.path.isdir(dataset_path):
            json_path = os.path.join(dataset_path, "data.json")
            pkl_path = os.path.join(dataset_path, "data.pkl")
            pt_path = os.path.join(dataset_path, "data.pt")
            
            if os.path.exists(json_path):
                dataset_path = json_path
            elif os.path.exists(pkl_path):
                dataset_path = pkl_path
            elif os.path.exists(pt_path):
                dataset_path = pt_path
            else:
                files = [f for f in os.listdir(dataset_path) if not f.startswith('.')]
                if files:
                    dataset_path = os.path.join(dataset_path, files[0])
                else:
                    raise FileNotFoundError(f"No data file found in directory: {dataset_path}")
        
        if not os.path.exists(dataset_path):
            raise FileNotFoundError(f"Dataset file not found: {dataset_path}")
        
        print(f"Loading dataset from: {dataset_path}", file=sys.stderr)
        
        if dataset_path.endswith('.json'):
            with open(dataset_path, 'r') as f:
                data = json.load(f)
            
            if isinstance(data, dict):
                if 'inputs' in data and 'targets' in data:
                    inputs = torch.tensor(data['inputs'], dtype=torch.float32)
                    targets = torch.tensor(data['targets'], dtype=torch.long)
                elif 'data' in data:
                    inputs = torch.tensor(data['data'], dtype=torch.float32)
                    targets = torch.tensor(data.get('labels', [0] * len(data['data'])), dtype=torch.long)
                else:
                    raise ValueError("Invalid JSON format")
            elif isinstance(data, list):
                if len(data) > 0 and isinstance(data[0], dict):
                    inputs = torch.tensor([item.get('input', item.get('x', item.get('features', [0]))) for item in data], dtype=torch.float32)
                    targets = torch.tensor([item.get('target', item.get('y', item.get('label', 0))) for item in data], dtype=torch.long)
                else:
                    inputs = torch.tensor(data, dtype=torch.float32)
                    targets = torch.zeros(len(data), dtype=torch.long)
            else:
                raise ValueError("Invalid data format")
        
        elif dataset_path.endswith('.pkl') or dataset_path.endswith('.pickle'):
            import pickle
            with open(dataset_path, 'rb') as f:
                data = pickle.load(f)
            
            if isinstance(data, tuple) and len(data) == 2:
                inputs, targets = data
                inputs = torch.tensor(inputs, dtype=torch.float32) if not isinstance(inputs, torch.Tensor) else inputs.float()
                targets = torch.tensor(targets, dtype=torch.long) if not isinstance(targets, torch.Tensor) else targets.long()
            elif isinstance(data, dict):
                inputs = torch.tensor(data.get('inputs', data.get('x', data.get('data', []))), dtype=torch.float32)
                targets = torch.tensor(data.get('targets', data.get('y', data.get('labels', []))), dtype=torch.long)
            else:
                raise ValueError("Invalid pickle format")
        
        elif dataset_path.endswith('.pt') or dataset_path.endswith('.pth'):
            data = torch.load(dataset_path)
            if isinstance(data, tuple) and len(data) == 2:
                inputs, targets = data
            elif isinstance(data, dict):
                inputs = data.get('inputs', data.get('x', data.get('data')))
                targets = data.get('targets', data.get('y', data.get('labels')))
            else:
                raise ValueError("Invalid torch format")
            
            if not isinstance(inputs, torch.Tensor):
                inputs = torch.tensor(inputs, dtype=torch.float32)
            if not isinstance(targets, torch.Tensor):
                targets = torch.tensor(targets, dtype=torch.long)
        
        else:
            raise ValueError(f"Unsupported dataset format: {dataset_path}")
        
        if inputs.dim() == 1:
            inputs = inputs.unsqueeze(1)
        
        if targets.dim() == 0:
            targets = targets.unsqueeze(0)
        
        if len(inputs) == 0:
            raise ValueError("Empty dataset")
        
        print(f"Dataset loaded: {inputs.shape[0]} samples", file=sys.stderr)
        return inputs, targets
    
    except Exception as e:
        print(f"Error loading dataset: {e}", file=sys.stderr)
        import traceback
        traceback.print_exc(file=sys.stderr)
        raise

def train_model(model, inputs, targets, epochs=%d, batch_size=%d, learning_rate=%.6f):
    """Train model"""
    device = next(model.parameters()).device
    
    inputs = inputs.to(device)
    targets = targets.to(device)
    
    dataset = TensorDataset(inputs, targets)
    dataloader = DataLoader(dataset, batch_size=batch_size, shuffle=True)
    
    is_classification = len(targets.shape) == 1 or (len(targets.shape) == 2 and targets.shape[1] == 1)
    
    if is_classification:
        num_classes = int(targets.max().item()) + 1 if len(targets.shape) == 1 else 1
        if num_classes > 1:
            criterion = nn.CrossEntropyLoss()
        else:
            criterion = nn.BCEWithLogitsLoss()
    else:
        criterion = nn.MSELoss()
    
    optimizer = optim.Adam(model.parameters(), lr=learning_rate)
    
    model.train()
    
    for epoch in range(epochs):
        epoch_loss = 0.0
        batch_count = 0
        
        for batch_idx, (batch_inputs, batch_targets) in enumerate(dataloader):
            try:
                optimizer.zero_grad()
                
                output = model(batch_inputs)
                
                if isinstance(output, tuple):
                    output = output[0]
                
                if output.shape != batch_targets.shape:
                    if len(batch_targets.shape) == 1 and len(output.shape) == 2:
                        if output.shape[1] == 1:
                            output = output.squeeze(1)
                            loss = criterion(output, batch_targets.float())
                        else:
                            batch_targets = batch_targets.long()
                            if num_classes > 1:
                                loss = criterion(output, batch_targets)
                            else:
                                loss = criterion(output.squeeze(), batch_targets.float())
                    else:
                        loss = criterion(output, batch_targets.float())
                else:
                    if is_classification and num_classes > 1:
                        loss = criterion(output, batch_targets.long())
                    else:
                        loss = criterion(output, batch_targets.float())
                
                loss.backward()
                torch.nn.utils.clip_grad_norm_(model.parameters(), max_norm=1.0)
                optimizer.step()
                
                epoch_loss += loss.item()
                batch_count += 1
            
            except Exception as e:
                print(f"Error in training batch {batch_idx}: {e}", file=sys.stderr)
                continue
        
        if batch_count > 0:
            avg_loss = epoch_loss / batch_count
            print(f"Epoch {epoch+1}/{epochs}, Loss: {avg_loss:.6f}", file=sys.stderr)
    
    checkpoint_path = 'checkpoint.pt'
    torch.save(model.state_dict(), checkpoint_path)
    print(f"Checkpoint saved to {checkpoint_path}", file=sys.stderr)

def main():
    model_path = '%s'
    dataset_path = '%s'
    
    try:
        print("=" * 50, file=sys.stderr)
        print("Atlas Node Training Script", file=sys.stderr)
        print("=" * 50, file=sys.stderr)
        
        model, device = load_model(model_path)
        print(f"Model loaded on {device}", file=sys.stderr)
        
        inputs, targets = load_dataset(dataset_path)
        
        print("Starting training...", file=sys.stderr)
        train_model(model, inputs, targets, epochs=%d, batch_size=%d, learning_rate=%.6f)
        print("Training completed", file=sys.stderr)
        sys.exit(0)
    
    except Exception as e:
        print(f"Training failed: {e}", file=sys.stderr)
        import traceback
        traceback.print_exc(file=sys.stderr)
        sys.exit(1)

if __name__ == '__main__':
    main()
`, epochs, batchSize, learningRate, modelPath, datasetPath, epochs, batchSize, learningRate)

	return os.WriteFile(scriptPath, []byte(script), 0755)
}
func (c *FLClient) SendGradients(ctx context.Context, jobID string, round int, gradients []float64) error {
	return c.protocol.SendGradients(ctx, jobID, round, gradients)
}

func (c *FLClient) ReceiveModel(ctx context.Context, jobID string, handler func([]float64)) error {
	return c.protocol.ReceiveModel(ctx, jobID, handler)
}

