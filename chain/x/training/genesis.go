package training

import (
	"github.com/atlas/chain/x/training/keeper"
	"github.com/atlas/chain/x/training/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.DefaultGenesis()
}

