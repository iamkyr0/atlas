package training

import (
	"context"
	"fmt"
	"github.com/atlas/lora/adapters"
)

type LoRAFLIntegration struct {
	adapter *adapters.LoRAAdapter
	trainer *LoRATrainer
}

func NewLoRAFLIntegration(rank int, alpha float64) *LoRAFLIntegration {
	adapter := adapters.NewLoRAAdapter(rank, alpha)
	trainer := NewLoRATrainer(adapter)
	
	return &LoRAFLIntegration{
		adapter: adapter,
		trainer: trainer,
	}
}

func (l *LoRAFLIntegration) TrainRound(ctx context.Context, datasetPath string) (map[string][]float64, error) {
	if err := l.trainer.Train(ctx, datasetPath); err != nil {
		return nil, fmt.Errorf("training failed: %w", err)
	}

	weights, err := l.trainer.GetAdapterWeights()
	if err != nil {
		return nil, fmt.Errorf("failed to get weights: %w", err)
	}

	return weights, nil
}

func (l *LoRAFLIntegration) UpdateAdapter(aggregatedWeights map[string][]float64) error {
	l.adapter.SetWeights(aggregatedWeights)
	return nil
}

func (l *LoRAFLIntegration) SaveAdapter(path string) error {
	return l.adapter.Save(path)
}

func (l *LoRAFLIntegration) LoadAdapter(path string) error {
	return l.adapter.Load(path)
}

