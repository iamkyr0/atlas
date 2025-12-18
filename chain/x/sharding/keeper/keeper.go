package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	storagekeeper "github.com/atlas/chain/x/storage/keeper"
	trainingkeeper "github.com/atlas/chain/x/training/keeper"
)

type Keeper struct {
	cdc codec.BinaryCodec
	storeKey storetypes.StoreKey
	memKey storetypes.StoreKey
	storageKeeper storagekeeper.Keeper
	trainingKeeper trainingkeeper.Keeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey, memKey storetypes.StoreKey,
	storageKeeper storagekeeper.Keeper,
	trainingKeeper trainingkeeper.Keeper,
) *Keeper {
	return &Keeper{
		cdc: cdc, storeKey: storeKey, memKey: memKey,
		storageKeeper: storageKeeper,
		trainingKeeper: trainingKeeper,
	}
}

