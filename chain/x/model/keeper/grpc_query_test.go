package keeper

import (
	"context"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/atlas/chain/x/model/types"
)

func setupQueryServer(t *testing.T) (QueryServer, sdk.Context) {
	storeKey := sdk.NewKVStoreKey("model")
	memStoreKey := storetypes.NewMemoryStoreKey("mem_model")

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	keeper := NewKeeper(cdc, storeKey, memStoreKey)
	qs := NewQueryServer(keeper)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	return qs, ctx
}

func TestGetModel(t *testing.T) {
	qs, ctx := setupQueryServer(t)

	model := types.Model{
		ID:        "model-1",
		Name:      "test-model",
		Version:   "1.0.0",
		CID:       "QmModel123",
		Metadata:  make(map[string]string),
		CreatedAt: time.Now(),
	}

	qs.Keeper.RegisterModel(ctx, model)

	req := &types.QueryGetModelRequest{ModelId: "model-1"}
	resp, err := qs.GetModel(sdk.WrapSDKContext(ctx), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, model.ID, resp.Model.Id)
	require.Equal(t, model.Name, resp.Model.Name)

	req.ModelId = ""
	_, err = qs.GetModel(sdk.WrapSDKContext(ctx), req)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))

	req.ModelId = "nonexistent"
	_, err = qs.GetModel(sdk.WrapSDKContext(ctx), req)
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))

	_, err = qs.GetModel(context.Background(), nil)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestListModels(t *testing.T) {
	qs, ctx := setupQueryServer(t)

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

	qs.Keeper.RegisterModel(ctx, model1)
	qs.Keeper.RegisterModel(ctx, model2)

	req := &types.QueryListModelsRequest{}
	resp, err := qs.ListModels(sdk.WrapSDKContext(ctx), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Models, 2)

	_, err = qs.ListModels(context.Background(), nil)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

