package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	computekeeper "github.com/atlas/chain/x/compute/keeper"
	computetypes "github.com/atlas/chain/x/compute/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	memKey     storetypes.StoreKey
	bankKeeper bankkeeper.Keeper
	computeKeeper computekeeper.Keeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	bankKeeper bankkeeper.Keeper,
	computeKeeper computekeeper.Keeper,
) *Keeper {
	return &Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		memKey:       memKey,
		bankKeeper:   bankKeeper,
		computeKeeper: computeKeeper,
	}
}

func (k Keeper) SelectNodeForInference(ctx sdk.Context, modelID string, strategy string) (string, error) {
	nodes := k.computeKeeper.GetAllNodes(ctx)
	
	if len(nodes) == 0 {
		return "", fmt.Errorf("no nodes available")
	}
	
	onlineNodes := []computetypes.Node{}
	for _, node := range nodes {
		if node.Status == "online" && len(node.ActiveTasks) < 10 {
			onlineNodes = append(onlineNodes, node)
		}
	}
	
	if len(onlineNodes) == 0 {
		return "", fmt.Errorf("no online nodes available")
	}
	
	switch strategy {
	case "round_robin":
		return k.selectRoundRobin(ctx, onlineNodes)
	case "least_loaded":
		return k.selectLeastLoaded(ctx, onlineNodes)
	case "best_reputation":
		return k.selectBestReputation(ctx, onlineNodes)
	default:
		return k.selectRoundRobin(ctx, onlineNodes)
	}
}

func (k Keeper) selectRoundRobin(ctx sdk.Context, nodes []computetypes.Node) (string, error) {
	store := ctx.KVStore(k.storeKey)
	key := []byte("round_robin_index")
	
	var index uint64
	bz := store.Get(key)
	if bz != nil {
		index = sdk.BigEndianToUint64(bz)
	}
	
	selected := nodes[index%uint64(len(nodes))]
	index++
	store.Set(key, sdk.Uint64ToBigEndian(index))
	
	return selected.ID, nil
}

func (k Keeper) selectLeastLoaded(ctx sdk.Context, nodes []computetypes.Node) (string, error) {
	leastLoaded := nodes[0]
	minTasks := len(nodes[0].ActiveTasks)
	
	for _, node := range nodes[1:] {
		taskCount := len(node.ActiveTasks)
		if taskCount < minTasks {
			minTasks = taskCount
			leastLoaded = node
		}
	}
	
	return leastLoaded.ID, nil
}

func (k Keeper) selectBestReputation(ctx sdk.Context, nodes []computetypes.Node) (string, error) {
	bestNode := nodes[0]
	bestReputation := nodes[0].Reputation
	
	for _, node := range nodes[1:] {
		if node.Reputation > bestReputation {
			bestReputation = node.Reputation
			bestNode = node
		}
	}
	
	return bestNode.ID, nil
}

func (k Keeper) RecordInferenceRequest(ctx sdk.Context, nodeID string, modelID string, latencyMs int64) {
	store := ctx.KVStore(k.storeKey)
	key := []byte(fmt.Sprintf("inference:%s:%s", nodeID, modelID))
	
	var count uint64
	bz := store.Get(key)
	if bz != nil {
		count = sdk.BigEndianToUint64(bz)
	}
	count++
	store.Set(key, sdk.Uint64ToBigEndian(count))
}

