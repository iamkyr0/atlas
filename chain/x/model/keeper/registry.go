package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/atlas/chain/x/model/types"
)

func (k Keeper) RegisterModel(ctx sdk.Context, model types.Model) error {
	store := ctx.KVStore(k.storeKey)
	
	existing := store.Get([]byte("model:" + model.ID))
	if existing != nil {
		return fmt.Errorf("model already exists")
	}

	bz := k.cdc.MustMarshal(&model)
	store.Set([]byte("model:"+model.ID), bz)

	return nil
}

func (k Keeper) GetModel(ctx sdk.Context, modelID string) (types.Model, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte("model:" + modelID))
	if bz == nil {
		return types.Model{}, false
	}

	var model types.Model
	k.cdc.MustUnmarshal(bz, &model)
	return model, true
}

func (k Keeper) GetAllModels(ctx sdk.Context) []types.Model {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte("model:"))
	defer iterator.Close()

	var models []types.Model
	for ; iterator.Valid(); iterator.Next() {
		var model types.Model
		k.cdc.MustUnmarshal(iterator.Value(), &model)
		models = append(models, model)
	}
	return models
}

func (k Keeper) UpdateModelVersion(ctx sdk.Context, modelID string, newVersion string, newCID string) error {
	model, found := k.GetModel(ctx, modelID)
	if !found {
		return fmt.Errorf("model not found")
	}

	model.Version = newVersion
	model.CID = newCID

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&model)
	store.Set([]byte("model:"+model.ID), bz)

	return nil
}

func (k Keeper) GetModelsByCID(ctx sdk.Context, cid string) []types.Model {
	allModels := k.GetAllModels(ctx)
	var matching []types.Model
	
	for _, model := range allModels {
		if model.CID == cid {
			matching = append(matching, model)
		}
	}
	return matching
}

