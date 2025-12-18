# Storage Layer

Storage layer provides IPFS integration, data sharding, and content validation for the Atlas platform.

## Components

### Manager (`manager/`)
IPFS manager for file operations with fallback support.

**Key Functions:**
- `NewIPFSManager`: Create IPFS manager with optional fallback nodes
- `AddFile`: Upload file to IPFS and get CID
- `GetFile`: Download file from IPFS by CID
- `GetFileWithFallback`: Download with fallback to other IPFS nodes
- `CheckDataAvailability`: Check if CID is available on IPFS network
- `Pin`: Pin content to prevent garbage collection
- `SetFallbackNodes`: Configure fallback IPFS nodes

**Fallback Mechanism:**
- Primary IPFS node first
- Fallback to secondary nodes with exponential backoff
- Retry delay: 2 seconds * 2^attempt

### Sharding (`sharding/`)
Data and model sharding for distributed processing.

**Key Functions:**
- `SplitDataset`: Split dataset into multiple shards
- `SplitModel`: Split model using specified strategy
- `CalculateShardHash`: Calculate SHA-256 hash of shard
- `SplitModelByLayers`: Split model by layers (placeholder)
- `SplitModelByChunks`: Split model into chunks

**Sharding Strategies:**
- `chunk`: Split into fixed-size chunks (4 chunks)
- `layer`: Split by model layers (placeholder)
- `default`: Upload as single file

**ShardManager:**
- `SplitDataset`: Split dataset into N shards
- `SplitModelByChunks`: Split model into 4 chunks
- `SplitModelByLayers`: Split model by layers
- `ValidateShardHash`: Validate shard integrity

### Validation (`validation/`)
Content validation and hash calculation.

**Key Functions:**
- `CalculateHash`: Calculate SHA-256 hash of file
- `ValidateHash`: Validate file hash matches expected

**Hash Algorithm:**
- SHA-256
- Returns hex-encoded string (64 characters)

### PubSub (`pubsub/`)
IPFS pub/sub messaging for real-time communication.

**Key Functions:**
- `NewEventSubscriber`: Create event subscriber
- `Subscribe`: Subscribe to topic with handler
- `PublishEvent`: Publish event to topic

**Topics:**
- `/atlas/tasks/assign/{node_id}`: Task assignments
- `/atlas/fl/gradients/{job_id}`: Federated learning gradients
- `/atlas/fl/model/{job_id}`: Aggregated models
- `/atlas/recovery/rollback/{job_id}`: Rollback events

## Usage Example

```go
ipfsManager := manager.NewIPFSManager("/ip4/127.0.0.1/tcp/5001", "/ip4/192.168.1.1/tcp/5001")

cid, err := ipfsManager.AddFile("/path/to/file")
if err != nil {
    return err
}

err = ipfsManager.GetFile(cid, "/path/to/output")
if err != nil {
    return err
}

cids, err := sharding.SplitDataset(ipfsManager, "/path/to/dataset", 4, "/tmp/shards")
if err != nil {
    return err
}

hash, err := validation.CalculateHash("/path/to/shard")
if err != nil {
    return err
}
```

## Testing

Run tests:
```bash
cd storage
go test ./... -v
```

Test coverage:
- Manager: File operations, fallback mechanism
- Sharding: Dataset splitting, model splitting, hash calculation
- Validation: Hash calculation and validation

