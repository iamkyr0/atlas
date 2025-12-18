# Node Usage Guide

Guide for running an Atlas compute node to provide compute and storage resources.

## Prerequisites

- Go 1.21+
- Docker
- IPFS node
- GPU (optional, recommended)
- Minimum 8GB RAM
- Minimum 100GB storage

## Installation

### 1. Setup IPFS Node

```bash
wget https://dist.ipfs.io/go-ipfs/v0.20.0/go-ipfs_v0.20.0_linux-amd64.tar.gz
tar -xvzf go-ipfs_v0.20.0_linux-amd64.tar.gz
cd go-ipfs
sudo ./install.sh

ipfs init
ipfs daemon
```

### 2. Build Node Binary

```bash
cd node
go mod download
go build -o atlas-node cmd/node/main.go
```

## Configuration

### Environment Variables

```bash
export CHAIN_RPC_URL=localhost:9090
export IPFS_API_URL=/ip4/127.0.0.1/tcp/5001
export NODE_ID=node-001
export NODE_ADDRESS=your_wallet_address
export WORK_DIR=/tmp/atlas-work
```

### Command Line Flags

```bash
atlas-node start \
  --chain-rpc localhost:9090 \
  --ipfs-api /ip4/127.0.0.1/tcp/5001 \
  --node-id node-001 \
  --address your_wallet_address
```

## Node Registration

### 1. Register on Blockchain

```bash
atlas-node register \
  --chain-rpc localhost:9090 \
  --node-id node-001 \
  --address your_wallet_address
```

Node will automatically detect:
- CPU cores
- GPU count and memory
- RAM and storage
- Network speed
- Geolocation (IP, country, region)

### 2. Verify Registration

```bash
atlasd query compute get-node node-001
```

## Running the Node

### Start Node

```bash
atlas-node start \
  --chain-rpc localhost:9090 \
  --ipfs-api /ip4/127.0.0.1/tcp/5001 \
  --node-id node-001 \
  --address your_wallet_address
```

### Node Operations

The node automatically:
- Sends heartbeat every 30 seconds
- Receives task assignments from blockchain
- Downloads models and datasets from IPFS
- Executes training tasks
- Uploads gradients and checkpoints to IPFS
- Receives rewards for completed tasks

### Check Node Status

```bash
atlas-node status
```

Output shows:
- CPU cores
- Memory and storage
- GPU count
- Network speed
- Location (IP, country, region)

## Monitoring

### Node Health

```bash
atlasd query health check-node-health node-001
```

### Task Status

```bash
atlasd query training get-tasks-by-node node-001
```

### Rewards

```bash
atlasd query reward node-rewards node-001
```

### Reputation

```bash
atlasd query compute node-reputation node-001
```

## Docker Deployment

### Build Image

```bash
docker build -t atlas-node:latest -f docker/Dockerfile.node .
```

### Run Container

```bash
docker run -d \
  --name atlas-node \
  --gpus all \
  -v /tmp/atlas-work:/work \
  -e CHAIN_RPC_URL=localhost:9090 \
  -e IPFS_API_URL=/ip4/127.0.0.1/tcp/5001 \
  -e NODE_ID=node-001 \
  atlas-node:latest
```

### View Logs

```bash
docker logs -f atlas-node
```

## Storage Node

To run as storage-only node:

```bash
atlasd tx storage register-storage-node \
  --node-id storage-node-001 \
  --storage-gb 1000 \
  --from mykey \
  --chain-id atlas-chain
```

## Troubleshooting

### Node Not Receiving Tasks

- Verify node is registered: `atlasd query compute get-node node-001`
- Check node status is "online"
- Verify heartbeat is being sent
- Check node has available resources

### Task Execution Failures

- Check IPFS connectivity: `ipfs id`
- Verify model/dataset CIDs are accessible
- Check disk space: `df -h`
- Review task logs in work directory

### Network Issues

- Verify IPFS node is accessible
- Check firewall rules
- Test IPFS API: `curl http://localhost:5001/api/v0/version`
- Verify chain gRPC connection

### Resource Detection

- GPU: Run `nvidia-smi` to verify GPU detection
- Storage: Check `df -h` output
- Network: Speed test runs automatically on start

## Best Practices

- Keep IPFS node running continuously
- Monitor disk space for checkpoints
- Set appropriate work directory size
- Use GPU for faster training
- Monitor node reputation
- Keep node software updated

