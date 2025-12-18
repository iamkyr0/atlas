package compute

import (
	"github.com/atlas/chain/x/compute/keeper"
	"github.com/atlas/chain/x/compute/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	
	nodes := k.GetAllNodes(ctx)
	genesis.Nodes = nodes
	
	return genesis
}

