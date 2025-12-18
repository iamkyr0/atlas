package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	trainingkeeper "github.com/atlas/chain/x/training/keeper"
	shardingkeeper "github.com/atlas/chain/x/sharding/keeper"
	computekeeper "github.com/atlas/chain/x/compute/keeper"
	healthkeeper "github.com/atlas/chain/x/health/keeper"
)

type Keeper struct {
	cdc codec.BinaryCodec
	storeKey storetypes.StoreKey
	memKey storetypes.StoreKey
	trainingKeeper trainingkeeper.Keeper
	shardingKeeper shardingkeeper.Keeper
	computeKeeper computekeeper.Keeper
	healthKeeper healthkeeper.Keeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey, memKey storetypes.StoreKey,
	trainingKeeper trainingkeeper.Keeper,
	shardingKeeper shardingkeeper.Keeper,
	computeKeeper computekeeper.Keeper,
	healthKeeper healthkeeper.Keeper,
) *Keeper {
	return &Keeper{
		cdc: cdc, storeKey: storeKey, memKey: memKey,
		trainingKeeper: trainingKeeper,
		shardingKeeper: shardingKeeper,
		computeKeeper: computeKeeper,
		healthKeeper: healthKeeper,
	}
}

