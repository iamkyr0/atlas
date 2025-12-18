package validator

import (
	"context"
	"fmt"
)

type Validator struct {
	chainRPCURL string
	client      BlockchainClient
}

func NewValidator(chainRPCURL string) *Validator {
	client := NewHTTPBlockchainClient(chainRPCURL)
	return &Validator{
		chainRPCURL: chainRPCURL,
		client:      client,
	}
}

func NewValidatorWithClient(chainRPCURL string, client BlockchainClient) *Validator {
	return &Validator{
		chainRPCURL: chainRPCURL,
		client:      client,
	}
}

func (v *Validator) ValidateAssignment(ctx context.Context, shardID string, nodeID string) error {
	duplicate, err := v.CheckDuplication(ctx, shardID)
	if err != nil {
		return fmt.Errorf("duplication check failed: %w", err)
	}
	if duplicate {
		return fmt.Errorf("shard already assigned")
	}

	capacity, err := v.client.QueryNodeCapacity(ctx, nodeID)
	if err != nil {
		return fmt.Errorf("capacity check failed: %w", err)
	}
	if capacity <= 0 {
		return fmt.Errorf("node has no available capacity")
	}

	reputation, err := v.client.QueryNodeReputation(ctx, nodeID)
	if err != nil {
		return fmt.Errorf("reputation check failed: %w", err)
	}
	if reputation < 0.0 {
		return fmt.Errorf("node has negative reputation")
	}

	return nil
}

func (v *Validator) CheckDuplication(ctx context.Context, shardID string) (bool, error) {
	assignments, err := v.client.QueryShardAssignments(ctx, shardID)
	if err != nil {
		return false, err
	}
	return len(assignments) > 0, nil
}

func ValidateAssignment(shardID string, nodeID string) error {
	validator := NewValidator("http://localhost:26657")
	return validator.ValidateAssignment(context.Background(), shardID, nodeID)
}

func CheckDuplication(shardID string) (bool, error) {
	validator := NewValidator("http://localhost:26657")
	return validator.CheckDuplication(context.Background(), shardID)
}

