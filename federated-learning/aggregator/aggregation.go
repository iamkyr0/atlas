package aggregator

import (
	"crypto/rand"
	"fmt"
	"math"
)

func FederatedAveraging(gradientsList [][]float64, weights []float64) ([]float64, error) {
	if len(gradientsList) == 0 {
		return nil, fmt.Errorf("no gradients to aggregate")
	}

	if len(weights) == 0 {
		weights = make([]float64, len(gradientsList))
		equalWeight := 1.0 / float64(len(gradientsList))
		for i := range weights {
			weights[i] = equalWeight
		}
	}

	if len(gradientsList) != len(weights) {
		return nil, fmt.Errorf("gradients and weights length mismatch")
	}

	totalWeight := 0.0
	for _, w := range weights {
		totalWeight += w
	}
	for i := range weights {
		weights[i] /= totalWeight
	}

	dim := len(gradientsList[0])
	for _, g := range gradientsList {
		if len(g) != dim {
			return nil, fmt.Errorf("gradient dimension mismatch")
		}
	}

	aggregated := make([]float64, dim)
	for i := 0; i < dim; i++ {
		for j, g := range gradientsList {
			aggregated[i] += g[i] * weights[j]
		}
	}

	return aggregated, nil
}

func SecureAggregation(gradientsList [][]float64, noiseScale float64) ([]float64, error) {
	aggregated, err := FederatedAveraging(gradientsList, nil)
	if err != nil {
		return nil, err
	}

	for i := range aggregated {
		noise := laplacianNoise(noiseScale)
		aggregated[i] += noise
	}

	return aggregated, nil
}

func laplacianNoise(scale float64) float64 {
	buf := make([]byte, 2)
	if _, err := rand.Read(buf); err != nil {
		return 0.0
	}
	
	u1 := float64(buf[0]) / 256.0
	u2 := float64(buf[1]) / 256.0
	
	u := u1 - 0.5
	if u == 0 {
		u = u2 / 256.0 - 0.5
	}
	
	noise := scale * math.Copysign(1.0, u) * math.Log(1.0 - 2.0*math.Abs(u))
	return noise
}

