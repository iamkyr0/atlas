package training

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/atlas/lora/adapters"
	"github.com/stretchr/testify/require"
)

func TestNewLoRATrainer(t *testing.T) {
	adapter := adapters.NewLoRAAdapter(4, 8.0)
	trainer := NewLoRATrainer(adapter)
	
	require.NotNil(t, trainer)
	require.Equal(t, adapter, trainer.adapter)
	require.Equal(t, "/tmp/lora-training", trainer.workDir)
}

func TestSetWorkDir(t *testing.T) {
	adapter := adapters.NewLoRAAdapter(4, 8.0)
	trainer := NewLoRATrainer(adapter)
	
	trainer.SetWorkDir("/custom/path")
	require.Equal(t, "/custom/path", trainer.workDir)
}

func TestGetAndSetAdapterWeights(t *testing.T) {
	adapter := adapters.NewLoRAAdapter(4, 8.0)
	trainer := NewLoRATrainer(adapter)
	
	weights, err := trainer.GetAdapterWeights()
	require.NoError(t, err)
	require.NotEmpty(t, weights)
	
	newWeights := map[string][]float64{
		"q_proj": {1.0, 2.0},
		"v_proj": {3.0, 4.0},
	}
	
	trainer.SetAdapterWeights(newWeights)
	updatedWeights, err := trainer.GetAdapterWeights()
	require.NoError(t, err)
	require.Equal(t, newWeights, updatedWeights)
}

func TestSimulateTraining(t *testing.T) {
	adapter := adapters.NewLoRAAdapter(4, 8.0)
	trainer := NewLoRATrainer(adapter)
	
	originalWeights := trainer.adapter.GetWeights()
	
	err := trainer.simulateTraining()
	require.NoError(t, err)
	
	updatedWeights := trainer.adapter.GetWeights()
	require.NotEqual(t, originalWeights, updatedWeights)
}

func TestSaveAndLoadAdapterConfig(t *testing.T) {
	adapter := adapters.NewLoRAAdapter(4, 8.0)
	trainer := NewLoRATrainer(adapter)
	
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")
	
	err := trainer.saveAdapterConfig(configPath)
	require.NoError(t, err)
	require.FileExists(t, configPath)
	
	err = trainer.loadWeightsFromFile(configPath)
	require.NoError(t, err)
}

