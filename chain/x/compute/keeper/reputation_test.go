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

	"github.com/atlas/chain/x/compute/types"
)

func setupReputationKeeper(t *testing.T) (*Keeper, sdk.Context) {
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

	k := NewKeeper(cdc, storeKey, memStoreKey, bankKeeper)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	return k, ctx
}

func TestUpdateReputation(t *testing.T) {
	k, ctx := setupReputationKeeper(t)

	node := types.Node{
		ID:            "node-1",
		Address:       "cosmos1abc123",
		Status:        "online",
		Resources:     make(map[string]string),
		Reputation:    50.0,
		UptimePercent: 50.0,
		LastHeartbeat: time.Now(),
		RegisteredAt:  time.Now(),
		ActiveTasks:   []string{},
	}

	k.SetNode(ctx, node)

	k.UpdateReputation(ctx, "node-1", 99.0)
	updatedNode, _ := k.GetNode(ctx, "node-1")
	require.Equal(t, 99.0, updatedNode.UptimePercent)
	require.Equal(t, 99.0, updatedNode.Reputation)

	k.UpdateReputation(ctx, "node-1", 40.0)
	updatedNode, _ = k.GetNode(ctx, "node-1")
	require.Equal(t, 40.0, updatedNode.UptimePercent)
	require.Less(t, updatedNode.Reputation, 40.0)

	k.UpdateReputation(ctx, "nonexistent", 99.0)
}

func TestGetNodeReputation(t *testing.T) {
	k, ctx := setupReputationKeeper(t)

	node := types.Node{
		ID:            "node-1",
		Address:       "cosmos1abc123",
		Status:        "online",
		Resources:     make(map[string]string),
		Reputation:    100.0,
		UptimePercent: 99.5,
		LastHeartbeat: time.Now(),
		RegisteredAt:  time.Now(),
		ActiveTasks:   []string{},
	}

	k.SetNode(ctx, node)

	reputation := k.GetNodeReputation(ctx, "node-1")
	require.Equal(t, 100.0, reputation)

	reputation = k.GetNodeReputation(ctx, "nonexistent")
	require.Equal(t, 0.0, reputation)
}

func TestUpdateHeartbeat(t *testing.T) {
	k, ctx := setupReputationKeeper(t)

	node := types.Node{
		ID:            "node-1",
		Address:       "cosmos1abc123",
		Status:        "offline",
		Resources:     make(map[string]string),
		Reputation:    100.0,
		UptimePercent: 99.5,
		LastHeartbeat: time.Now().Add(-100 * time.Second),
		RegisteredAt:  time.Now(),
		ActiveTasks:   []string{},
	}

	k.SetNode(ctx, node)

	k.UpdateHeartbeat(ctx, "node-1")
	updatedNode, _ := k.GetNode(ctx, "node-1")
	require.Equal(t, "online", updatedNode.Status)

	k.UpdateHeartbeat(ctx, "nonexistent")
}

