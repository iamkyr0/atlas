package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/atlas/storage/manager"
)

type InferenceExecutor struct {
	executor    *Executor
	workDir     string
	ipfsManager *manager.IPFSManager
	modelCache  map[string]string
}

type InferenceInput struct {
	Data      interface{} `json:"data"`
	ModelType string      `json:"model_type,omitempty"`
	Options   map[string]interface{} `json:"options,omitempty"`
}

type InferenceOutput struct {
	Result    interface{} `json:"result"`
	LatencyMs int64       `json:"latency_ms"`
	ModelID   string      `json:"model_id"`
}

func NewInferenceExecutor(e *Executor, workDir string, ipfsAPIURL string) *InferenceExecutor {
	return &InferenceExecutor{
		executor:    e,
		workDir:     workDir,
		ipfsManager: manager.NewIPFSManager(ipfsAPIURL),
		modelCache:  make(map[string]string),
	}
}

func (ie *InferenceExecutor) ExecuteInference(ctx context.Context, task *Task, modelPath string, inputData []byte) (*InferenceOutput, error) {
	startTime := time.Now()
	
	taskDir := filepath.Join(ie.workDir, task.ID)
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create task directory: %w", err)
	}

	modelLocalPath, err := ie.ensureModelDownloaded(ctx, modelPath, taskDir)
	if err != nil {
		return nil, fmt.Errorf("failed to download model: %w", err)
	}

	inputPath := filepath.Join(taskDir, "input.json")
	if err := os.WriteFile(inputPath, inputData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write input file: %w", err)
	}

	outputPath := filepath.Join(taskDir, "output.json")
	
	var input InferenceInput
	if err := json.Unmarshal(inputData, &input); err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	modelType := input.ModelType
	if modelType == "" {
		modelType = "auto"
	}

	scriptPath := filepath.Join(taskDir, "inference.py")
	if err := ie.createInferenceScript(scriptPath, modelLocalPath, inputPath, outputPath, modelType); err != nil {
		return nil, fmt.Errorf("failed to create inference script: %w", err)
	}

	cmd := exec.CommandContext(ctx, "python3", scriptPath)
	cmd.Dir = taskDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		cmd = exec.CommandContext(ctx, "python", scriptPath)
		cmd.Dir = taskDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("inference script execution failed: %w", err)
		}
	}

	outputData, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read output file: %w", err)
	}

	var result interface{}
	if err := json.Unmarshal(outputData, &result); err != nil {
		return nil, fmt.Errorf("failed to parse output: %w", err)
	}

	latency := time.Since(startTime).Milliseconds()

	return &InferenceOutput{
		Result:    result,
		LatencyMs: latency,
		ModelID:   task.ModelPath,
	}, nil
}

func (ie *InferenceExecutor) ensureModelDownloaded(ctx context.Context, modelPath string, taskDir string) (string, error) {
	if cached, ok := ie.modelCache[modelPath]; ok {
		if _, err := os.Stat(cached); err == nil {
			return cached, nil
		}
	}

	modelDir := filepath.Join(taskDir, "model")
	if err := os.MkdirAll(modelDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create model directory: %w", err)
	}

	if len(modelPath) >= 2 && (modelPath[:2] == "Qm" || modelPath[:2] == "ba") {
		if err := ie.ipfsManager.GetFile(modelPath, modelDir); err != nil {
			return "", fmt.Errorf("IPFS download failed: %w", err)
		}
	} else {
		if _, err := os.Stat(modelPath); os.IsNotExist(err) {
			return "", fmt.Errorf("local model file not found: %s", modelPath)
		}
		if err := copyFile(modelPath, filepath.Join(modelDir, filepath.Base(modelPath))); err != nil {
			return "", fmt.Errorf("failed to copy model: %w", err)
		}
	}

	modelFiles, err := findModelFiles(modelDir)
	if err != nil {
		return "", fmt.Errorf("failed to find model files: %w", err)
	}

	if len(modelFiles) == 0 {
		return "", fmt.Errorf("no model files found in %s", modelDir)
	}

	modelLocalPath := modelFiles[0]
	ie.modelCache[modelPath] = modelLocalPath

	return modelLocalPath, nil
}

