package proof

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type ProofOfComputation struct {
	TaskID      string
	NodeID      string
	Timestamp   time.Time
	Hash        string
	Signature   string
	Metrics     ComputationMetrics
}

type ComputationMetrics struct {
	Iterations   int
	TimeElapsed  time.Duration
	MemoryUsed   uint64
	GPUUtilization float64
}

func GenerateProof(taskID string, nodeID string, metrics ComputationMetrics) (*ProofOfComputation, error) {
	timestamp := time.Now()
	
	data := fmt.Sprintf("%s:%s:%d:%d", taskID, nodeID, timestamp.Unix(), metrics.Iterations)
	hash := sha256.Sum256([]byte(data))
	
	proof := &ProofOfComputation{
		TaskID:    taskID,
		NodeID:    nodeID,
		Timestamp: timestamp,
		Hash:      fmt.Sprintf("%x", hash),
		Metrics:   metrics,
	}
	
	return proof, nil
}

func VerifyProof(proof *ProofOfComputation) (bool, error) {
	data := fmt.Sprintf("%s:%s:%d:%d", proof.TaskID, proof.NodeID, proof.Timestamp.Unix(), proof.Metrics.Iterations)
	hash := sha256.Sum256([]byte(data))
	expectedHash := fmt.Sprintf("%x", hash)
	
	return proof.Hash == expectedHash, nil
}

