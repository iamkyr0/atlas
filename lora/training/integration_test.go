package training

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLoRAFLIntegration(t *testing.T) {
	integration := NewLoRAFLIntegration(4, 8.0)
	require.NotNil(t, integration)
	require.NotNil(t, integration.adapter)
	require.NotNil(t, integration.trainer)
}

func TestUpdateAdapter(t *testing.T) {
	integration := NewLoRAFLIntegration(4, 8.0)
	
	aggregatedWeights := map[string][]float64{
		"q_proj": {1.0, 2.0, 3.0},
		"v_proj": {4.0, 5.0, 6.0},
	}
	
	err := integration.UpdateAdapter(aggregatedWeights)
	require.NoError(t, err)
	
	weights, err := integration.trainer.GetAdapterWeights()
	require.NoError(t, err)
	require.Equal(t, aggregatedWeights, weights)
}

func TestSaveAndLoadAdapter(t *testing.T) {
	integration := NewLoRAFLIntegration(4, 8.0)
	
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "adapter.json")
	
	err := integration.SaveAdapter(path)
	require.NoError(t, err)
	require.FileExists(t, path)
	
	err = integration.LoadAdapter(path)
	require.NoError(t, err)
}

