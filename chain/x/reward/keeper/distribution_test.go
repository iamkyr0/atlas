package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	computekeeper "github.com/atlas/chain/x/compute/keeper"
	computetypes "github.com/atlas/chain/x/compute/types"
	storagekeeper "github.com/atlas/chain/x/storage/keeper"
)

func setupRewardKeeper(t *testing.T) (*Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey("reward")
	memStoreKey := storetypes.NewMemoryStoreKey("mem_reward")
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

	k := NewKeeper(cdc, storeKey, memStoreKey, bankKeeper, computeKeeper, storageKeeper)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	moduleAddr := authtypes.NewModuleAddress("reward")
	coins := sdk.NewCoins(sdk.NewCoin("uatlas", sdk.NewInt(1000000)))
	bankKeeper.MintCoins(ctx, moduleAddr, coins)

	return k, ctx
}

func TestCalculateReward(t *testing.T) {
	k, ctx := setupRewardKeeper(t)

	node := computetypes.Node{
		ID:            "node-1",
		Address:       "cosmos1abc123",
		Status:        "online",
		Resources:     make(map[string]string),
		Reputation:    100.0,
		UptimePercent: 99.5,
	}
	k.computeKeeper.SetNode(ctx, node)

	baseReward := sdk.NewCoin("uatlas", sdk.NewInt(1000))
	reward := k.CalculateReward(ctx, "node-1", 1.0, baseReward)
	require.Equal(t, baseReward.Denom, reward.Denom)
	require.True(t, reward.Amount.IsPositive())

	reward = k.CalculateReward(ctx, "node-1", 0.5, baseReward)
	require.True(t, reward.Amount.LT(baseReward.Amount))

	reward = k.CalculateReward(ctx, "node-1", 0.0, baseReward)
	require.True(t, reward.Amount.IsZero() || reward.Amount.IsPositive())

	node.Reputation = 50.0
	k.computeKeeper.SetNode(ctx, node)
	reward = k.CalculateReward(ctx, "node-1", 1.0, baseReward)
	require.True(t, reward.Amount.LT(baseReward.Amount))
}

func TestDistributeReward(t *testing.T) {
	k, ctx := setupRewardKeeper(t)

	nodeAddr, err := sdk.AccAddressFromBech32("cosmos1abc123")
	require.NoError(t, err)

	amount := sdk.NewCoin("uatlas", sdk.NewInt(100))
	err = k.DistributeReward(ctx, nodeAddr.String(), amount, "test reward")
	require.NoError(t, err)
}

