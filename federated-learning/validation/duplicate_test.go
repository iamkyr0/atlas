package validation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckDuplicateGradients(t *testing.T) {
	gradients := []float64{1.0, 2.0, 3.0}
	existing := [][]float64{
		{1.0, 2.0, 3.0},
		{4.0, 5.0, 6.0},
	}
	
	duplicate, err := CheckDuplicateGradients(gradients, existing)
	require.NoError(t, err)
	require.True(t, duplicate)
	
	gradients = []float64{7.0, 8.0, 9.0}
	duplicate, err = CheckDuplicateGradients(gradients, existing)
	require.NoError(t, err)
	require.False(t, duplicate)
}

func TestValidateAggregation(t *testing.T) {
	gradients := [][]float64{
		{1.0, 2.0, 3.0},
		{4.0, 5.0, 6.0},
		{7.0, 8.0, 9.0},
	}
	
	err := ValidateAggregation(gradients)
	require.NoError(t, err)
	
	err = ValidateAggregation([][]float64{})
	require.Error(t, err)
	
	mismatched := [][]float64{
		{1.0, 2.0, 3.0},
		{4.0, 5.0},
	}
	err = ValidateAggregation(mismatched)
	require.Error(t, err)
	
	withNaN := [][]float64{
		{1.0, 2.0, 3.0},
		{4.0, math.NaN(), 6.0},
	}
	err = ValidateAggregation(withNaN)
	require.Error(t, err)
}

