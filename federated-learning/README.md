# Federated Learning

Federated Learning implementation for distributed model training across multiple nodes without centralizing data.

## Components

### Client (`client/`)
Client-side federated learning implementation that performs local training on data shards.

**Key Functions:**
- `NewFLClient`: Create new FL client
- `Train`: Perform training on a shard and extract gradients
- `SendGradients`: Send gradients to aggregator via IPFS pub/sub
- `ReceiveModel`: Receive aggregated model updates
- `SetKeepWorkDir`: Control whether to preserve training directory
- `GetTrainScriptPath`: Get expected path of training script

**Training Flow:**
1. Download model and shard data from IPFS
2. Create Python training script (`train.py`)
3. Execute training script
4. Read gradients from `gradients.json`
5. Send gradients to aggregator

**Python Script:**
- Loads model and dataset
- Performs training loop
- Extracts gradients from model parameters
- Saves gradients to `gradients.json`

### Aggregator (`aggregator/`)
Distributed aggregator that collects gradients and performs federated averaging.

**Key Functions:**
- `NewAggregator`: Create new aggregator instance
- `BecomeAggregator`: Announce this node as aggregator
- `Start`: Start listening for gradients via IPFS pub/sub
- `ReceiveGradients`: Receive gradients from clients
- `Aggregate`: Perform federated averaging

**Aggregation Flow:**
1. Subscribe to gradient topic via IPFS pub/sub
2. Collect gradients from all clients
3. Perform federated averaging
4. Broadcast aggregated model to clients
5. Increment round counter

**Algorithms:**
- `FederatedAveraging`: Standard federated averaging with weighted aggregation
- `SecureAggregation`: Secure aggregation with differential privacy (Laplacian noise)

### Protocols (`protocols/`)
IPFS pub/sub communication protocol for FL.

**Key Functions:**
- `NewFLProtocol`: Create FL protocol instance
- `SendGradients`: Send gradients to aggregator
- `ReceiveModel`: Receive model updates from aggregator
- `GetAPI`: Get underlying IPFS API for advanced operations

**Topics:**
- `/atlas/fl/gradients/{jobID}`: Gradient messages from clients
- `/atlas/fl/model/{jobID}`: Aggregated model updates
- `/atlas/fl/aggregator/{jobID}`: Aggregator announcements

### Sharding (`sharding/`)
Coordinates shard assignment and tracks completion status.

**Key Functions:**
- `NewCoordinator`: Create shard coordinator
- `AssignShard`: Assign shard to a node
- `WaitForShards`: Wait for all shards to complete
- `UpdateShardStatus`: Update shard status and progress
- `GetShardStatus`: Get current shard status

**Shard States:**
- `assigned`: Shard assigned to node
- `in_progress`: Node processing shard
- `completed`: Shard processing complete
- `failed`: Shard processing failed

### Validation (`validation/`)
Validates gradients and aggregation results.

**Key Functions:**
- `CheckDuplicateGradients`: Check for duplicate gradient submissions
- `ValidateAggregation`: Validate aggregated gradients for anomalies

**Validation Checks:**
- Dimension consistency across gradients
- No NaN or Inf values
- Duplicate detection using epsilon comparison

## Federated Averaging Algorithm

Weighted average of gradients from multiple clients:

```
aggregated[i] = Σ(gradients[j][i] * weights[j]) / Σ(weights[j])
```

Where:
- `gradients[j]` is gradient vector from client j
- `weights[j]` is weight for client j (default: equal weights)
- Result is normalized by total weight

## Secure Aggregation

Adds Laplacian noise for differential privacy:

```
noise = Laplacian(0, scale)
aggregated[i] = federated_average[i] + noise
```

Where:
- `scale` is noise scale parameter
- Larger scale = more privacy, less accuracy

## Usage Example

```go
client := client.NewFLClient("node-1", "/ip4/127.0.0.1/tcp/5001", "/tmp/fl")

gradients, err := client.Train(ctx, "shard-1", "QmModel123")
if err != nil {
    return err
}

err = client.SendGradients(ctx, "job-1", 1, gradients)
if err != nil {
    return err
}

aggregator := aggregator.NewAggregator("job-1", "node-1", "/ip4/127.0.0.1/tcp/5001")
aggregator.BecomeAggregator()
aggregator.Start(ctx)

aggregated, err := aggregator.Aggregate()
if err != nil {
    return err
}
```

## Testing

Run tests:
```bash
cd federated-learning
go test ./... -v
```

Test coverage:
- Aggregation: Federated averaging, secure aggregation
- Validation: Duplicate detection, anomaly detection
- Sharding: Coordinator operations, wait for completion