func (ie *InferenceExecutor) createInferenceScript(scriptPath string, modelPath string, inputPath string, outputPath string, modelType string) error {
	script := fmt.Sprintf(`#!/usr/bin/env python3
import json
import os
import sys
import time
import torch
import numpy as np
from pathlib import Path

def load_model(model_path):
    """Load model from file"""
    device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
    print(f"Loading model from: {model_path}", file=sys.stderr)
    
    if os.path.isdir(model_path):
        model_files = [f for f in os.listdir(model_path) if f.endswith(('.pt', '.pth', '.onnx', '.h5'))]
        if model_files:
            model_path = os.path.join(model_path, model_files[0])
        else:
            raise FileNotFoundError(f"No model file found in directory: {model_path}")
    
    extension = Path(model_path).suffix.lower()
    
    if extension in ['.pt', '.pth']:
        model_data = torch.load(model_path, map_location=device)
        if isinstance(model_data, torch.nn.Module):
            model = model_data.to(device)
            model.eval()
            return model, device, 'pytorch'
        elif isinstance(model_data, dict):
            if 'model' in model_data:
                model = model_data['model'].to(device)
                model.eval()
                return model, device, 'pytorch'
            else:
                raise ValueError("Invalid PyTorch model format")
        else:
            raise ValueError(f"Unsupported PyTorch model type: {type(model_data)}")
    
    elif extension == '.onnx':
        try:
            import onnxruntime as ort
            session = ort.InferenceSession(model_path)
            return session, device, 'onnx'
        except ImportError:
            raise ImportError("ONNX Runtime not installed. Install with: pip install onnxruntime")
    
    elif extension == '.h5':
        try:
            import tensorflow as tf
            model = tf.keras.models.load_model(model_path)
            return model, device, 'tensorflow'
        except ImportError:
            raise ImportError("TensorFlow not installed. Install with: pip install tensorflow")
    
    else:
        raise ValueError(f"Unsupported model format: {extension}")

def prepare_input(input_data, model_type, framework):
    """Prepare input data for model"""
    if framework == 'pytorch':
        if isinstance(input_data, list):
            tensor = torch.tensor(input_data, dtype=torch.float32)
        elif isinstance(input_data, dict):
            if 'input_ids' in input_data:
                tensor = torch.tensor(input_data['input_ids'], dtype=torch.long)
            else:
                tensor = torch.tensor(list(input_data.values())[0], dtype=torch.float32)
        else:
            tensor = torch.tensor(input_data, dtype=torch.float32)
        
        if tensor.dim() == 1:
            tensor = tensor.unsqueeze(0)
        
        return tensor
    
    elif framework == 'onnx':
        if isinstance(input_data, list):
            array = np.array(input_data, dtype=np.float32)
        elif isinstance(input_data, dict):
            array = np.array(list(input_data.values())[0], dtype=np.float32)
        else:
            array = np.array(input_data, dtype=np.float32)
        
        if array.ndim == 1:
            array = array.reshape(1, -1)
        
        return {list(input_data.keys())[0] if isinstance(input_data, dict) else 'input': array}
    
    elif framework == 'tensorflow':
        if isinstance(input_data, list):
            array = np.array(input_data, dtype=np.float32)
        elif isinstance(input_data, dict):
            array = np.array(list(input_data.values())[0], dtype=np.float32)
        else:
            array = np.array(input_data, dtype=np.float32)
        
        if array.ndim == 1:
            array = array.reshape(1, -1)
        
        return array
    
    return input_data

def run_inference(model, input_data, framework, device):
    """Run inference on model"""
    if framework == 'pytorch':
        model_input = prepare_input(input_data, 'pytorch', framework)
        model_input = model_input.to(device)
        
        with torch.no_grad():
            output = model(model_input)
        
        if isinstance(output, torch.Tensor):
            result = output.cpu().numpy().tolist()
        elif isinstance(output, (list, tuple)):
            result = [o.cpu().numpy().tolist() if isinstance(o, torch.Tensor) else o for o in output]
        else:
            result = output
        
        return result
    
    elif framework == 'onnx':
        model_input = prepare_input(input_data, 'onnx', framework)
        output_names = [output.name for output in model.get_outputs()]
        outputs = model.run(output_names, model_input)
        
        if len(outputs) == 1:
            result = outputs[0].tolist()
        else:
            result = [o.tolist() for o in outputs]
        
        return result
    
    elif framework == 'tensorflow':
        model_input = prepare_input(input_data, 'tensorflow', framework)
        output = model.predict(model_input, verbose=0)
        return output.tolist()
    
    return None

def main():
    model_path = '%s'
    input_path = '%s'
    output_path = '%s'
    model_type = '%s'
    
    try:
        print("=" * 50, file=sys.stderr)
        print("Atlas Inference Script", file=sys.stderr)
        print("=" * 50, file=sys.stderr)
        
        with open(input_path, 'r') as f:
            input_data = json.load(f)
        
        data = input_data.get('data', input_data)
        
        print(f"Loading model...", file=sys.stderr)
        model, device, framework = load_model(model_path)
        print(f"Model loaded on {device}, framework: {framework}", file=sys.stderr)
        
        print(f"Running inference...", file=sys.stderr)
        start_time = time.time()
        result = run_inference(model, data, framework, device)
        inference_time = (time.time() - start_time) * 1000
        
        print(f"Inference completed in {inference_time:.2f}ms", file=sys.stderr)
        
        output = {
            "result": result,
            "latency_ms": int(inference_time),
            "framework": framework
        }
        
        with open(output_path, 'w') as f:
            json.dump(output, f, indent=2)
        
        print(f"Output saved to {output_path}", file=sys.stderr)
        sys.exit(0)
    
    except Exception as e:
        print(f"Inference failed: {e}", file=sys.stderr)
        import traceback
        traceback.print_exc(file=sys.stderr)
        sys.exit(1)

if __name__ == '__main__':
    main()
`, modelPath, inputPath, outputPath, modelType)

	return os.WriteFile(scriptPath, []byte(script), 0755)
}

func findModelFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := filepath.Ext(path)
			if ext == ".pt" || ext == ".pth" || ext == ".onnx" || ext == ".h5" {
				files = append(files, path)
			}
		}
		return nil
	})
	return files, err
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

	_, err = destFile.ReadFrom(sourceFile)
	return err
}

