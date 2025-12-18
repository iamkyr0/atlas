package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

type StorageNode struct {
	ID          string
	Address     string
	IPFSAddress string
	Capacity    int64
	Used        int64
	Status      string
}

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
	memKey   storetypes.StoreKey
	bankKeeper bankkeeper.Keeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey, memKey storetypes.StoreKey,
	bankKeeper bankkeeper.Keeper,
) *Keeper {
	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		bankKeeper: bankKeeper,
	}
}

func (k Keeper) RegisterStorageNode(ctx sdk.Context, node StorageNode) error {
	store := ctx.KVStore(k.storeKey)
	
	existing := store.Get([]byte("node:" + node.ID))
	if existing != nil {
		return fmt.Errorf("storage node already exists")
	}

	bz := k.cdc.MustMarshal(&node)
	store.Set([]byte("node:"+node.ID), bz)
	return nil
}

func (k Keeper) GetStorageNode(ctx sdk.Context, nodeID string) (StorageNode, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte("node:" + nodeID))
	if bz == nil {
		return StorageNode{}, false
	}

	var node StorageNode
	k.cdc.MustUnmarshal(bz, &node)
	return node, true
}

func (k Keeper) GetAllStorageNodes(ctx sdk.Context) []StorageNode {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte("node:"))
	defer iterator.Close()

	var nodes []StorageNode
	for ; iterator.Valid(); iterator.Next() {
		var node StorageNode
		k.cdc.MustUnmarshal(iterator.Value(), &node)
		nodes = append(nodes, node)
	}
	return nodes
}

func (k Keeper) UpdateStorageNodeCapacity(ctx sdk.Context, nodeID string, used int64) error {
	node, found := k.GetStorageNode(ctx, nodeID)
	if !found {
		return fmt.Errorf("storage node not found")
	}

	node.Used = used
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&node)
	store.Set([]byte("node:"+node.ID), bz)
	return nil
}

func (k Keeper) GetAvailableStorageNodes(ctx sdk.Context) []StorageNode {
	allNodes := k.GetAllStorageNodes(ctx)
	var available []StorageNode
	
	for _, node := range allNodes {
		if node.Status == "online" && node.Used < node.Capacity {
			available = append(available, node)
		}
	}
	return available
}

