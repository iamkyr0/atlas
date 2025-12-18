# Atlas Architecture

Decentralized AI Platform Architecture Overview

## System Overview

Atlas is a decentralized Infrastructure-as-a-Service (IaaS) platform for AI model fine-tuning and serving. It uses blockchain for coordination and rewards, with distributed compute and storage nodes.

## Core Components

### Blockchain Layer (`chain/`)
Cosmos SDK blockchain with custom modules:
- **Compute Module**: Node registration and management
- **Training Module**: Job and task coordination
- **Storage Module**: Storage node management
- **Reward Module**: Reward calculation and distribution
- **Model Module**: Model registry and versioning
- **Health Module**: Node health monitoring
- **Recovery Module**: Task rollback and reassignment
- **Sharding Module**: Data shard management
- **Validation Module**: Assignment and content validation

### Compute Node (`node/`)
Distributed compute nodes that execute training tasks:
- **Executor**: Task execution engine
- **Resource Manager**: CPU/GPU/RAM/Storage detection
- **Health Monitor**: Heartbeat and status reporting
- **Recovery**: Checkpoint management and rollback
- **Proof**: Computation proof generation
- **Validator**: Assignment validation
- **Network**: Speed testing and geolocation

### Federated Learning (`federated-learning/`)
Privacy-preserving distributed training:
- **Client**: Client-side training logic
- **Aggregator**: Gradient aggregation server
- **Protocols**: IPFS pub/sub communication
- **Sharding**: Shard coordination
- **Validation**: Gradient validation

### LoRA Fine-Tuning (`lora/`)
Efficient fine-tuning with Low-Rank Adaptation:
- **Adapters**: LoRA adapter implementation
- **Training**: LoRA training logic

### Storage Layer (`storage/`)
IPFS-based decentralized storage:
- **Manager**: IPFS file operations with fallback
- **Sharding**: Dataset and model sharding
- **Validation**: Content hash validation
- **PubSub**: IPFS pub/sub messaging

### Client SDK (`sdk/`)
Python SDK for client interaction:
- **P2P Communication**: Direct node communication
- **Blockchain Client**: gRPC chain interaction
- **Job Management**: Job submission and monitoring
- **Model Serving**: Model loading and inference
- **CLI**: Command-line interface
- **Daemon**: REST API for application integration

## Data Flow

### Training Job Flow

1. **Client Submission**:
   - Client uploads dataset to IPFS
   - Client submits job via blockchain transaction
   - Blockchain creates job and tasks

2. **Task Assignment**:
   - Blockchain assigns tasks to available nodes
   - Nodes validate assignments
   - Nodes download model and shard from IPFS

3. **Training Execution**:
   - Nodes execute training scripts
   - Gradients uploaded to IPFS
   - Aggregator collects and aggregates gradients

4. **Model Update**:
   - Aggregated model uploaded to IPFS
   - Blockchain updates job status
   - Client receives updates via pub/sub

5. **Completion**:
   - Final model CID stored on blockchain
   - Rewards distributed to nodes
   - Client downloads trained model

### Node Discovery

- **Primary**: IPFS DHT for node discovery
- **Fallback**: Blockchain query for registered nodes
- **P2P**: Direct node-to-node communication

### Storage Strategy

- **Datasets**: Sharded and distributed across IPFS
- **Models**: Stored as single files or chunks
- **Checkpoints**: Periodic checkpoints uploaded to IPFS
- **Fallback**: Multiple IPFS nodes for redundancy

## Communication Protocols

### IPFS Pub/Sub Topics

- `/atlas/tasks/assign/{node_id}`: Task assignments
- `/atlas/fl/gradients/{job_id}`: Federated learning gradients
- `/atlas/fl/model/{job_id}`: Aggregated models
- `/atlas/recovery/rollback/{job_id}`: Rollback events
- `/atlas/fl/aggregator/{job_id}`: Aggregator announcements

### Blockchain Transactions

- `RegisterNode`: Node registration
- `SubmitJob`: Job submission
- `UpdateTaskStatus`: Task status updates
- `DistributeReward`: Reward distribution

## Security

- **Encryption**: Private datasets encrypted before IPFS upload
- **Private Networks**: Support for private IPFS networks
- **Proof of Computation**: Cryptographic proofs for task completion
- **Validation**: Content hash validation for shards
- **Reputation**: Node reputation system for trust

## Fault Tolerance

- **Health Monitoring**: Periodic heartbeat checks
- **Task Reassignment**: Automatic reassignment for offline nodes
- **Checkpoint Recovery**: Resume from checkpoints
- **Fallback Nodes**: Multiple IPFS nodes for redundancy
- **Rollback Mechanism**: Task rollback on failure

## Scalability

- **Sharding**: Large datasets split into shards
- **Distributed Training**: Multiple nodes train in parallel
- **Horizontal Scaling**: Add more nodes for capacity
- **Load Balancing**: Task distribution across nodes

