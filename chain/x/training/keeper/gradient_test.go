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
	storagekeeper "github.com/atlas/chain/x/storage/keeper"
	"github.com/atlas/chain/x/training/types"
)

func setupGradientKeeper(t *testing.T) (*Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	computeStoreKey := sdk.NewKVStoreKey(computetypes.StoreKey)
	storageStoreKey := sdk.NewKVStoreKey("storage")

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(computeStoreKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(storageStoreKey, storetypes.StoreTypeIAVL, db)
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
	storageKeeper := storagekeeper.NewKeeper(cdc, storageStoreKey, storetypes.NewMemoryStoreKey("mem_storage"), bankKeeper)

	k := NewKeeper(cdc, storeKey, memStoreKey, computeKeeper, storageKeeper, bankKeeper)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	return k, ctx
}

func TestTrackGradientContribution(t *testing.T) {
	k, ctx := setupGradientKeeper(t)

	err := k.TrackGradientContribution(ctx, "node-1", "job-1", 1, "QmGradient123", 0.5)
	require.NoError(t, err)

	contributions := k.GetGradientContributions(ctx, "job-1", 1)
	require.Len(t, contributions, 1)
	require.Equal(t, "node-1", contributions[0].NodeID)
	require.Equal(t, "job-1", contributions[0].JobID)
	require.Equal(t, 1, contributions[0].Round)
	require.Equal(t, "QmGradient123", contributions[0].GradientCID)
	require.Equal(t, 0.5, contributions[0].Contribution)
}

func TestGetGradientContributions(t *testing.T) {
	k, ctx := setupGradientKeeper(t)

	k.TrackGradientContribution(ctx, "node-1", "job-1", 1, "QmGradient123", 0.5)
	k.TrackGradientContribution(ctx, "node-2", "job-1", 1, "QmGradient456", 0.3)
	k.TrackGradientContribution(ctx, "node-1", "job-1", 2, "QmGradient789", 0.4)

	contributions := k.GetGradientContributions(ctx, "job-1", 1)
	require.Len(t, contributions, 2)

	contributions = k.GetGradientContributions(ctx, "job-1", 2)
	require.Len(t, contributions, 1)

	contributions = k.GetGradientContributions(ctx, "job-2", 1)
	require.Len(t, contributions, 0)
}

func TestCalculateFairRewards(t *testing.T) {
	k, ctx := setupGradientKeeper(t)

	k.TrackGradientContribution(ctx, "node-1", "job-1", 1, "QmGradient123", 0.5)
	k.TrackGradientContribution(ctx, "node-2", "job-1", 1, "QmGradient456", 0.3)
	k.TrackGradientContribution(ctx, "node-3", "job-1", 1, "QmGradient789", 0.2)

	rewards := k.CalculateFairRewards(ctx, "job-1", 1)
	require.Len(t, rewards, 3)
	require.InDelta(t, 0.5, rewards["node-1"], 0.01)
	require.InDelta(t, 0.3, rewards["node-2"], 0.01)
	require.InDelta(t, 0.2, rewards["node-3"], 0.01)

	rewards = k.CalculateFairRewards(ctx, "job-1", 2)
	require.Len(t, rewards, 0)
}

