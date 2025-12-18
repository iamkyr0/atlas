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

func setupKeeper(t *testing.T) (*Keeper, sdk.Context) {
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

func TestGetNode(t *testing.T) {
	k, ctx := setupKeeper(t)

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

	retrievedNode, found := k.GetNode(ctx, "node-1")
	require.True(t, found)
	require.Equal(t, node.ID, retrievedNode.ID)
	require.Equal(t, node.Address, retrievedNode.Address)

	_, found = k.GetNode(ctx, "nonexistent")
	require.False(t, found)
}

func TestSetNode(t *testing.T) {
	k, ctx := setupKeeper(t)

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

	retrievedNode, found := k.GetNode(ctx, "node-1")
	require.True(t, found)
	require.Equal(t, node.ID, retrievedNode.ID)
}

func TestGetAllNodes(t *testing.T) {
	k, ctx := setupKeeper(t)

	node1 := types.Node{
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

	node2 := types.Node{
		ID:            "node-2",
		Address:       "cosmos1def456",
		Status:        "online",
		Resources:     make(map[string]string),
		Reputation:    95.0,
		UptimePercent: 98.0,
		LastHeartbeat: time.Now(),
		RegisteredAt:  time.Now(),
		ActiveTasks:   []string{},
	}

	k.SetNode(ctx, node1)
	k.SetNode(ctx, node2)

	allNodes := k.GetAllNodes(ctx)
	require.Len(t, allNodes, 2)
}

func TestIterateNodes(t *testing.T) {
	k, ctx := setupKeeper(t)

	node1 := types.Node{
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

	node2 := types.Node{
		ID:            "node-2",
		Address:       "cosmos1def456",
		Status:        "online",
		Resources:     make(map[string]string),
		Reputation:    95.0,
		UptimePercent: 98.0,
		LastHeartbeat: time.Now(),
		RegisteredAt:  time.Now(),
		ActiveTasks:   []string{},
	}

	k.SetNode(ctx, node1)
	k.SetNode(ctx, node2)

	count := 0
	k.IterateNodes(ctx, func(node types.Node) (stop bool) {
		count++
		return false
	})

	require.Equal(t, 2, count)

	count = 0
	k.IterateNodes(ctx, func(node types.Node) (stop bool) {
		count++
		return true
	})

	require.Equal(t, 1, count)
}

