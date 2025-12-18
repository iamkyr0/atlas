package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) ValidateShardAssignment(ctx sdk.Context, shardID string, nodeID string) error {
	shard, found := k.shardingKeeper.GetShard(ctx, shardID)
	if !found {
		return fmt.Errorf("shard not found")
	}

	if shard.NodeID != "" && shard.NodeID != nodeID {
		return fmt.Errorf("shard already assigned to node %s", shard.NodeID)
	}

	nodeShards := k.shardingKeeper.GetShardsByNode(ctx, nodeID)
	for _, nodeShard := range nodeShards {
		if nodeShard.ID == shardID {
			continue
		}
		if nodeShard.Hash != "" && shard.Hash != "" && nodeShard.Hash == shard.Hash {
			return fmt.Errorf("duplicate shard content detected (same hash)")
		}
	}

	return nil
}

func (k Keeper) CheckDuplicateShard(ctx sdk.Context, shardHash string) (bool, error) {
	if shardHash == "" {
		return false, nil
	}
	
	shards := k.shardingKeeper.GetShardsByHash(ctx, shardHash)
	
	return len(shards) > 0, nil
}

func (k Keeper) ValidateTaskAssignment(ctx sdk.Context, taskID string, nodeID string) error {
	task, found := k.trainingKeeper.GetTask(ctx, taskID)
	if !found {
		return fmt.Errorf("task not found")
	}

	if task.NodeID != "" && task.NodeID != nodeID {
		return fmt.Errorf("task already assigned to node %s", task.NodeID)
	}

	node, nodeFound := k.computeKeeper.GetNode(ctx, nodeID)
	if !nodeFound {
		return fmt.Errorf("node not found")
	}

	if node.Status != "online" {
		return fmt.Errorf("node is not online")
	}

	isHealthy, err := k.healthKeeper.CheckNodeHealth(ctx, nodeID)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	if !isHealthy {
		return fmt.Errorf("node is not healthy")
	}

	return nil
}

