package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	computekeeper "github.com/atlas/chain/x/compute/keeper"
	storagekeeper "github.com/atlas/chain/x/storage/keeper"
)

type Keeper struct {
	cdc codec.BinaryCodec
	storeKey storetypes.StoreKey
	memKey storetypes.StoreKey
	bankKeeper bankkeeper.Keeper
	computeKeeper computekeeper.Keeper
	storageKeeper storagekeeper.Keeper
	accountKeeper interface{}
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey, memKey storetypes.StoreKey,
	bankKeeper bankkeeper.Keeper,
	computeKeeper computekeeper.Keeper,
	storageKeeper storagekeeper.Keeper,
) *Keeper {
	return &Keeper{
		cdc: cdc, storeKey: storeKey, memKey: memKey,
		bankKeeper: bankKeeper,
		computeKeeper: computeKeeper,
		storageKeeper: storageKeeper,
	}
}

