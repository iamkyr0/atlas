package sharding

import (
	"fmt"
	"github.com/atlas/storage/manager"
	"github.com/atlas/storage/validation"
)

func SplitDataset(ipfsManager *manager.IPFSManager, datasetPath string, numShards int, outputDir string) ([]string, error) {
	shardMgr := NewShardManager(ipfsManager)
	cids, err := shardMgr.SplitDataset(datasetPath, numShards, outputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to split dataset: %w", err)
	}
	return cids, nil
}

func SplitModel(ipfsManager *manager.IPFSManager, modelPath string, strategy string, outputDir string) ([]string, error) {
	shardMgr := NewShardManager(ipfsManager)
	
	switch strategy {
	case "layer":
		return shardMgr.SplitModelByLayers(modelPath, outputDir)
	case "chunk":
		return shardMgr.SplitModelByChunks(modelPath, outputDir)
	default:
		cid, err := ipfsManager.AddFile(modelPath)
		if err != nil {
			return nil, fmt.Errorf("failed to upload model: %w", err)
		}
		return []string{cid}, nil
	}
}

func CalculateShardHash(shardPath string) (string, error) {
	return validation.CalculateHash(shardPath)
}

