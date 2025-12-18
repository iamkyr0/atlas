# Atlas Compute Node

The Atlas compute node is responsible for executing training tasks, managing resources, monitoring health, and participating in the decentralized AI training network.

## Components

### Executor (`executor/`)
Manages task execution lifecycle and coordinates with training executors.

**Key Files:**
- `task.go`: Task management, lifecycle, and execution coordination
- `training.go`: Training task execution using Python scripts
- `inference.go`: Inference task execution for model serving

**Key Functions:**
- `NewExecutor`: Create new executor with resource manager
- `AddTask`: Add a new task to the executor
- `GetTask`: Retrieve task by ID
- `ListTasks`: List all tasks
- `Start`: Start task processing loop
- `Stop`: Stop executor and cancel all tasks
- `StopTask`: Stop a specific task
- `SetWorkDir`: Set working directory for tasks
- `SetIPFSAPIURL`: Set IPFS API URL for downloads
- `InitializeTrainingExecutor`: Initialize training executor
- `InitializeInferenceExecutor`: Initialize inference executor

**Task Types:**
- `TaskTypeTraining`: Training tasks that execute Python scripts
- `TaskTypeInference`: Inference tasks for model serving

**Task Lifecycle:**
1. `pending`: Task created, waiting to start
2. `in_progress`: Task currently executing
3. `completed`: Task finished successfully
4. `failed`: Task failed with error
5. `paused`: Task stopped/cancelled

### Resource Manager (`resource/`)
Auto-detects and manages system resources (CPU, GPU, RAM, storage, network).

**Key Functions:**
- `NewManager`: Create resource manager with auto-detection
- `DetectResources`: Perform full resource detection including network tests
- `GetResources`: Get resource information as map
- `AllocateResources`: Allocate CPU and memory for a task
- `ReleaseResources`: Release allocated resources

**Auto-Detection:**
- CPU cores: Uses `runtime.NumCPU()`
- Memory: Uses `syscall.Sysinfo` (Linux) or fallback
- Storage: Uses `syscall.Statfs` or `df` command
- GPUs: Detects NVIDIA GPUs via `nvidia-smi`
- Network: Speed test using Cloudflare CDN
- Geolocation: IP and location via ip-api.com

**Resource Tracking:**
- Tracks allocated CPU and memory per task
- Prevents overallocation
- Validates resource availability before allocation

### Validator (`validator/`)
Validates shard and task assignments by querying the blockchain.

**Key Functions:**
- `NewValidator`: Create validator with blockchain client
- `ValidateAssignment`: Validate shard can be assigned to node
- `CheckDuplication`: Check if shard already assigned

**Validation Checks:**
- Shard not already assigned to another node
- Node has available capacity
- Node has positive reputation

**Blockchain Client:**
- `BlockchainClient` interface for blockchain queries
- `HTTPBlockchainClient` placeholder implementation
- Requires gRPC client with protobuf stubs for full implementation

### Health Monitor (`health/`)
Monitors node health and sends heartbeats.

**Key Functions:**
- `NewMonitor`: Create health monitor
- `Start`: Start heartbeat loop
- `SendHeartbeat`: Send heartbeat to blockchain

**Heartbeat:**
- Sends heartbeat every 30 seconds
- Updates node status on blockchain
- Used by blockchain to detect offline nodes

### Recovery (`recovery/`)
Handles checkpoint management and rollback operations.

**Key Components:**
- `checkpoint.go`: Checkpoint save/load operations
- `checkpoint_manager.go`: Checkpoint manager with IPFS integration
- `rollback_handler.go`: Task rollback and cleanup
- `emergency.go`: Graceful shutdown with emergency checkpoint

**Key Functions:**
- `SaveCheckpoint`: Save checkpoint to IPFS
- `LoadCheckpoint`: Load checkpoint from IPFS
- `ValidateCheckpoint`: Validate checkpoint signature and age
- `HandleRollback`: Handle task rollback and notify blockchain
- `CleanupTaskState`: Clean up task directory after rollback
- `SetupGracefulShutdown`: Setup signal handler for emergency checkpoint

**Checkpoint Structure:**
- Task ID, epoch, iteration
- IPFS CID of checkpoint data
- Timestamp and signature for validation
- Maximum age: 7 days

### Proof (`proof/`)
Generates proof of computation for task verification.

**Key Functions:**
- `GenerateProof`: Generate proof of computation
- `VerifyProof`: Verify proof integrity

**Proof Structure:**
- Task ID and node ID
- Timestamp
- Hash of computation data
- Computation metrics (iterations, time, memory, GPU)

### Network (`network/`)
Network speed testing and geolocation detection.

**Key Functions:**
- `SpeedTest`: Perform network speed test
- `GetGeolocation`: Get IP and geographic location

**Speed Test:**
- Download speed: Uses Cloudflare CDN test file
- Upload speed: Estimated as 10% of download
- Latency: Tests latency to Google

**Geolocation:**
- Public IP: Uses ipify.org
- Location: Uses ip-api.com (free tier)
- Returns: IP, country, region, city, coordinates

## CLI Commands

### Start Node
```bash
atlas-node start --chain-rpc http://localhost:26657 --ipfs-api /ip4/127.0.0.1/tcp/5001
```

### Show Status
```bash
atlas-node status
```

### Register Node
```bash
atlas-node register --node-id node-1 --address cosmos1abc123
```

### Show Config
```bash
atlas-node config
```

## Configuration

**Flags:**
- `--chain-rpc`: Blockchain RPC URL (default: http://localhost:26657)
- `--ipfs-api`: IPFS API URL (default: /ip4/127.0.0.1/tcp/5001)
- `--node-id`: Node identifier
- `--address`: Node wallet address

## Task Execution Flow

### Training Tasks
1. **Task Received**: Task assigned via IPFS pub/sub or blockchain
2. **Resource Allocation**: Allocate CPU and memory
3. **Download Data**: Download model and dataset from IPFS
4. **Execute Training**: Run Python training script
5. **Save Checkpoint**: Periodically save checkpoints to IPFS
6. **Upload Results**: Upload gradients/results to IPFS
7. **Release Resources**: Free allocated resources
8. **Update Status**: Update task status on blockchain

### Inference Tasks
1. **Task Received**: Inference request received
2. **Model Download**: Download model from IPFS (cached if available)
3. **Input Preparation**: Prepare input data for model
4. **Execute Inference**: Run Python inference script
5. **Return Results**: Return inference results with latency
6. **Update Status**: Update task status on blockchain

## Error Handling

- **Task Failure**: Task marked as failed, resources released
- **Node Offline**: Tasks rolled back and reassigned
- **Resource Exhaustion**: Task rejected if insufficient resources
- **Network Failure**: Retry with exponential backoff

## Testing

Run tests:
```bash
cd node
go test ./... -v
```

Test coverage:
- Resource manager: Allocation, release, detection
- Validator: Assignment validation, duplication check
- Proof: Generation and verification
- Recovery: Rollback and cleanup

## Integration

**With Blockchain:**
- Registers node on blockchain
- Sends heartbeats for health monitoring
- Updates task status
- Queries for assignments

**With IPFS:**
- Downloads models and datasets
- Uploads checkpoints and results
- Uses pub/sub for task assignments
- Publishes rollback events

**With Python:**
- Executes Python training scripts
- Executes Python inference scripts
- Passes model and dataset paths
- Reads gradients from JSON output (training)
- Reads inference results from JSON (inference)
- Supports python3 and python commands
- Supports multiple model formats (PyTorch, ONNX, TensorFlow)

