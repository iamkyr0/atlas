# Atlas Chain Modules

Atlas blockchain is built on Cosmos SDK and consists of 9 custom modules that handle compute node management, training job coordination, storage, rewards, models, health monitoring, recovery, sharding, and validation.

## Module Structure

### x/compute
Manages compute node registration, heartbeat tracking, and reputation system.

**Key Components:**
- `keeper/keeper.go`: Node CRUD operations, iteration
- `keeper/reputation.go`: Reputation calculation and updates
- `keeper/grpc_query.go`: gRPC query server for node queries
- `keeper/msg_server.go`: Message server for node registration and heartbeat
- `types/node.go`: Node type definitions

**Key Functions:**
- `RegisterNode`: Register a new compute node
- `GetNode`: Retrieve node by ID
- `GetAllNodes`: List all registered nodes
- `IterateNodes`: Iterate through nodes with handler
- `UpdateReputation`: Update node reputation based on uptime
- `GetNodeReputation`: Get current reputation score
- `UpdateHeartbeat`: Update node heartbeat timestamp

### x/training
Manages training jobs and tasks, coordinates federated learning workflows.

**Key Components:**
- `keeper/keeper.go`: Job and task CRUD operations
- `keeper/gradient.go`: Gradient contribution tracking and fair reward calculation
- `keeper/grpc_query.go`: gRPC query server for job/task queries
- `keeper/msg_server.go`: Message server for job submission and task management
- `types/job.go`: Job type definitions
- `types/task.go`: Task type definitions

**Key Functions:**
- `SubmitJob`: Create a new training job
- `GetJob`: Retrieve job by ID
- `SetJob`: Store job in state
- `CreateTask`: Create a task for a job
- `GetTask`: Retrieve task by ID
- `SetTask`: Store task in state
- `IterateTasks`: Iterate through tasks with handler
- `TrackGradientContribution`: Track gradient contributions for fair rewards
- `GetGradientContributions`: Get contributions for a job round
- `CalculateFairRewards`: Calculate proportional rewards based on contributions

### x/storage
Manages storage node registration and capacity tracking.

**Key Components:**
- `keeper/keeper.go`: Storage node management
- `types/keys.go`: Store key definitions

**Key Functions:**
- `RegisterStorageNode`: Register a new storage node
- `GetStorageNode`: Retrieve storage node by ID
- `GetAllStorageNodes`: List all storage nodes
- `UpdateStorageNodeCapacity`: Update used capacity for a node
- `GetAvailableStorageNodes`: Get nodes with available capacity

### x/reward
Handles reward calculation and distribution to compute nodes.

**Key Components:**
- `keeper/distribution.go`: Reward calculation and distribution logic
- `keeper/keeper.go`: Keeper structure with dependencies

**Key Functions:**
- `CalculateReward`: Calculate reward based on work completed and reputation
- `DistributeReward`: Send reward tokens to node address

**Reward Formula:**
- Base reward multiplied by work completed (0.0-1.0)
- Adjusted by reputation multiplier (0.0-1.0)
- Reputation based on uptime percentage

### x/model
Manages AI model registry and versioning.

**Key Components:**
- `keeper/registry.go`: Model registration and versioning
- `keeper/grpc_query.go`: gRPC query server for model queries
- `keeper/msg_server.go`: Message server for model registration
- `types/model.go`: Model type definitions

**Key Functions:**
- `RegisterModel`: Register a new model version
- `GetModel`: Retrieve model by ID
- `GetAllModels`: List all registered models
- `UpdateModelVersion`: Update model version and CID
- `GetModelsByCID`: Find models by IPFS CID

### x/health
Monitors node health status based on heartbeat timestamps.

**Key Components:**
- `keeper/monitor.go`: Health checking logic
- `keeper/keeper.go`: Keeper structure

**Key Functions:**
- `CheckNodeHealth`: Check if node is healthy (heartbeat within timeout)
- `UpdateHeartbeat`: Update node heartbeat timestamp
- `GetOfflineNodes`: Get list of nodes that haven't sent heartbeat

**Health Check:**
- Heartbeat timeout: 90 seconds
- Nodes without heartbeat within timeout are marked offline

### x/recovery
Handles task rollback and reassignment when nodes go offline.

