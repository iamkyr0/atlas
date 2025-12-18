package sharding

import (
	"context"
	"fmt"
	"time"
)

type Coordinator struct {
	shards map[string]*ShardState
}

type ShardState struct {
	ID       string
	NodeID   string
	Status   string
	Progress float64
}

func NewCoordinator() *Coordinator {
	return &Coordinator{
		shards: make(map[string]*ShardState),
	}
}

func (c *Coordinator) AssignShard(shardID string, nodeID string) error {
	c.shards[shardID] = &ShardState{
		ID:     shardID,
		NodeID: nodeID,
		Status: "assigned",
	}
	return nil
}

func (c *Coordinator) WaitForShards(ctx context.Context, timeoutSeconds int) error {
	timeout := time.Duration(timeoutSeconds) * time.Second
	deadline := time.Now().Add(timeout)
	
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			allCompleted := true
			for _, shard := range c.shards {
				if shard.Status != "completed" {
					allCompleted = false
					break
				}
			}
			
			if allCompleted {
				return nil
			}
			
			if time.Now().After(deadline) {
				return fmt.Errorf("timeout waiting for shards")
			}
		}
	}
}

func (c *Coordinator) UpdateShardStatus(shardID string, status string, progress float64) {
	if shard, ok := c.shards[shardID]; ok {
		shard.Status = status
		shard.Progress = progress
	}
}

func (c *Coordinator) GetShardStatus(shardID string) (*ShardState, bool) {
	shard, ok := c.shards[shardID]
	return shard, ok
}

