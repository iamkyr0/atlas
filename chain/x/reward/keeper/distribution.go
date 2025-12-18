package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (k Keeper) DistributeReward(ctx sdk.Context, nodeID string, amount sdk.Coin, reason string) error {
	nodeAddr, err := sdk.AccAddressFromBech32(nodeID)
	if err != nil {
		return fmt.Errorf("invalid node address: %w", err)
	}

	moduleAddr := authtypes.NewModuleAddress("reward")
	
	if err := k.bankKeeper.SendCoins(ctx, moduleAddr, nodeAddr, sdk.NewCoins(amount)); err != nil {
		return fmt.Errorf("failed to send reward: %w", err)
	}

	return nil
}

func (k Keeper) CalculateReward(ctx sdk.Context, nodeID string, workCompleted float64, baseReward sdk.Coin) sdk.Coin {
	reputation := k.computeKeeper.GetNodeReputation(ctx, nodeID)
	reputationMultiplier := reputation / 100.0
	if reputationMultiplier > 1.0 {
		reputationMultiplier = 1.0
	}
	if reputationMultiplier < 0.0 {
		reputationMultiplier = 0.0
	}
	
	if workCompleted < 0.0 {
		workCompleted = 0.0
	}
	if workCompleted > 1.0 {
		workCompleted = 1.0
	}
	
	workAmount := sdk.NewDecFromInt(baseReward.Amount).Mul(sdk.NewDecFromFloat64(workCompleted))
	reputationAmount := workAmount.Mul(sdk.NewDecFromFloat64(reputationMultiplier))
	
	amount := reputationAmount.TruncateInt()
	if amount.IsNegative() {
		amount = sdk.ZeroInt()
	}
	
	return sdk.NewCoin(baseReward.Denom, amount)
}

