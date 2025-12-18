package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	computekeeper "github.com/atlas/chain/x/compute/keeper"
	storagekeeper "github.com/atlas/chain/x/storage/keeper"
	"github.com/atlas/chain/x/training/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
	memKey   storetypes.StoreKey

	bankKeeper    bankkeeper.Keeper
	computeKeeper computekeeper.Keeper
	storageKeeper storagekeeper.Keeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	computeKeeper computekeeper.Keeper,
	storageKeeper storagekeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
) *Keeper {
	return &Keeper{
		cdc:           cdc,
		storeKey:     storeKey,
		memKey:        memKey,
		computeKeeper: computeKeeper,
		storageKeeper: storageKeeper,
		bankKeeper:   bankKeeper,
	}
}

func (k Keeper) GetJob(ctx sdk.Context, id string) (types.Job, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte("job:" + id))
	if bz == nil {
		return types.Job{}, false
	}

	var job types.Job
	k.cdc.MustUnmarshal(bz, &job)
	return job, true
}

func (k Keeper) SetJob(ctx sdk.Context, job types.Job) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&job)
	store.Set([]byte("job:"+job.ID), bz)
}

func (k Keeper) GetTask(ctx sdk.Context, id string) (types.Task, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte("task:" + id))
	if bz == nil {
		return types.Task{}, false
	}

	var task types.Task
	k.cdc.MustUnmarshal(bz, &task)
	return task, true
}

func (k Keeper) SetTask(ctx sdk.Context, task types.Task) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&task)
	store.Set([]byte("task:"+task.ID), bz)
}

func (k Keeper) IterateTasks(ctx sdk.Context, handler func(task types.Task) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte("task:"))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var task types.Task
		k.cdc.MustUnmarshal(iterator.Value(), &task)
		if handler(task) {
			break
		}
	}
}

