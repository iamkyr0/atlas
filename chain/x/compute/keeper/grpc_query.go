package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/atlas/chain/x/compute/types"
)

type QueryServer struct {
	Keeper
}

func NewQueryServer(keeper Keeper) QueryServer {
	return QueryServer{Keeper: keeper}
}

func (qs QueryServer) GetNode(ctx context.Context, req *types.QueryGetNodeRequest) (*types.QueryGetNodeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.NodeId == "" {
		return nil, status.Error(codes.InvalidArgument, "node_id cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	node, found := qs.Keeper.GetNode(sdkCtx, req.NodeId)
	if !found {
		return nil, status.Error(codes.NotFound, "node not found")
	}

	return &types.QueryGetNodeResponse{Node: &node}, nil
}

func (qs QueryServer) ListNodes(ctx context.Context, req *types.QueryListNodesRequest) (*types.QueryListNodesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	nodes := qs.Keeper.GetAllNodes(sdkCtx)

	return &types.QueryListNodesResponse{Nodes: nodes}, nil
}

