package keeper

import (
	"fmt"
	"time"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/atlas/chain/x/compute/types"
	computekeeper "github.com/atlas/chain/x/compute/keeper"
)

const (
	HeartbeatTimeout = 90 * time.Second
)

func (k Keeper) CheckNodeHealth(ctx sdk.Context, nodeID string) (bool, error) {
	node, found := k.computeKeeper.GetNode(ctx, nodeID)
	if !found {
		return false, fmt.Errorf("node not found")
	}

	timeSinceLastHeartbeat := ctx.BlockTime().Sub(node.LastHeartbeat)
	if timeSinceLastHeartbeat > HeartbeatTimeout {
		node.Status = "offline"
		k.computeKeeper.SetNode(ctx, node)
		return false, nil
	}

	return true, nil
}

func (k Keeper) UpdateHeartbeat(ctx sdk.Context, nodeID string) error {
	node, found := k.computeKeeper.GetNode(ctx, nodeID)
	if !found {
		return fmt.Errorf("node not found")
	}

	node.LastHeartbeat = ctx.BlockTime()
	node.Status = "online"
	k.computeKeeper.SetNode(ctx, node)

	return nil
}

func (k Keeper) GetOfflineNodes(ctx sdk.Context) []types.Node {
	allNodes := k.computeKeeper.GetAllNodes(ctx)
	var offlineNodes []types.Node

	for _, node := range allNodes {
		timeSinceLastHeartbeat := ctx.BlockTime().Sub(node.LastHeartbeat)
		if timeSinceLastHeartbeat > HeartbeatTimeout {
			offlineNodes = append(offlineNodes, node)
		}
	}

	return offlineNodes
}

