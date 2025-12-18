package aggregator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFederatedAveraging(t *testing.T) {
	gradients1 := []float64{1.0, 2.0, 3.0}
	gradients2 := []float64{2.0, 3.0, 4.0}
	gradients3 := []float64{3.0, 4.0, 5.0}
	
	gradientsList := [][]float64{gradients1, gradients2, gradients3}
	
	result, err := FederatedAveraging(gradientsList, nil)
	require.NoError(t, err)
	require.Len(t, result, 3)
	require.InDelta(t, 2.0, result[0], 0.01)
	require.InDelta(t, 3.0, result[1], 0.01)
	require.InDelta(t, 4.0, result[2], 0.01)
	
	weights := []float64{0.5, 0.3, 0.2}
	result, err = FederatedAveraging(gradientsList, weights)
	require.NoError(t, err)
	require.Len(t, result, 3)
	
	_, err = FederatedAveraging([][]float64{}, nil)
	require.Error(t, err)
	
	_, err = FederatedAveraging(gradientsList, []float64{0.5, 0.3})
	require.Error(t, err)
	
	mismatched := [][]float64{gradients1, []float64{1.0, 2.0}}
	_, err = FederatedAveraging(mismatched, nil)
	require.Error(t, err)
}

func TestSecureAggregation(t *testing.T) {
	gradients1 := []float64{1.0, 2.0, 3.0}
	gradients2 := []float64{2.0, 3.0, 4.0}
	
	gradientsList := [][]float64{gradients1, gradients2}
	
	result, err := SecureAggregation(gradientsList, 0.1)
	require.NoError(t, err)
	require.Len(t, result, 3)
	
	_, err = SecureAggregation([][]float64{}, 0.1)
	require.Error(t, err)
}

func TestLaplacianNoise(t *testing.T) {
	noise1 := laplacianNoise(1.0)
	noise2 := laplacianNoise(1.0)
	
	require.NotEqual(t, noise1, noise2)
}

