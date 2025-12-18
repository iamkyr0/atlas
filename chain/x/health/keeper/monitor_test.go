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

	computekeeper "github.com/atlas/chain/x/compute/keeper"
	computetypes "github.com/atlas/chain/x/compute/types"
)

func setupKeeper(t *testing.T) (*Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey("health")
	memStoreKey := storetypes.NewMemoryStoreKey("mem_health")
	computeStoreKey := sdk.NewKVStoreKey(computetypes.StoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(computeStoreKey, storetypes.StoreTypeIAVL, db)
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

	computeKeeper := computekeeper.NewKeeper(cdc, computeStoreKey, storetypes.NewMemoryStoreKey("mem_compute"), bankKeeper)

	k := NewKeeper(cdc, storeKey, memStoreKey, computeKeeper)

	ctx := sdk.NewContext(stateStore, tmproto.Header{Time: time.Now()}, false, log.NewNopLogger())

	return k, ctx
}

func TestCheckNodeHealth(t *testing.T) {
	k, ctx := setupKeeper(t)

	node := computetypes.Node{
		ID:            "node-1",
		Address:       "cosmos1abc123",
		Status:        "online",
		Resources:     make(map[string]string),
		Reputation:    100.0,
		UptimePercent: 99.5,
		LastHeartbeat: ctx.BlockTime(),
	}
	k.computeKeeper.SetNode(ctx, node)

	healthy, err := k.CheckNodeHealth(ctx, "node-1")
	require.NoError(t, err)
	require.True(t, healthy)

	node.LastHeartbeat = ctx.BlockTime().Add(-100 * time.Second)
	k.computeKeeper.SetNode(ctx, node)

	healthy, err = k.CheckNodeHealth(ctx, "node-1")
	require.NoError(t, err)
	require.False(t, healthy)

	_, err = k.CheckNodeHealth(ctx, "nonexistent")
	require.Error(t, err)
}

func TestUpdateHeartbeat(t *testing.T) {
	k, ctx := setupKeeper(t)

	node := computetypes.Node{
		ID:            "node-1",
		Address:       "cosmos1abc123",
		Status:        "offline",
		Resources:     make(map[string]string),
		Reputation:    100.0,
		UptimePercent: 99.5,
		LastHeartbeat: ctx.BlockTime().Add(-200 * time.Second),
	}
	k.computeKeeper.SetNode(ctx, node)

	err := k.UpdateHeartbeat(ctx, "node-1")
	require.NoError(t, err)

	updatedNode, found := k.computeKeeper.GetNode(ctx, "node-1")
	require.True(t, found)
	require.Equal(t, "online", updatedNode.Status)

	err = k.UpdateHeartbeat(ctx, "nonexistent")
	require.Error(t, err)
}

func TestGetOfflineNodes(t *testing.T) {
	k, ctx := setupKeeper(t)

	node1 := computetypes.Node{
		ID:            "node-1",
		Address:       "cosmos1abc123",
		Status:        "online",
		Resources:     make(map[string]string),
		Reputation:    100.0,
		UptimePercent: 99.5,
		LastHeartbeat: ctx.BlockTime(),
	}

	node2 := computetypes.Node{
		ID:            "node-2",
		Address:       "cosmos1def456",
		Status:        "online",
		Resources:     make(map[string]string),
		Reputation:    95.0,
		UptimePercent: 98.0,
		LastHeartbeat: ctx.BlockTime().Add(-100 * time.Second),
	}

	k.computeKeeper.SetNode(ctx, node1)
	k.computeKeeper.SetNode(ctx, node2)

	offlineNodes := k.GetOfflineNodes(ctx)
	require.Len(t, offlineNodes, 1)
	require.Equal(t, "node-2", offlineNodes[0].ID)
}

