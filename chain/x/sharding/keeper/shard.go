package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

type Shard struct {
	ID       string
	JobID    string
	CID      string
	Hash     string
	NodeID   string
	Status   string
	Size     int64
}

func (k Keeper) RegisterShard(ctx sdk.Context, shard *Shard) error {
	store := ctx.KVStore(k.storeKey)
	
	existing := store.Get([]byte("shard:" + shard.ID))
	if existing != nil {
		return fmt.Errorf("shard already exists")
	}

	bz := k.cdc.MustMarshal(&shard)
	store.Set([]byte("shard:"+shard.ID), bz)

	return nil
}

func (k Keeper) GetShard(ctx sdk.Context, shardID string) (*Shard, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte("shard:" + shardID))
	if bz == nil {
		return nil, false
	}

	var shard Shard
	k.cdc.MustUnmarshal(bz, &shard)
	return &shard, true
}

func (k Keeper) GetShardForValidation(ctx sdk.Context, shardID string) (*Shard, bool) {
	return k.GetShard(ctx, shardID)
}

func (k Keeper) AssignShardToNode(ctx sdk.Context, shardID string, nodeID string) error {
	shard, found := k.GetShard(ctx, shardID)
	if !found {
		return fmt.Errorf("shard not found")
	}

	if shard.NodeID != "" && shard.NodeID != nodeID {
		return fmt.Errorf("shard already assigned to another node")
	}

	shard.NodeID = nodeID
	shard.Status = "assigned"

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&shard)
	store.Set([]byte("shard:"+shard.ID), bz)

	return nil
}

func (k Keeper) GetShardsForJob(ctx sdk.Context, jobID string) []*Shard {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte("shard:"))
	defer iterator.Close()

	var shards []*Shard
	for ; iterator.Valid(); iterator.Next() {
		var shard Shard
		k.cdc.MustUnmarshal(iterator.Value(), &shard)
		if shard.JobID == jobID {
			shards = append(shards, &shard)
		}
	}
	return shards
}

func (k Keeper) GetShardsByNode(ctx sdk.Context, nodeID string) []*Shard {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte("shard:"))
	defer iterator.Close()

	var shards []*Shard
	for ; iterator.Valid(); iterator.Next() {
		var shard Shard
		k.cdc.MustUnmarshal(iterator.Value(), &shard)
		if shard.NodeID == nodeID {
			shards = append(shards, &shard)
		}
	}
	return shards
}

func (k Keeper) GetShardsByHash(ctx sdk.Context, hash string) []*Shard {
	if hash == "" {
		return []*Shard{}
	}
	
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte("shard:"))
	defer iterator.Close()

	var shards []*Shard
	for ; iterator.Valid(); iterator.Next() {
		var shard Shard
		k.cdc.MustUnmarshal(iterator.Value(), &shard)
		if shard.Hash == hash {
			shards = append(shards, &shard)
		}
	}
	return shards
}

func (k Keeper) GetAllShards(ctx sdk.Context) []*Shard {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte("shard:"))
	defer iterator.Close()

	var shards []*Shard
	for ; iterator.Valid(); iterator.Next() {
		var shard Shard
		k.cdc.MustUnmarshal(iterator.Value(), &shard)
		shards = append(shards, &shard)
	}
	return shards
}

