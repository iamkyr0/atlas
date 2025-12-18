package sharding

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/atlas/storage/manager"
	"github.com/stretchr/testify/require"
)

func TestSplitDataset(t *testing.T) {
	tempDir := t.TempDir()
	datasetPath := filepath.Join(tempDir, "dataset.bin")
	
	data := make([]byte, 1000)
	for i := range data {
		data[i] = byte(i % 256)
	}
	
	err := os.WriteFile(datasetPath, data, 0644)
	require.NoError(t, err)
	
	outputDir := filepath.Join(tempDir, "shards")
	err = os.MkdirAll(outputDir, 0755)
	require.NoError(t, err)
	
	ipfsManager := manager.NewIPFSManager("/ip4/127.0.0.1/tcp/5001")
	
	cids, err := SplitDataset(ipfsManager, datasetPath, 4, outputDir)
	require.NoError(t, err)
	require.Len(t, cids, 4)
}

func TestSplitModel(t *testing.T) {
	tempDir := t.TempDir()
	modelPath := filepath.Join(tempDir, "model.pt")
	
	data := make([]byte, 1000)
	err := os.WriteFile(modelPath, data, 0644)
	require.NoError(t, err)
	
	outputDir := filepath.Join(tempDir, "chunks")
	err = os.MkdirAll(outputDir, 0755)
	require.NoError(t, err)
	
	ipfsManager := manager.NewIPFSManager("/ip4/127.0.0.1/tcp/5001")
	
	cids, err := SplitModel(ipfsManager, modelPath, "chunk", outputDir)
	require.NoError(t, err)
	require.NotEmpty(t, cids)
	
	cids, err = SplitModel(ipfsManager, modelPath, "layer", outputDir)
	require.NoError(t, err)
	require.NotEmpty(t, cids)
	
	cids, err = SplitModel(ipfsManager, modelPath, "default", outputDir)
	require.NoError(t, err)
	require.Len(t, cids, 1)
}

func TestCalculateShardHash(t *testing.T) {
	tempDir := t.TempDir()
	shardPath := filepath.Join(tempDir, "shard.bin")
	
	data := []byte("test shard data")
	err := os.WriteFile(shardPath, data, 0644)
	require.NoError(t, err)
	
	hash, err := CalculateShardHash(shardPath)
	require.NoError(t, err)
	require.NotEmpty(t, hash)
	require.Len(t, hash, 64)
	
	hash2, err := CalculateShardHash(shardPath)
	require.NoError(t, err)
	require.Equal(t, hash, hash2)
}

