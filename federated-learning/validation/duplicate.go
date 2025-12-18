package validation

import (
	"fmt"
	"math"
)

func CheckDuplicateGradients(gradients []float64, existingGradients [][]float64) (bool, error) {
	for _, existing := range existingGradients {
		if compareGradients(gradients, existing) {
			return true, nil
		}
	}
	return false, nil
}

func compareGradients(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	const epsilon = 1e-9
	for i := range a {
		if abs(a[i]-b[i]) > epsilon {
			return false
		}
	}
	return true
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func ValidateAggregation(gradients [][]float64) error {
	if len(gradients) == 0 {
		return fmt.Errorf("no gradients to validate")
	}
	
	dim := len(gradients[0])
	for i, g := range gradients {
		if len(g) != dim {
			return fmt.Errorf("gradient %d has inconsistent dimension", i)
		}
	}
	
	for i, g := range gradients {
		for j, val := range g {
			if math.IsNaN(val) || math.IsInf(val, 0) {
				return fmt.Errorf("invalid value at gradient %d, index %d", i, j)
			}
		}
	}
	
	return nil
}

