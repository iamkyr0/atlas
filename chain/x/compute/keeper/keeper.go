package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/atlas/chain/x/compute/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
	memKey   storetypes.StoreKey

	bankKeeper bankkeeper.Keeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	bankKeeper bankkeeper.Keeper,
) *Keeper {
	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		bankKeeper: bankKeeper,
	}
}

func (k Keeper) GetNode(ctx sdk.Context, id string) (types.Node, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(id))
	if bz == nil {
		return types.Node{}, false
	}

	var node types.Node
	k.cdc.MustUnmarshal(bz, &node)
	return node, true
}

func (k Keeper) SetNode(ctx sdk.Context, node types.Node) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&node)
	store.Set([]byte(node.ID), bz)
}

func (k Keeper) GetAllNodes(ctx sdk.Context) []types.Node {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	var nodes []types.Node
	for ; iterator.Valid(); iterator.Next() {
		var node types.Node
		k.cdc.MustUnmarshal(iterator.Value(), &node)
		nodes = append(nodes, node)
	}
	return nodes
}

func (k Keeper) IterateNodes(ctx sdk.Context, handler func(node types.Node) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var node types.Node
		k.cdc.MustUnmarshal(iterator.Value(), &node)
		if handler(node) {
			break
		}
	}
}

