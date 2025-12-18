package adapters

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLoRAAdapter(t *testing.T) {
	adapter := NewLoRAAdapter(8, 16.0)
	require.NotNil(t, adapter)
	require.Equal(t, 8, adapter.rank)
	require.Equal(t, 16.0, adapter.alpha)
	require.Equal(t, 0.1, adapter.dropout)
	require.Len(t, adapter.targetModules, 2)
	require.NotEmpty(t, adapter.weights)
}

func TestApply(t *testing.T) {
	adapter := NewLoRAAdapter(4, 8.0)
	
	err := adapter.Apply("test-model")
	require.NoError(t, err)
	
	err = adapter.Apply(nil)
	require.Error(t, err)
}

func TestSaveAndLoad(t *testing.T) {
	adapter := NewLoRAAdapter(4, 8.0)
	
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "adapter.json")
	
	err := adapter.Save(path)
	require.NoError(t, err)
	
	require.FileExists(t, path)
	
	newAdapter := NewLoRAAdapter(4, 8.0)
	newAdapter.weights = make(map[string][]float64)
	
	err = newAdapter.Load(path)
	require.NoError(t, err)
	require.Equal(t, adapter.rank, newAdapter.rank)
	require.Equal(t, adapter.alpha, newAdapter.alpha)
	require.NotEmpty(t, newAdapter.weights)
	
	err = newAdapter.Load("nonexistent.json")
	require.Error(t, err)
}

func TestGetAndSetWeights(t *testing.T) {
	adapter := NewLoRAAdapter(4, 8.0)
	
	weights := adapter.GetWeights()
	require.NotEmpty(t, weights)
	
	newWeights := map[string][]float64{
		"q_proj": {1.0, 2.0, 3.0},
		"v_proj": {4.0, 5.0, 6.0},
	}
	
	adapter.SetWeights(newWeights)
	require.Equal(t, newWeights, adapter.GetWeights())
}

