package adapters

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

type LoRAAdapter struct {
	rank         int
	alpha        float64
	dropout      float64
	targetModules []string
	weights      map[string][]float64
}

func NewLoRAAdapter(rank int, alpha float64) *LoRAAdapter {
	adapter := &LoRAAdapter{
		rank:         rank,
		alpha:        alpha,
		dropout:      0.1,
		targetModules: []string{"q_proj", "v_proj"},
		weights:      make(map[string][]float64),
	}
	
	for _, module := range adapter.targetModules {
		weightSize := rank * rank * 2
		weights := make([]float64, weightSize)
		
		scale := 1.0 / float64(rank)
		for i := range weights {
			weights[i] = (rand.Float64() - 0.5) * 2.0 * scale
		}
		
		adapter.weights[module] = weights
	}
	
	return adapter
}

func (a *LoRAAdapter) Apply(model interface{}) error {
	if model == nil {
		return fmt.Errorf("model cannot be nil")
	}
	
	if len(a.weights) == 0 {
		return fmt.Errorf("LoRA adapter weights not initialized")
	}
	
	for _, module := range a.targetModules {
		if _, ok := a.weights[module]; !ok {
			return fmt.Errorf("weights not found for target module: %s", module)
		}
	}
	
	return nil
}

func (a *LoRAAdapter) Save(path string) error {
	data := map[string]interface{}{
		"rank":          a.rank,
		"alpha":         a.alpha,
		"dropout":       a.dropout,
		"target_modules": a.targetModules,
		"weights":       a.weights,
	}
	
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal adapter: %w", err)
	}
	
	if err := os.WriteFile(path, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write adapter file: %w", err)
	}
	
	return nil
}

func (a *LoRAAdapter) Load(path string) error {
	jsonData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read adapter file: %w", err)
	}
	
	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return fmt.Errorf("failed to unmarshal adapter: %w", err)
	}
	
	if weights, ok := data["weights"].(map[string]interface{}); ok {
		a.weights = make(map[string][]float64)
		for k, v := range weights {
			if weightsList, ok := v.([]interface{}); ok {
				floats := make([]float64, len(weightsList))
				for i, w := range weightsList {
					if f, ok := w.(float64); ok {
						floats[i] = f
					}
				}
				a.weights[k] = floats
			}
		}
	}
	
	return nil
}

func (a *LoRAAdapter) GetWeights() map[string][]float64 {
	return a.weights
}

func (a *LoRAAdapter) SetWeights(weights map[string][]float64) {
	a.weights = weights
}

