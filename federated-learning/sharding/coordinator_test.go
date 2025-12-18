package sharding

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAssignShard(t *testing.T) {
	coordinator := NewCoordinator()
	
	err := coordinator.AssignShard("shard-1", "node-1")
	require.NoError(t, err)
	
	shard, found := coordinator.GetShardStatus("shard-1")
	require.True(t, found)
	require.Equal(t, "shard-1", shard.ID)
	require.Equal(t, "node-1", shard.NodeID)
	require.Equal(t, "assigned", shard.Status)
}

func TestUpdateShardStatus(t *testing.T) {
	coordinator := NewCoordinator()
	coordinator.AssignShard("shard-1", "node-1")
	
	coordinator.UpdateShardStatus("shard-1", "completed", 1.0)
	
	shard, found := coordinator.GetShardStatus("shard-1")
	require.True(t, found)
	require.Equal(t, "completed", shard.Status)
	require.Equal(t, 1.0, shard.Progress)
}

func TestWaitForShards(t *testing.T) {
	coordinator := NewCoordinator()
	coordinator.AssignShard("shard-1", "node-1")
	coordinator.AssignShard("shard-2", "node-2")
	
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	
	go func() {
		time.Sleep(100 * time.Millisecond)
		coordinator.UpdateShardStatus("shard-1", "completed", 1.0)
		coordinator.UpdateShardStatus("shard-2", "completed", 1.0)
	}()
	
	err := coordinator.WaitForShards(ctx, 5)
	require.NoError(t, err)
}

func TestWaitForShardsTimeout(t *testing.T) {
	coordinator := NewCoordinator()
	coordinator.AssignShard("shard-1", "node-1")
	
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	
	err := coordinator.WaitForShards(ctx, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "timeout")
}

