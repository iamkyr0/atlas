package keeper

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/atlas/chain/x/model/types"
)

func setupKeeper(t *testing.T) (*Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey("model")
	memStoreKey := storetypes.NewMemoryStoreKey("mem_model")

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	k := NewKeeper(cdc, storeKey, memStoreKey)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	return k, ctx
}

func TestRegisterModel(t *testing.T) {
	k, ctx := setupKeeper(t)

	model := types.Model{
		ID:        "model-1",
		Name:      "test-model",
		Version:   "1.0.0",
		CID:       "QmModel123",
		Metadata:  make(map[string]string),
		CreatedAt: time.Now(),
	}

	err := k.RegisterModel(ctx, model)
	require.NoError(t, err)

	retrievedModel, found := k.GetModel(ctx, "model-1")
	require.True(t, found)
	require.Equal(t, model.ID, retrievedModel.ID)
	require.Equal(t, model.Name, retrievedModel.Name)

	err = k.RegisterModel(ctx, model)
	require.Error(t, err)
}

func TestGetModel(t *testing.T) {
	k, ctx := setupKeeper(t)

	model := types.Model{
		ID:        "model-1",
		Name:      "test-model",
		Version:   "1.0.0",
		CID:       "QmModel123",
		Metadata:  make(map[string]string),
		CreatedAt: time.Now(),
	}

	k.RegisterModel(ctx, model)

	retrievedModel, found := k.GetModel(ctx, "model-1")
	require.True(t, found)
	require.Equal(t, model.ID, retrievedModel.ID)

	_, found = k.GetModel(ctx, "nonexistent")
	require.False(t, found)
}

func TestGetAllModels(t *testing.T) {
	k, ctx := setupKeeper(t)

	model1 := types.Model{
		ID:        "model-1",
		Name:      "test-model-1",
		Version:   "1.0.0",
		CID:       "QmModel123",
		Metadata:  make(map[string]string),
		CreatedAt: time.Now(),
	}

	model2 := types.Model{
		ID:        "model-2",
		Name:      "test-model-2",
		Version:   "2.0.0",
		CID:       "QmModel456",
		Metadata:  make(map[string]string),
		CreatedAt: time.Now(),
	}

	k.RegisterModel(ctx, model1)
	k.RegisterModel(ctx, model2)

	allModels := k.GetAllModels(ctx)
	require.Len(t, allModels, 2)
}

func TestUpdateModelVersion(t *testing.T) {
	k, ctx := setupKeeper(t)

	model := types.Model{
		ID:        "model-1",
		Name:      "test-model",
		Version:   "1.0.0",
		CID:       "QmModel123",
		Metadata:  make(map[string]string),
		CreatedAt: time.Now(),
	}

	k.RegisterModel(ctx, model)

	err := k.UpdateModelVersion(ctx, "model-1", "2.0.0", "QmModel456")
	require.NoError(t, err)

	updatedModel, found := k.GetModel(ctx, "model-1")
	require.True(t, found)
	require.Equal(t, "2.0.0", updatedModel.Version)
	require.Equal(t, "QmModel456", updatedModel.CID)

	err = k.UpdateModelVersion(ctx, "nonexistent", "2.0.0", "QmModel456")
	require.Error(t, err)
}

func TestGetModelsByCID(t *testing.T) {
	k, ctx := setupKeeper(t)

	model1 := types.Model{
		ID:        "model-1",
		Name:      "test-model-1",
		Version:   "1.0.0",
		CID:       "QmModel123",
		Metadata:  make(map[string]string),
		CreatedAt: time.Now(),
	}

	model2 := types.Model{
		ID:        "model-2",
		Name:      "test-model-2",
		Version:   "2.0.0",
		CID:       "QmModel123",
		Metadata:  make(map[string]string),
		CreatedAt: time.Now(),
	}

	model3 := types.Model{
		ID:        "model-3",
		Name:      "test-model-3",
		Version:   "3.0.0",
		CID:       "QmModel456",
		Metadata:  make(map[string]string),
		CreatedAt: time.Now(),
	}

	k.RegisterModel(ctx, model1)
	k.RegisterModel(ctx, model2)
	k.RegisterModel(ctx, model3)

	models := k.GetModelsByCID(ctx, "QmModel123")
	require.Len(t, models, 2)

	models = k.GetModelsByCID(ctx, "QmModel456")
	require.Len(t, models, 1)

	models = k.GetModelsByCID(ctx, "nonexistent")
	require.Len(t, models, 0)
}

