# LoRA (Low-Rank Adaptation)

LoRA implementation for efficient fine-tuning of large language models using low-rank matrix decomposition.

## Components

### Adapters (`adapters/`)
LoRA adapter structure and operations for applying low-rank adaptations to models.

**Key Functions:**
- `NewLoRAAdapter`: Create new LoRA adapter with specified rank and alpha
- `Apply`: Apply adapter to model (validates weights are initialized)
- `Save`: Save adapter weights to JSON file
- `Load`: Load adapter weights from JSON file
- `GetWeights`: Get current adapter weights
- `SetWeights`: Set adapter weights

**LoRA Structure:**
- `rank`: Rank of low-rank matrices (typically 4-64)
- `alpha`: Scaling factor for LoRA weights
- `dropout`: Dropout rate (default: 0.1)
- `targetModules`: List of modules to apply LoRA to (default: ["q_proj", "v_proj"])
- `weights`: Map of module names to weight arrays

**Weight Initialization:**
- Xavier initialization with scale = 1/rank
- Weight size per module: rank * rank * 2 (B and A matrices)
- Random initialization in range [-scale, scale]

### Training (`training/`)
LoRA training logic that integrates with Python ML frameworks.

**Key Functions:**
- `NewLoRATrainer`: Create LoRA trainer with adapter
- `SetWorkDir`: Set working directory for training scripts
- `Train`: Execute LoRA training via Python script
- `GetAdapterWeights`: Get current adapter weights
- `SetAdapterWeights`: Set adapter weights
- `simulateTraining`: Fallback simulated training if Python fails

**Training Flow:**
1. Create training directory
2. Save adapter configuration to JSON
3. Create Python training script (`train_lora.py`)
4. Execute Python script (python3 or python)
5. Load updated weights from `adapter_weights.json`
6. Fallback to simulated training if script fails

**Python Script:**
- Loads adapter configuration
- Performs LoRA training loop
- Updates LoRA weights (B and A matrices)
- Saves updated weights to `adapter_weights.json`

### Integration (`training/integration.go`)
Integration with Federated Learning for distributed LoRA training.

**Key Functions:**
- `NewLoRAFLIntegration`: Create LoRA-FL integration
- `TrainRound`: Train LoRA adapter for one FL round
- `UpdateAdapter`: Update adapter with aggregated weights
- `SaveAdapter`: Save adapter to file
- `LoadAdapter`: Load adapter from file

**Federated LoRA Flow:**
1. Each node trains LoRA adapter on local data
2. Extract LoRA weights (B and A matrices)
3. Aggregate weights across nodes
4. Update adapter with aggregated weights
5. Repeat for next round

## LoRA Algorithm

Low-rank decomposition of weight updates:

```
W = W0 + ΔW
ΔW = B * A * (alpha / rank)
```

Where:
- `W0`: Original weight matrix
- `B`: Low-rank matrix (rank x hidden_dim)
- `A`: Low-rank matrix (hidden_dim x rank)
- `alpha`: Scaling factor
- `rank`: Rank of decomposition

**Advantages:**
- Only train rank * rank * 2 parameters per module
- Much smaller than full fine-tuning
- Can be merged into base model after training

## Usage Example

```go
adapter := adapters.NewLoRAAdapter(8, 16.0)

err := adapter.Apply(model)
if err != nil {
    return err
}

trainer := training.NewLoRATrainer(adapter)
trainer.SetWorkDir("/tmp/lora")

err = trainer.Train(ctx, "/path/to/dataset")
if err != nil {
    return err
}

weights := trainer.GetAdapterWeights()

integration := training.NewLoRAFLIntegration(8, 16.0)
weights, err := integration.TrainRound(ctx, "/path/to/dataset")
if err != nil {
    return err
}

integration.UpdateAdapter(aggregatedWeights)
```

## Testing

Run tests:
```bash
cd lora
go test ./... -v
```

Test coverage:
- Adapters: Creation, save/load, apply validation
- Training: Weight management, simulated training
- Integration: FL integration, adapter updates

