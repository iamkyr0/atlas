package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/atlas/chain/x/training/types"
)

type GradientContribution struct {
	NodeID      string
	JobID       string
	Round       int
	GradientCID string
	Contribution float64
	Timestamp   int64
}

func (k Keeper) TrackGradientContribution(ctx sdk.Context, nodeID string, jobID string, round int, gradientCID string, contribution float64) error {
	key := fmt.Sprintf("gradient:%s:%s:%d:%s", jobID, nodeID, round, gradientCID)
	
	contributionData := GradientContribution{
		NodeID:      nodeID,
		JobID:       jobID,
		Round:       round,
		GradientCID: gradientCID,
		Contribution: contribution,
		Timestamp:   ctx.BlockTime().Unix(),
	}
	
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&contributionData)
	store.Set([]byte(key), bz)
	
	return nil
}

func (k Keeper) GetGradientContributions(ctx sdk.Context, jobID string, round int) []GradientContribution {
	store := ctx.KVStore(k.storeKey)
	prefix := fmt.Sprintf("gradient:%s:", jobID)
	iterator := sdk.KVStorePrefixIterator(store, []byte(prefix))
	defer iterator.Close()
	
	var contributions []GradientContribution
	for ; iterator.Valid(); iterator.Next() {
		var contribution GradientContribution
		k.cdc.MustUnmarshal(iterator.Value(), &contribution)
		
		if contribution.Round == round {
			contributions = append(contributions, contribution)
		}
	}
	
	return contributions
}

func (k Keeper) CalculateFairRewards(ctx sdk.Context, jobID string, round int) map[string]float64 {
	contributions := k.GetGradientContributions(ctx, jobID, round)
	
	if len(contributions) == 0 {
		return map[string]float64{}
	}
	
	totalContribution := 0.0
	for _, c := range contributions {
		totalContribution += c.Contribution
	}
	
	rewards := make(map[string]float64)
	for _, c := range contributions {
		if totalContribution > 0 {
			rewards[c.NodeID] = c.Contribution / totalContribution
		}
	}
	
	return rewards
}

