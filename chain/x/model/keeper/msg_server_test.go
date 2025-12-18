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

	"github.com/atlas/chain/x/model/types"
)

func setupMsgServer(t *testing.T) (MsgServer, sdk.Context) {
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
	ms := NewMsgServer(keeper)

	ctx := sdk.NewContext(stateStore, tmproto.Header{Time: time.Now()}, false, log.NewNopLogger())

	return ms, ctx
}

func TestRegisterModel(t *testing.T) {
	ms, ctx := setupMsgServer(t)

	msg := &types.MsgRegisterModel{
		Creator:  "cosmos1abc123",
		Name:     "test-model",
		Version:  "1.0.0",
		Cid:      "QmModel123",
		Metadata: make(map[string]string),
	}

	resp, err := ms.RegisterModel(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.ModelId)

	model, found := ms.Keeper.GetModel(ctx, resp.ModelId)
	require.True(t, found)
	require.Equal(t, msg.Name, model.Name)
	require.Equal(t, msg.Version, model.Version)
	require.Equal(t, msg.Cid, model.CID)

	_, err = ms.RegisterModel(context.Background(), nil)
	require.Error(t, err)
}

