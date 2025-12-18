package keeper

import (
	"testing"

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
)

func setupKeeper(t *testing.T) (*Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey("storage")
	memStoreKey := storetypes.NewMemoryStoreKey("mem_storage")

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

func TestRegisterStorageNode(t *testing.T) {
	k, ctx := setupKeeper(t)

	node := StorageNode{
		ID:          "node-1",
		Address:     "cosmos1abc123",
		IPFSAddress: "/ip4/127.0.0.1/tcp/5001",
		Capacity:    1000,
		Used:        100,
		Status:      "online",
	}

	err := k.RegisterStorageNode(ctx, node)
	require.NoError(t, err)

	retrievedNode, found := k.GetStorageNode(ctx, "node-1")
	require.True(t, found)
	require.Equal(t, node.ID, retrievedNode.ID)
	require.Equal(t, node.Address, retrievedNode.Address)

	err = k.RegisterStorageNode(ctx, node)
	require.Error(t, err)
}

func TestGetStorageNode(t *testing.T) {
	k, ctx := setupKeeper(t)

	node := StorageNode{
		ID:          "node-1",
		Address:     "cosmos1abc123",
		IPFSAddress: "/ip4/127.0.0.1/tcp/5001",
		Capacity:    1000,
		Used:        100,
		Status:      "online",
	}

	k.RegisterStorageNode(ctx, node)

	retrievedNode, found := k.GetStorageNode(ctx, "node-1")
	require.True(t, found)
	require.Equal(t, node.ID, retrievedNode.ID)

	_, found = k.GetStorageNode(ctx, "nonexistent")
	require.False(t, found)
}

func TestGetAllStorageNodes(t *testing.T) {
	k, ctx := setupKeeper(t)

	node1 := StorageNode{
		ID:          "node-1",
		Address:     "cosmos1abc123",
		IPFSAddress: "/ip4/127.0.0.1/tcp/5001",
		Capacity:    1000,
		Used:        100,
		Status:      "online",
	}

	node2 := StorageNode{
		ID:          "node-2",
		Address:     "cosmos1def456",
		IPFSAddress: "/ip4/127.0.0.1/tcp/5002",
		Capacity:    2000,
		Used:        200,
		Status:      "online",
	}

	k.RegisterStorageNode(ctx, node1)
	k.RegisterStorageNode(ctx, node2)

	allNodes := k.GetAllStorageNodes(ctx)
	require.Len(t, allNodes, 2)
}

func TestUpdateStorageNodeCapacity(t *testing.T) {
	k, ctx := setupKeeper(t)

	node := StorageNode{
		ID:          "node-1",
		Address:     "cosmos1abc123",
		IPFSAddress: "/ip4/127.0.0.1/tcp/5001",
		Capacity:    1000,
		Used:        100,
		Status:      "online",
	}

	k.RegisterStorageNode(ctx, node)

	err := k.UpdateStorageNodeCapacity(ctx, "node-1", 200)
	require.NoError(t, err)

	updatedNode, found := k.GetStorageNode(ctx, "node-1")
	require.True(t, found)
	require.Equal(t, int64(200), updatedNode.Used)

	err = k.UpdateStorageNodeCapacity(ctx, "nonexistent", 200)
	require.Error(t, err)
}

func TestGetAvailableStorageNodes(t *testing.T) {
	k, ctx := setupKeeper(t)

	node1 := StorageNode{
		ID:          "node-1",
		Address:     "cosmos1abc123",
		IPFSAddress: "/ip4/127.0.0.1/tcp/5001",
		Capacity:    1000,
		Used:        100,
		Status:      "online",
	}

	node2 := StorageNode{
		ID:          "node-2",
		Address:     "cosmos1def456",
		IPFSAddress: "/ip4/127.0.0.1/tcp/5002",
		Capacity:    2000,
		Used:        2000,
		Status:      "online",
	}

	node3 := StorageNode{
		ID:          "node-3",
		Address:     "cosmos1ghi789",
		IPFSAddress: "/ip4/127.0.0.1/tcp/5003",
		Capacity:    3000,
		Used:        500,
		Status:      "offline",
	}

	k.RegisterStorageNode(ctx, node1)
	k.RegisterStorageNode(ctx, node2)
	k.RegisterStorageNode(ctx, node3)

	available := k.GetAvailableStorageNodes(ctx)
	require.Len(t, available, 1)
	require.Equal(t, "node-1", available[0].ID)
}

