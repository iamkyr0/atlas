package proof

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGenerateProof(t *testing.T) {
	metrics := ComputationMetrics{
		Iterations:    100,
		TimeElapsed:   5 * time.Second,
		MemoryUsed:    1024,
		GPUUtilization: 0.85,
	}
	
	proof, err := GenerateProof("task-1", "node-1", metrics)
	require.NoError(t, err)
	require.NotNil(t, proof)
	require.Equal(t, "task-1", proof.TaskID)
	require.Equal(t, "node-1", proof.NodeID)
	require.NotEmpty(t, proof.Hash)
	require.Equal(t, metrics, proof.Metrics)
}

func TestVerifyProof(t *testing.T) {
	metrics := ComputationMetrics{
		Iterations:    100,
		TimeElapsed:   5 * time.Second,
		MemoryUsed:    1024,
		GPUUtilization: 0.85,
	}
	
	proof, err := GenerateProof("task-1", "node-1", metrics)
	require.NoError(t, err)
	
	valid, err := VerifyProof(proof)
	require.NoError(t, err)
	require.True(t, valid)
	
	proof.Hash = "invalid"
	valid, err = VerifyProof(proof)
	require.NoError(t, err)
	require.False(t, valid)
}

