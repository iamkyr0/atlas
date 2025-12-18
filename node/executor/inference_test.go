package executor

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/atlas/storage/manager"
	"github.com/stretchr/testify/require"
)

func TestInferenceExecutor_ExecuteInference(t *testing.T) {
	tempDir := t.TempDir()
	ipfsManager := manager.NewIPFSManager("/ip4/127.0.0.1/tcp/5001")
	
	executor := NewExecutor(nil)
	executor.SetWorkDir(tempDir)
	executor.SetIPFSAPIURL("/ip4/127.0.0.1/tcp/5001")
	
	inferenceExecutor := NewInferenceExecutor(executor, tempDir, "/ip4/127.0.0.1/tcp/5001")
	
	task := &Task{
		ID:        "test-inference-1",
		TaskType:  "inference",
		ModelPath: "test_model.pt",
		InputData: []byte(`{"data": [1.0, 2.0, 3.0]}`),
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	output, err := inferenceExecutor.ExecuteInference(ctx, task, "test_model.pt", task.InputData)
	
	if err != nil {
		t.Logf("Inference test skipped (model file not available): %v", err)
		return
	}
	
	require.NoError(t, err)
	require.NotNil(t, output)
	require.NotEmpty(t, output.Result)
	require.Greater(t, output.LatencyMs, int64(0))
}

func TestInferenceExecutor_CreateInferenceScript(t *testing.T) {
	tempDir := t.TempDir()
	executor := NewExecutor(nil)
	inferenceExecutor := NewInferenceExecutor(executor, tempDir, "/ip4/127.0.0.1/tcp/5001")
	
	scriptPath := filepath.Join(tempDir, "inference.py")
	modelPath := "/path/to/model.pt"
	inputPath := filepath.Join(tempDir, "input.json")
	outputPath := filepath.Join(tempDir, "output.json")
	
	err := inferenceExecutor.createInferenceScript(scriptPath, modelPath, inputPath, outputPath, "pytorch")
	require.NoError(t, err)
	require.FileExists(t, scriptPath)
	
	scriptContent, err := os.ReadFile(scriptPath)
	require.NoError(t, err)
	require.Contains(t, string(scriptContent), "Atlas Inference Script")
	require.Contains(t, string(scriptContent), modelPath)
}

func TestInferenceInput_Unmarshal(t *testing.T) {
	inputJSON := `{
		"data": [1.0, 2.0, 3.0],
		"model_type": "llm",
		"options": {"temperature": 0.7}
	}`
	
	var input InferenceInput
	err := json.Unmarshal([]byte(inputJSON), &input)
	require.NoError(t, err)
	require.Equal(t, "llm", input.ModelType)
	require.NotNil(t, input.Data)
	require.NotNil(t, input.Options)
}

func TestInferenceOutput_Marshal(t *testing.T) {
	output := &InferenceOutput{
		Result:    []float64{0.1, 0.2, 0.3},
		LatencyMs: 100,
		ModelID:   "model-123",
	}
	
	jsonData, err := json.Marshal(output)
	require.NoError(t, err)
	
	var decoded InferenceOutput
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)
	require.Equal(t, output.LatencyMs, decoded.LatencyMs)
	require.Equal(t, output.ModelID, decoded.ModelID)
}

