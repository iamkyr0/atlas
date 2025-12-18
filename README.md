# Atlas - Decentralized AI Platform

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![Python Version](https://img.shields.io/badge/Python-3.10+-blue.svg)](https://www.python.org/)

**Atlas** is a decentralized Infrastructure-as-a-Service (IaaS) platform for fine-tuning and serving AI models. The platform uses blockchain for coordination and reward distribution, with nodes contributing compute and storage resources similar to cryptocurrency mining systems.

## 🌟 Features

- **Decentralized Training**: Fine-tune AI models using distributed compute nodes
- **Federated Learning**: Privacy-preserving distributed training with gradient aggregation
- **LoRA Fine-Tuning**: Efficient fine-tuning with Low-Rank Adaptation
- **Model Serving**: Serve models for inference via HTTP/gRPC API
- **Inference Network**: Distributed network for LLM, Vision, Speech-to-text, and Embedding services
- **Blockchain Coordination**: Cosmos SDK blockchain for job coordination and reward distribution
- **IPFS Storage**: Decentralized storage for datasets, models, and checkpoints
- **P2P Communication**: Direct peer-to-peer communication without API Gateway
- **Auto Resource Detection**: Automatic detection of CPU, GPU, RAM, storage, and network speed
- **Fault Tolerance**: Automatic task reassignment, checkpoint recovery, and graceful degradation

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    ATLAS PLATFORM                            │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────┐      ┌──────────┐      ┌──────────┐        │
│  │  Client  │      │  Node 1  │      │  Node 2  │        │
│  │   SDK    │      │          │      │          │        │
│  └────┬─────┘      └────┬─────┘      └────┬─────┘        │
│       │                 │                 │               │
│       ├─────────────────┴─────────────────┤               │
│       │                                     │               │
│       │         IPFS Pub/Sub               │               │
│       │    (Real-time messaging)           │               │
│       │                                     │               │
│       ├─────────────────────────────────────┤               │
│       │                                     │               │
│       │         IPFS Storage                │               │
│       │    (Files: datasets, models)       │               │
│       │                                     │               │
│       └─────────────────┬───────────────────┘               │
│                        │                                    │
│                        │ gRPC/RPC                           │
│                        │                                    │
│                 ┌──────▼──────┐                            │
│                 │   Cosmos    │                            │
│                 │   Chain     │                            │
│                 │  (Blockchain)│                           │
│                 └─────────────┘                            │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Core Components

- **Blockchain Layer** (`chain/`): Cosmos SDK with 10 custom modules
- **Compute Node** (`node/`): Distributed compute nodes for training & inference
- **Federated Learning** (`federated-learning/`): Privacy-preserving distributed training
- **LoRA** (`lora/`): Efficient fine-tuning with Low-Rank Adaptation
- **Storage Layer** (`storage/`): IPFS-based decentralized storage
- **Client SDK** (`sdk/`): Python SDK with CLI and daemon mode

## 🚀 Quick Start

### Prerequisites

- **Go** 1.21+
- **Python** 3.10+
- **Docker** (optional, for containerized deployment)
- **IPFS node** (or use Infura/IPFS gateway)
- **Cosmos SDK** (for blockchain)

### Installation

```bash
# Clone repository
git clone https://github.com/iamkyr0/atlas.git
cd atlas

# Setup Go dependencies
cd chain
go mod download

# Setup Python SDK
cd ../sdk/python
pip install -r requirements.txt

# Start IPFS node (or use Infura)
ipfs daemon
```

### Docker Compose (Recommended)

```bash
# Start all services (chain, IPFS, node)
docker-compose up -d

# Check status
docker-compose ps
```

## 📖 Usage Guide

### 🎯 As a Client (Using Atlas for Fine-Tuning)

#### 1. Install Python SDK

```bash
cd sdk/python
pip install -r requirements.txt
```

#### 2. Initialize Client

```python
from atlas import AtlasClient

client = AtlasClient(
    ipfs_api_url="/ip4/127.0.0.1/tcp/5001",  # or "https://ipfs.infura.io:5001"
    chain_grpc_url="localhost:9090",  # Format: host:port
    creator="your_wallet_address"
)
```

#### 3. Upload Dataset & Submit Training Job

```python
import asyncio

async def train_model():
    async with AtlasClient() as client:
        # 1. Upload dataset
        dataset_cid = await client.upload_dataset("dataset.zip")
        print(f"Dataset CID: {dataset_cid}")
        
        # 2. Submit training job
        job_id = await client.submit_job(
            model_id="model-123",
            dataset_cid=dataset_cid,
            config={
                "epochs": 10,
                "batch_size": 32,
                "learning_rate": 0.001,
                "lora_rank": 8
            }
        )
        print(f"Job ID: {job_id}")
        
        # 3. Monitor progress (real-time via IPFS Pub/Sub)
        async for update in client.subscribe_to_job_updates(job_id):
            status = update.get('status')
            progress = update.get('progress', 0.0) * 100
            print(f"Status: {status}, Progress: {progress:.1f}%")
            
            if status == "completed":
                break
        
        # 4. Download trained model
        job = await client.get_job(job_id)
        await client.download_model(
            model_cid=job.get('model_cid'),
            output_path="./trained_model.pt"
        )
        print("Model downloaded!")

asyncio.run(train_model())
```

#### 4. Using Model for Inference

```python
from atlas.serving.predictor import Predictor
from atlas.serving.loader import ModelLoader

# Load model from IPFS (with caching)
loader = ModelLoader(ipfs_api_url="/ip4/127.0.0.1/tcp/5001")
model_path = loader.load(model_cid="QmXXX...")

# Create predictor
predictor = Predictor(model_path=model_path)

# Run inference
result = predictor.predict(
    input_data="Hello, how are you?",
    model_type="llm",  # or "vision", "speech", "embedding", "auto"
    options={
        "max_length": 100,
        "temperature": 0.7
    }
)

print(f"Result: {result}")
```

#### 5. Serve Model as HTTP Server

```bash
# Via CLI
atlas serve-model model-123 --port 8000

# Or via Python
from atlas.serving.server import ModelServer

server = ModelServer(chain_client=chain_client)
await server.serve(model_id="model-123", host="0.0.0.0", port=8000)
```

#### 6. CLI Commands

```bash
# Upload dataset
atlas upload-dataset dataset.zip --encrypt

# Submit training job
atlas submit-job model-123 QmXXX... --config config.json

# Monitor job
atlas get-job job-123
atlas list-jobs

# Download model
atlas download-model model-123 --output ./models/

# Serve model
atlas serve-model model-123 --port 8000

# Start daemon (REST API)
atlas daemon --port 8080
```

### 🖥️ As a Node (Providing Compute & Storage)

#### 1. Prerequisites

- Go 1.21+
- Docker
- IPFS node
- GPU (optional, for faster training)
- Minimum 8GB RAM
- Minimum 100GB storage

#### 2. Setup IPFS Node

```bash
# Install IPFS
wget https://dist.ipfs.io/go-ipfs/v0.20.0/go-ipfs_v0.20.0_linux-amd64.tar.gz
tar -xvzf go-ipfs_v0.20.0_linux-amd64.tar.gz
cd go-ipfs
sudo ./install.sh

# Initialize & start
ipfs init
ipfs daemon
```

#### 3. Build & Run Node

```bash
# Build node binary
cd node
go mod download
go build -o atlas-node cmd/node/main.go

# Run node
./atlas-node start \
  --chain-rpc localhost:9090 \
  --ipfs-api /ip4/127.0.0.1/tcp/5001 \
  --node-id node-001 \
  --address your_wallet_address
```

#### 4. Register Node to Blockchain

```bash
# Register via blockchain CLI
atlasd tx compute register-node \
  --node-id "node-001" \
  --address "your_wallet_address" \
  --from mykey \
  --chain-id atlas-chain
```

#### 5. Monitor Node Status

```bash
# Check node status
atlas-node status

# Check rewards
atlasd query reward node-rewards node-001

# Check reputation
atlasd query compute node-reputation node-001
```

#### 6. Docker Deployment

```bash
# Build image
docker build -t atlas-node:latest -f docker/Dockerfile.node .

# Run container
docker run -d \
  --name atlas-node \
  --gpus all \
  -v /tmp/atlas-work:/work \
  -e CHAIN_RPC_URL=localhost:9090 \
  -e IPFS_API_URL=/ip4/127.0.0.1/tcp/5001 \
  -e NODE_ID=node-001 \
  atlas-node:latest
```

## 🔧 Advanced Features

### Federated Learning

```python
job = await client.submit_job(
    model_id="model-123",
    dataset_cid=dataset_cid,
    config={
        "training_type": "federated",
        "min_clients": 5,
        "rounds": 10,
        "lora_enabled": True
    }
)
```

### Custom Training Script

```python
# Upload custom training script
script_cid = await client.upload_dataset("custom_train.py")

# Submit job with custom script
job = await client.submit_job(
    model_id="model-123",
    dataset_cid=dataset_cid,
    config={
        "training_script_cid": script_cid,
        "epochs": 10
    }
)
```

### Private Datasets

```python
# Upload with encryption
dataset_cid = await client.upload_dataset(
    "private_dataset.zip",
    encrypt=True,
    private_network=True
)
```

### Inference Network

```python
# LLM Inference
result = predictor.predict(
    input_data="Explain quantum computing",
    model_type="llm",
    options={"max_length": 200, "temperature": 0.8}
)

# Vision Model
result = predictor.predict(
    input_data="./image.jpg",
    model_type="vision"
)

# Speech-to-Text
result = predictor.predict(
    input_data="./audio.wav",
    model_type="speech"
)

# Embedding Service
result = predictor.predict(
    input_data="text to embed",
    model_type="embedding"
)
```

## 📊 Monitoring & Progress Tracking

### Real-time Progress Monitoring

```python
# Subscribe to job updates via IPFS Pub/Sub
async for update in client.subscribe_to_job_updates(job_id):
    print(f"Status: {update.get('status')}")
    print(f"Progress: {update.get('progress', 0.0) * 100:.1f}%")
    print(f"Tasks: {update.get('tasks', [])}")
    
    if update.get("status") == "completed":
        break
```

### Manual Query

```python
# Get job status
job = await client.get_job(job_id)
print(f"Status: {job.get('status')}")
print(f"Progress: {job.get('progress', 0.0) * 100:.1f}%")

# List all jobs
jobs = await client.list_jobs()
for job in jobs:
    print(f"{job.get('id')}: {job.get('status')} - {job.get('progress', 0.0) * 100:.1f}%")
```

## 🚢 Deployment

### Local Development

```bash
# Start IPFS
ipfs daemon

# Start chain (in separate terminal)
cd chain
atlasd start

# Start node (in separate terminal)
cd node
./atlas-node start
```

### Docker Compose

```bash
# Start all services
docker-compose up -d

# Check logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Kubernetes

```bash
# Deploy chain
kubectl apply -f k8s/chain-deployment.yaml

# Deploy IPFS
kubectl apply -f k8s/ipfs-deployment.yaml

# Deploy node
kubectl apply -f k8s/node-deployment.yaml
```

## 🔐 Security

- **Encryption**: Private datasets are encrypted before upload to IPFS
- **Private Networks**: Support for private IPFS networks
- **Proof of Computation**: Cryptographic proofs for task completion
- **Validation**: Content hash validation for shards
- **Reputation System**: Node reputation for trust management
- **Authentication**: API key and JWT support
- **Authorization**: RBAC (Role-Based Access Control)

## 🛠️ Configuration

### Environment Variables

```bash
# Client
export ATLAS_IPFS_API="/ip4/127.0.0.1/tcp/5001"
export ATLAS_CHAIN_GRPC="localhost:9090"
export ATLAS_CREATOR="your_wallet_address"

# Node
export CHAIN_RPC_URL="localhost:9090"
export IPFS_API_URL="/ip4/127.0.0.1/tcp/5001"
export NODE_ID="node-001"
export WORK_DIR="/tmp/atlas-work"
```

### Config File

```yaml
# config.yaml
node:
  id: "node-001"
  address: "your_wallet_address"
  chain_rpc_url: "localhost:9090"
  ipfs_api_url: "/ip4/127.0.0.1/tcp/5001"

resources:
  cpu_cores: 4
  gpu_enabled: true
  memory_gb: 16
  storage_gb: 500

heartbeat:
  interval_seconds: 30
  timeout_seconds: 90
```

## 📚 Documentation

### Main Documentation

- **[Client Usage Guide](client-usage.md)** - Complete guide for using Atlas as a client
- **[Node Usage Guide](node-usage.md)** - Guide for running compute nodes
- **[Architecture](architecture.md)** - System architecture overview
- **[Developer Guide](developer-guide.md)** - Guide for developers

### Component Documentation

- **[Chain](chain/README.md)** - Blockchain modules and implementation
- **[Node](node/README.md)** - Compute node implementation
- **[Federated Learning](federated-learning/README.md)** - FL implementation
- **[LoRA](lora/README.md)** - LoRA fine-tuning
- **[Storage](storage/README.md)** - Storage layer
- **[SDK](sdk/README.md)** - Client SDK documentation

## 🧪 Testing

```bash
# Run unit tests
cd chain
go test ./x/... -v

cd node
go test ./... -v

# Run integration tests
cd tests/integration
pytest -v
```

## 🐛 Troubleshooting

### Connection Issues

- **IPFS not connecting**: Check `ipfs id` and ensure IPFS daemon is running
- **Chain gRPC error**: Verify chain gRPC URL and port (default: 9090)
- **Network timeout**: Check firewall rules and network connectivity

### Job Failures

- **Job stuck**: Check node status and resource availability
- **Model download failed**: Verify model CID in IPFS network
- **Training error**: Check task logs in work directory

### Node Issues

- **Node not receiving tasks**: Verify node registration and status "online"
- **Resource exhaustion**: Check CPU/GPU/Memory utilization
- **Heartbeat timeout**: Verify network connectivity to blockchain

## 🤝 Contributing

Contributions are welcome! Please read our contributing guidelines and submit pull requests.

## 📄 License

MIT License - see LICENSE file for details

## 🔗 Links

- **Repository**: [GitHub](https://github.com/iamkyr0/atlas)
- **Documentation**: See [Documentation](#-documentation) section above
- **Issues**: [GitHub Issues](https://github.com/iamkyr0/atlas/issues)

## 🙏 Acknowledgments

- Cosmos SDK for blockchain framework
- IPFS for decentralized storage
- PyTorch and TensorFlow for ML frameworks

---

**Made with ❤️ by the Atlas Team**
