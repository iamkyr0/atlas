package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockBlockchainClient struct {
	shardAssignments map[string][]string
	nodeCapacity     map[string]int
	nodeReputation   map[string]float64
}

func (m *mockBlockchainClient) QueryShardAssignments(ctx context.Context, shardID string) ([]string, error) {
	if assignments, ok := m.shardAssignments[shardID]; ok {
		return assignments, nil
	}
	return []string{}, nil
}

func (m *mockBlockchainClient) QueryNodeCapacity(ctx context.Context, nodeID string) (int, error) {
	if capacity, ok := m.nodeCapacity[nodeID]; ok {
		return capacity, nil
	}
	return 10, nil
}

func (m *mockBlockchainClient) QueryNodeReputation(ctx context.Context, nodeID string) (float64, error) {
	if reputation, ok := m.nodeReputation[nodeID]; ok {
		return reputation, nil
	}
	return 100.0, nil
}

func TestValidateAssignment(t *testing.T) {
	mockClient := &mockBlockchainClient{
		shardAssignments: make(map[string][]string),
		nodeCapacity:     make(map[string]int),
		nodeReputation:   make(map[string]float64),
	}
	
	mockClient.nodeCapacity["node-1"] = 5
	mockClient.nodeReputation["node-1"] = 95.0
	
	validator := NewValidatorWithClient("http://localhost:26657", mockClient)
	
	err := validator.ValidateAssignment(context.Background(), "shard-1", "node-1")
	require.NoError(t, err)
	
	mockClient.shardAssignments["shard-1"] = []string{"node-2"}
	err = validator.ValidateAssignment(context.Background(), "shard-1", "node-1")
	require.Error(t, err)
	
	mockClient.nodeCapacity["node-1"] = 0
	err = validator.ValidateAssignment(context.Background(), "shard-2", "node-1")
	require.Error(t, err)
	
	mockClient.nodeCapacity["node-1"] = 5
	mockClient.nodeReputation["node-1"] = -10.0
	err = validator.ValidateAssignment(context.Background(), "shard-2", "node-1")
	require.Error(t, err)
}

func TestCheckDuplication(t *testing.T) {
	mockClient := &mockBlockchainClient{
		shardAssignments: make(map[string][]string),
		nodeCapacity:     make(map[string]int),
		nodeReputation:   make(map[string]float64),
	}
	
	validator := NewValidatorWithClient("http://localhost:26657", mockClient)
	
	duplicate, err := validator.CheckDuplication(context.Background(), "shard-1")
	require.NoError(t, err)
	require.False(t, duplicate)
	
	mockClient.shardAssignments["shard-1"] = []string{"node-1"}
	duplicate, err = validator.CheckDuplication(context.Background(), "shard-1")
	require.NoError(t, err)
	require.True(t, duplicate)
}

