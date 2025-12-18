package keeper

import (
	"time"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/atlas/chain/x/compute/types"
)

func (k Keeper) UpdateReputation(ctx sdk.Context, nodeID string, uptimePercent float64) {
	node, found := k.GetNode(ctx, nodeID)
	if !found {
		return
	}
	
	node.UptimePercent = uptimePercent
	node.Reputation = uptimePercent
	
	if uptimePercent < 50.0 {
		node.Reputation *= 0.5
	}
	
	k.SetNode(ctx, node)
}

func (k Keeper) GetNodeReputation(ctx sdk.Context, nodeID string) float64 {
	node, found := k.GetNode(ctx, nodeID)
	if !found {
		return 0.0
	}
	return node.Reputation
}

func (k Keeper) UpdateHeartbeat(ctx sdk.Context, nodeID string) {
	node, found := k.GetNode(ctx, nodeID)
	if !found {
		return
	}
	
	node.LastHeartbeat = time.Now()
	node.Status = "online"
	
	k.SetNode(ctx, node)
}

