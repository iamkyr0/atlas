package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/atlas/chain/x/compute/types"
)

type MsgServer struct {
	Keeper
}

func NewMsgServer(keeper Keeper) MsgServer {
	return MsgServer{Keeper: keeper}
}

func (ms MsgServer) RegisterNode(ctx context.Context, msg *types.MsgRegisterNode) (*types.MsgRegisterNodeResponse, error) {
	if msg == nil {
		return nil, fmt.Errorf("invalid message")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	_, found := ms.Keeper.GetNode(sdkCtx, msg.NodeId)
	if found {
		return nil, sdkerrors.Wrapf(types.ErrNodeExists, "node %s already exists", msg.NodeId)
	}

	node := types.Node{
		ID:            msg.NodeId,
		Address:       msg.Address,
		Status:        "online",
		Resources:      make(map[string]string),
		Reputation:     0.0,
		UptimePercent:  0.0,
		LastHeartbeat:  sdkCtx.BlockTime(),
		RegisteredAt:   sdkCtx.BlockTime(),
		ActiveTasks:    []string{},
	}
	
	node.Resources["cpu_cores"] = fmt.Sprintf("%d", msg.CpuCores)
	node.Resources["gpu_count"] = fmt.Sprintf("%d", msg.GpuCount)
	node.Resources["memory_gb"] = fmt.Sprintf("%d", msg.MemoryGb)
	node.Resources["storage_gb"] = fmt.Sprintf("%d", msg.StorageGb)

	ms.Keeper.SetNode(sdkCtx, node)

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeNodeRegistered,
			sdk.NewAttribute(types.AttributeKeyNodeID, msg.NodeId),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address),
		),
	)

	return &types.MsgRegisterNodeResponse{}, nil
}

func (ms MsgServer) UpdateHeartbeat(ctx context.Context, msg *types.MsgUpdateHeartbeat) (*types.MsgUpdateHeartbeatResponse, error) {
	if msg == nil {
		return nil, fmt.Errorf("invalid message")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	node, found := ms.Keeper.GetNode(sdkCtx, msg.NodeId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrNodeNotFound, "node %s not found", msg.NodeId)
	}

	node.LastHeartbeat = sdkCtx.BlockTime()
	node.Status = "online"
	ms.Keeper.SetNode(sdkCtx, node)

	return &types.MsgUpdateHeartbeatResponse{}, nil
}

