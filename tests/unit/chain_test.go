package unit

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/atlas/chain/x/compute/types"
)

func TestNodeValidation(t *testing.T) {
	node := types.Node{
		ID:      "test-node-1",
		Address: "localhost:8080",
		Status:  "online",
	}
	
	err := node.Validate()
	assert.NoError(t, err)
}

func TestNodeValidationEmptyID(t *testing.T) {
	node := types.Node{
		ID:      "",
		Address: "localhost:8080",
	}
	
	err := node.Validate()
	assert.Error(t, err)
}

