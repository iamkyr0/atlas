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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/atlas/chain/x/compute/types"
)

func setupQueryServer(t *testing.T) (QueryServer, sdk.Context) {
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
	qs := NewQueryServer(keeper)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	return qs, ctx
}

func TestGetNode(t *testing.T) {
	qs, ctx := setupQueryServer(t)

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

	qs.Keeper.SetNode(ctx, node)

	req := &types.QueryGetNodeRequest{NodeId: "node-1"}
	resp, err := qs.GetNode(sdk.WrapSDKContext(ctx), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, node.ID, resp.Node.ID)

	req.NodeId = ""
	_, err = qs.GetNode(sdk.WrapSDKContext(ctx), req)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))

	req.NodeId = "nonexistent"
	_, err = qs.GetNode(sdk.WrapSDKContext(ctx), req)
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))

	_, err = qs.GetNode(context.Background(), nil)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestListNodes(t *testing.T) {
	qs, ctx := setupQueryServer(t)

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

	qs.Keeper.SetNode(ctx, node1)
	qs.Keeper.SetNode(ctx, node2)

	req := &types.QueryListNodesRequest{}
	resp, err := qs.ListNodes(sdk.WrapSDKContext(ctx), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Nodes, 2)

	_, err = qs.ListNodes(context.Background(), nil)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