**Key Components:**
- `keeper/recovery.go`: Rollback and reassignment logic
- `keeper/keeper.go`: Keeper structure with dependencies

**Key Functions:**
- `RollbackTasksForNode`: Rollback all in-progress/assigned tasks for a node
- `ReassignTask`: Reassign a task to a new node
- `HandleNodeOffline`: Complete recovery workflow when node goes offline

**Recovery Flow:**
1. Rollback all tasks assigned to offline node
2. Find available healthy nodes
3. Reassign rolled-back tasks to available nodes

### x/sharding
Manages data and model sharding for distributed training.

**Key Components:**
- `keeper/shard.go`: Shard management operations
- `keeper/keeper.go`: Keeper structure

**Key Functions:**
- `RegisterShard`: Register a new shard
- `GetShard`: Retrieve shard by ID
- `AssignShardToNode`: Assign shard to a compute node
- `GetShardsForJob`: Get all shards for a job
- `GetShardsByNode`: Get all shards assigned to a node
- `GetShardsByHash`: Find shards by content hash (for deduplication)
- `GetAllShards`: List all registered shards

**Shard Structure:**
- ID: Unique shard identifier
- JobID: Associated training job
- CID: IPFS content identifier
- Hash: Content hash for deduplication
- NodeID: Assigned compute node
- Status: Current status (pending, assigned, completed)
- Size: Shard size in bytes

### x/validation
Validates shard and task assignments, checks for duplicates.

**Key Components:**
- `keeper/validator.go`: Validation logic
- `keeper/keeper.go`: Keeper structure with dependencies

**Key Functions:**
- `ValidateShardAssignment`: Validate shard can be assigned to node
- `CheckDuplicateShard`: Check if shard content already exists (by hash)
- `ValidateTaskAssignment`: Validate task can be assigned to node

**Validation Checks:**
- Shard not already assigned to another node
- No duplicate shard content (same hash)
- Node exists and is online
- Node is healthy (heartbeat check)

## Module Dependencies

```
compute (no dependencies)
  ↓
training → compute, storage
  ↓
reward → compute, storage
  ↓
model (no dependencies)
  ↓
health → compute
  ↓
recovery → training, compute, health
  ↓
sharding (no dependencies)
  ↓
validation → sharding, training, compute, health
  ↓
storage (no dependencies)
```

## State Storage

All modules use Cosmos SDK KVStore for persistent state:
- Jobs: `job:{jobID}`
- Tasks: `task:{taskID}`
- Nodes: `{nodeID}`
- Models: `model:{modelID}`
- Shards: `shard:{shardID}`
- Gradients: `gradient:{jobID}:{nodeID}:{round}:{gradientCID}`

## gRPC Services

Three modules expose gRPC services:

### Training Module
- `GetJob`: Get job by ID
- `ListJobs`: List all jobs
- `GetTask`: Get task by ID
- `GetTasksByJob`: Get all tasks for a job

### Compute Module
- `GetNode`: Get node by ID
- `ListNodes`: List all nodes

### Model Module
- `GetModel`: Get model by ID
- `ListModels`: List all models

## Message Handlers

### Training Module
- `MsgSubmitJob`: Create a new training job
- `MsgCreateTask`: Create a task for a job
- `MsgUpdateTaskStatus`: Update task status and progress

### Compute Module
- `MsgRegisterNode`: Register a new compute node
- `MsgUpdateHeartbeat`: Update node heartbeat

### Model Module
- `MsgRegisterModel`: Register a new model version

## Testing

All keepers have comprehensive unit tests:
- `*_test.go` files for each keeper
- Tests cover CRUD operations, edge cases, and error handling
- Use in-memory stores for fast test execution

Run tests:
```bash
cd chain
go test ./x/... -v
```

## Integration

Modules are integrated in `app.go`:
- All keepers initialized with proper dependencies
- gRPC query servers registered
- Message servers registered
- Genesis initialization for each module

## Usage Example

```go
keeper := trainingkeeper.NewKeeper(cdc, storeKey, memKey, computeKeeper, storageKeeper, bankKeeper)

job := types.Job{
    ID: "job-1",
    ModelID: "model-1",
    DatasetCID: "QmABC123",
    Status: types.TaskStatusPending,
}
keeper.SetJob(ctx, job)

retrievedJob, found := keeper.GetJob(ctx, "job-1")
```

