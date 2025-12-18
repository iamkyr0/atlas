package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrNodeNotFound = sdkerrors.Register(ModuleName, 1, "node not found")
	ErrNodeExists   = sdkerrors.Register(ModuleName, 2, "node already exists")
	ErrInvalidNode  = sdkerrors.Register(ModuleName, 3, "invalid node")
)

const (
	EventTypeNodeRegistered = "node_registered"
	EventTypeHeartbeatUpdated = "heartbeat_updated"
	
	AttributeKeyNodeID  = "node_id"
	AttributeKeyAddress = "address"
)
