package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	computekeeper "github.com/atlas/chain/x/compute/keeper"
	trainingkeeper "github.com/atlas/chain/x/training/keeper"
	healthkeeper "github.com/atlas/chain/x/health/keeper"
)

type Keeper struct {
	cdc codec.BinaryCodec
	storeKey storetypes.StoreKey
	memKey storetypes.StoreKey
	trainingKeeper trainingkeeper.Keeper
	computeKeeper computekeeper.Keeper
	healthKeeper healthkeeper.Keeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey, memKey storetypes.StoreKey,
	trainingKeeper trainingkeeper.Keeper,
	computeKeeper computekeeper.Keeper,
	healthKeeper healthkeeper.Keeper,
) *Keeper {
	return &Keeper{
		cdc: cdc, storeKey: storeKey, memKey: memKey,
		trainingKeeper: trainingKeeper,
		computeKeeper: computeKeeper,
		healthKeeper: healthKeeper,
	}
}

