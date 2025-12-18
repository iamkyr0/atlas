package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/atlas/chain/x/model/keeper"
	"github.com/atlas/chain/x/model/types"
)

type MsgServer struct {
	Keeper
}

func NewMsgServer(keeper Keeper) MsgServer {
	return MsgServer{Keeper: keeper}
}

func (ms MsgServer) RegisterModel(ctx context.Context, msg *types.MsgRegisterModel) (*types.MsgRegisterModelResponse, error) {
	if msg == nil {
		return nil, fmt.Errorf("invalid message")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	modelID := fmt.Sprintf("model-%s-%s", msg.Name, msg.Version)

	metadata := msg.Metadata
	if metadata == nil {
		metadata = make(map[string]string)
	}

	model := types.Model{
		ID:        modelID,
		Name:      msg.Name,
		Version:   msg.Version,
		CID:       msg.Cid,
		Metadata:  metadata,
		CreatedAt: sdkCtx.BlockTime(),
	}

	ms.Keeper.RegisterModel(sdkCtx, model)

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeModelRegistered,
			sdk.NewAttribute(types.AttributeKeyModelID, modelID),
			sdk.NewAttribute(types.AttributeKeyName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyVersion, msg.Version),
		),
	)

	return &types.MsgRegisterModelResponse{ModelId: modelID}, nil
}

