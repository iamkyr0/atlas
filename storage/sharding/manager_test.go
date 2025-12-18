package sharding

import (
	"fmt"
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
	shardMgr := NewShardManager(ipfsManager)
	
	cids, err := shardMgr.SplitDataset(datasetPath, 4, outputDir)
	require.NoError(t, err)
	require.Len(t, cids, 4)
	
	for i := 0; i < 4; i++ {
		shardPath := filepath.Join(outputDir, fmt.Sprintf("shard_%d", i))
		require.FileExists(t, shardPath)
	}
}

func TestSplitModelByChunks(t *testing.T) {
	tempDir := t.TempDir()
	modelPath := filepath.Join(tempDir, "model.pt")
	
	data := make([]byte, 1000)
	err := os.WriteFile(modelPath, data, 0644)
	require.NoError(t, err)
	
	outputDir := filepath.Join(tempDir, "chunks")
	err = os.MkdirAll(outputDir, 0755)
	require.NoError(t, err)
	
	ipfsManager := manager.NewIPFSManager("/ip4/127.0.0.1/tcp/5001")
	shardMgr := NewShardManager(ipfsManager)
	
	cids, err := shardMgr.SplitModelByChunks(modelPath, outputDir)
	require.NoError(t, err)
	require.Len(t, cids, 4)
}

func TestValidateShardHash(t *testing.T) {
	tempDir := t.TempDir()
	shardPath := filepath.Join(tempDir, "shard.bin")
	
	data := []byte("test shard data")
	err := os.WriteFile(shardPath, data, 0644)
	require.NoError(t, err)
	
	ipfsManager := manager.NewIPFSManager("/ip4/127.0.0.1/tcp/5001")
	shardMgr := NewShardManager(ipfsManager)
	
	hash, err := CalculateShardHash(shardPath)
	require.NoError(t, err)
	
	valid, err := shardMgr.ValidateShardHash(shardPath, hash)
	require.NoError(t, err)
	require.True(t, valid)
	
	valid, err = shardMgr.ValidateShardHash(shardPath, "invalid_hash")
	require.NoError(t, err)
	require.False(t, valid)
}

