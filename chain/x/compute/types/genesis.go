package types

import (
	"fmt"
)

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Nodes: []Node{},
	}
}

func (gs GenesisState) Validate() error {
	for _, node := range gs.Nodes {
		if err := node.Validate(); err != nil {
			return fmt.Errorf("invalid node: %w", err)
		}
	}
	return nil
}

