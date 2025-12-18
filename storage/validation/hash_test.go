package validation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalculateHash(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.bin")
	
	data := []byte("test data")
	err := os.WriteFile(filePath, data, 0644)
	require.NoError(t, err)
	
	hash, err := CalculateHash(filePath)
	require.NoError(t, err)
	require.NotEmpty(t, hash)
	require.Len(t, hash, 64)
	
	hash2, err := CalculateHash(filePath)
	require.NoError(t, err)
	require.Equal(t, hash, hash2)
}

func TestValidateHash(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.bin")
	
	data := []byte("test data")
	err := os.WriteFile(filePath, data, 0644)
	require.NoError(t, err)
	
	hash, err := CalculateHash(filePath)
	require.NoError(t, err)
	
	valid, err := ValidateHash(filePath, hash)
	require.NoError(t, err)
	require.True(t, valid)
	
	valid, err = ValidateHash(filePath, "invalid_hash")
	require.NoError(t, err)
	require.False(t, valid)
}

