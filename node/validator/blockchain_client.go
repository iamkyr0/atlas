package validator

import (
	"context"
	"fmt"
)

type BlockchainClient interface {
	QueryShardAssignments(ctx context.Context, shardID string) ([]string, error)
	QueryNodeCapacity(ctx context.Context, nodeID string) (int, error)
	QueryNodeReputation(ctx context.Context, nodeID string) (float64, error)
}

type HTTPBlockchainClient struct {
	rpcURL string
}

func NewHTTPBlockchainClient(rpcURL string) *HTTPBlockchainClient {
	return &HTTPBlockchainClient{rpcURL: rpcURL}
}

func (c *HTTPBlockchainClient) QueryShardAssignments(ctx context.Context, shardID string) ([]string, error) {
	return nil, fmt.Errorf("blockchain client not fully implemented: requires gRPC client with protobuf stubs")
}

func (c *HTTPBlockchainClient) QueryNodeCapacity(ctx context.Context, nodeID string) (int, error) {
	return 0, fmt.Errorf("blockchain client not fully implemented: requires gRPC client with protobuf stubs")
}

func (c *HTTPBlockchainClient) QueryNodeReputation(ctx context.Context, nodeID string) (float64, error) {
	return 0.0, fmt.Errorf("blockchain client not fully implemented: requires gRPC client with protobuf stubs")
}

