package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/atlas/chain/x/model/types"
)

type QueryServer struct {
	Keeper
}

func NewQueryServer(keeper Keeper) QueryServer {
	return QueryServer{Keeper: keeper}
}

func (qs QueryServer) GetModel(ctx context.Context, req *types.QueryGetModelRequest) (*types.QueryGetModelResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ModelId == "" {
		return nil, status.Error(codes.InvalidArgument, "model_id cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	keeperModel, found := qs.Keeper.GetModel(sdkCtx, req.ModelId)
	if !found {
		return nil, status.Error(codes.NotFound, "model not found")
	}

	modelProto := &types.ModelProto{
		Id:        keeperModel.ID,
		Name:      keeperModel.Name,
		Version:   keeperModel.Version,
		Cid:       keeperModel.CID,
		CreatedAt: keeperModel.CreatedAt,
		Metadata:  keeperModel.Metadata,
	}

	return &types.QueryGetModelResponse{Model: modelProto}, nil
}

func (qs QueryServer) ListModels(ctx context.Context, req *types.QueryListModelsRequest) (*types.QueryListModelsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	keeperModels := qs.Keeper.GetAllModels(sdkCtx)

	var models []types.ModelProto
	for _, m := range keeperModels {
		models = append(models, types.ModelProto{
			Id:        m.ID,
			Name:      m.Name,
			Version:   m.Version,
			Cid:       m.CID,
			CreatedAt: m.CreatedAt,
			Metadata:  m.Metadata,
		})
	}

	return &types.QueryListModelsResponse{Models: models}, nil
}

