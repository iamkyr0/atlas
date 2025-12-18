package sharding

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"github.com/atlas/storage/validation"
)

type ShardManager struct {
	ipfsManager interface{} // IPFS manager interface
}

func NewShardManager(ipfsManager interface{}) *ShardManager {
	return &ShardManager{
		ipfsManager: ipfsManager,
	}
}

func (sm *ShardManager) SplitDataset(datasetPath string, numShards int, outputDir string) ([]string, error) {
	file, err := os.Open(datasetPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open dataset: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := stat.Size()
	shardSize := fileSize / int64(numShards)

	var shardCIDs []string
	for i := 0; i < numShards; i++ {
		shardPath := filepath.Join(outputDir, fmt.Sprintf("shard_%d", i))
		shardFile, err := os.Create(shardPath)
		if err != nil {
			return nil, err
		}

		start := int64(i) * shardSize
		end := start + shardSize
		if i == numShards-1 {
			end = fileSize
		}

		file.Seek(start, 0)
		io.CopyN(shardFile, file, end-start)
		shardFile.Close()

		if ipfsMgr, ok := sm.ipfsManager.(interface{ AddFile(string) (string, error) }); ok {
			cid, err := ipfsMgr.AddFile(shardPath)
			if err != nil {
				return nil, fmt.Errorf("failed to upload shard %d to IPFS: %w", i, err)
			}
			shardCIDs = append(shardCIDs, cid)
		} else {
			shardCIDs = append(shardCIDs, fmt.Sprintf("shard_cid_%d", i))
		}
	}

	return shardCIDs, nil
}

func (sm *ShardManager) SplitModelByLayers(modelPath string, outputDir string) ([]string, error) {
	cid, err := sm.ipfsManager.(interface{ AddFile(string) (string, error) }).AddFile(modelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to upload model: %w", err)
	}
	return []string{cid}, nil
}

func (sm *ShardManager) SplitModelByChunks(modelPath string, outputDir string) ([]string, error) {
	file, err := os.Open(modelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open model: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := stat.Size()
	chunkSize := fileSize / 4

	var chunkCIDs []string
	for i := 0; i < 4; i++ {
		chunkPath := filepath.Join(outputDir, fmt.Sprintf("chunk_%d", i))
		chunkFile, err := os.Create(chunkPath)
		if err != nil {
			return nil, err
		}

		start := int64(i) * chunkSize
		end := start + chunkSize
		if i == 3 {
			end = fileSize
		}

		file.Seek(start, 0)
		io.CopyN(chunkFile, file, end-start)
		chunkFile.Close()

		if ipfsMgr, ok := sm.ipfsManager.(interface{ AddFile(string) (string, error) }); ok {
			cid, err := ipfsMgr.AddFile(chunkPath)
			if err != nil {
				return nil, fmt.Errorf("failed to upload chunk %d to IPFS: %w", i, err)
			}
			chunkCIDs = append(chunkCIDs, cid)
		} else {
			chunkCIDs = append(chunkCIDs, fmt.Sprintf("chunk_cid_%d", i))
		}
	}

	return chunkCIDs, nil
}

func (sm *ShardManager) ValidateShardHash(shardPath string, expectedHash string) (bool, error) {
	hash, err := validation.CalculateHash(shardPath)
	if err != nil {
		return false, fmt.Errorf("failed to calculate hash: %w", err)
	}
	return hash == expectedHash, nil
}

