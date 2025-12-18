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
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/atlas/chain/x/compute/types"
)

func setupMsgServer(t *testing.T) (MsgServer, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	bankKeeper := bankkeeper.NewBaseKeeper(
		cdc,
		storeKey,
		nil,
		nil,
		banktypes.DefaultGenesisState().DenomMetadata,
	)

	keeper := NewKeeper(cdc, storeKey, memStoreKey, bankKeeper)
	ms := NewMsgServer(keeper)

	ctx := sdk.NewContext(stateStore, tmproto.Header{Time: time.Now()}, false, log.NewNopLogger())

	return ms, ctx
}

func TestRegisterNode(t *testing.T) {
	ms, ctx := setupMsgServer(t)

	msg := &types.MsgRegisterNode{
		Creator:   "cosmos1abc123",
		NodeId:    "node-1",
		Address:   "cosmos1abc123",
		CpuCores:  8,
		GpuCount:  2,
		MemoryGb:  32,
		StorageGb: 500,
	}

	resp, err := ms.RegisterNode(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
	require.NotNil(t, resp)

	node, found := ms.Keeper.GetNode(ctx, "node-1")
	require.True(t, found)
	require.Equal(t, msg.NodeId, node.ID)
	require.Equal(t, msg.Address, node.Address)
	require.Equal(t, "online", node.Status)

	_, err = ms.RegisterNode(sdk.WrapSDKContext(ctx), msg)
	require.Error(t, err)

	_, err = ms.RegisterNode(context.Background(), nil)
	require.Error(t, err)
}

func TestUpdateHeartbeat(t *testing.T) {
	ms, ctx := setupMsgServer(t)

	node := types.Node{
		ID:            "node-1",
		Address:       "cosmos1abc123",
		Status:        "offline",
		Resources:     make(map[string]string),
		Reputation:    100.0,
		UptimePercent: 99.5,
		LastHeartbeat: ctx.BlockTime().Add(-100 * time.Second),
		RegisteredAt:  ctx.BlockTime(),
		ActiveTasks:   []string{},
	}

	ms.Keeper.SetNode(ctx, node)

	msg := &types.MsgUpdateHeartbeat{
		Creator: "cosmos1abc123",
		NodeId:  "node-1",
	}

	resp, err := ms.UpdateHeartbeat(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
	require.NotNil(t, resp)

	updatedNode, found := ms.Keeper.GetNode(ctx, "node-1")
	require.True(t, found)
	require.Equal(t, "online", updatedNode.Status)

	msg.NodeId = "nonexistent"
	_, err = ms.UpdateHeartbeat(sdk.WrapSDKContext(ctx), msg)
	require.Error(t, err)

	_, err = ms.UpdateHeartbeat(context.Background(), nil)
	require.Error(t, err)
}

