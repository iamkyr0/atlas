# Atlas - Decentralized AI Platform

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![Python Version](https://img.shields.io/badge/Python-3.10+-blue.svg)](https://www.python.org/)

**Atlas** adalah platform Infrastructure-as-a-Service (IaaS) terdesentralisasi untuk fine-tuning dan serving model AI. Platform ini menggunakan blockchain untuk koordinasi dan reward distribution, dengan node-node yang berkontribusi compute dan storage resources seperti sistem mining cryptocurrency.

## 🌟 Features

- **Decentralized Training**: Fine-tune AI models menggunakan distributed compute nodes
- **Federated Learning**: Privacy-preserving distributed training dengan gradient aggregation
- **LoRA Fine-Tuning**: Efficient fine-tuning dengan Low-Rank Adaptation
- **Model Serving**: Serve models untuk inference via HTTP/gRPC API
- **Inference Network**: Distributed network untuk LLM, Vision, Speech-to-text, dan Embedding services
- **Blockchain Coordination**: Cosmos SDK blockchain untuk job coordination dan reward distribution
- **IPFS Storage**: Decentralized storage untuk datasets, models, dan checkpoints
- **P2P Communication**: Direct peer-to-peer communication tanpa API Gateway
- **Auto Resource Detection**: Automatic detection CPU, GPU, RAM, storage, dan network speed
- **Fault Tolerance**: Automatic task reassignment, checkpoint recovery, dan graceful degradation

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

- **Blockchain Layer** (`chain/`): Cosmos SDK dengan 10 custom modules
- **Compute Node** (`node/`): Distributed compute nodes untuk training & inference
- **Federated Learning** (`federated-learning/`): Privacy-preserving distributed training
- **LoRA** (`lora/`): Efficient fine-tuning dengan Low-Rank Adaptation
- **Storage Layer** (`storage/`): IPFS-based decentralized storage
- **Client SDK** (`sdk/`): Python SDK dengan CLI dan daemon mode

## 🚀 Quick Start

### Prerequisites

- **Go** 1.21+
- **Python** 3.10+
- **Docker** (optional, untuk containerized deployment)
- **IPFS node** (atau gunakan Infura/IPFS gateway)
- **Cosmos SDK** (untuk blockchain)

### Installation

```bash
# Clone repository
git clone <repository-url>
cd atlas

# Setup Go dependencies
cd chain
go mod download

# Setup Python SDK
cd ../sdk/python
pip install -r requirements.txt

# Start IPFS node (atau gunakan Infura)
ipfs daemon
```

### Docker Compose (Recommended)

```bash
# Start semua services (chain, IPFS, node)
docker-compose up -d

# Check status
docker-compose ps
```

## 📖 Usage Guide

### 🎯 Sebagai Client (Menggunakan Atlas untuk Fine-Tuning)

#### 1. Install Python SDK

```bash
cd sdk/python
pip install -r requirements.txt
```

#### 2. Initialize Client

```python
from atlas import AtlasClient

client = AtlasClient(
    ipfs_api_url="/ip4/127.0.0.1/tcp/5001",  # atau "https://ipfs.infura.io:5001"
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

#### 4. Menggunakan Model untuk Inference

```python
from atlas.serving.predictor import Predictor
from atlas.serving.loader import ModelLoader

# Load model dari IPFS (dengan caching)
loader = ModelLoader(ipfs_api_url="/ip4/127.0.0.1/tcp/5001")
model_path = loader.load(model_cid="QmXXX...")

# Create predictor
predictor = Predictor(model_path=model_path)

# Run inference
result = predictor.predict(
    input_data="Hello, how are you?",
    model_type="llm",  # atau "vision", "speech", "embedding", "auto"
    options={
        "max_length": 100,
        "temperature": 0.7
    }
)

print(f"Result: {result}")
```

#### 5. Serve Model sebagai HTTP Server

```bash
# Via CLI
atlas serve-model model-123 --port 8000

# Atau via Python
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

### 🖥️ Sebagai Node (Menyediakan Compute & Storage)

#### 1. Prerequisites

- Go 1.21+
- Docker
- IPFS node
- GPU (optional, untuk training lebih cepat)
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

#### 4. Register Node ke Blockchain

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

# Submit job dengan custom script
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
# Upload dengan encryption
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

# Start chain (dalam terminal terpisah)
cd chain
atlasd start

# Start node (dalam terminal terpisah)
cd node
./atlas-node start
```

### Docker Compose

```bash
# Start semua services
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

- **Encryption**: Private datasets di-encrypt sebelum upload ke IPFS
- **Private Networks**: Support untuk private IPFS networks
- **Proof of Computation**: Cryptographic proofs untuk task completion
- **Validation**: Content hash validation untuk shards
- **Reputation System**: Node reputation untuk trust management
- **Authentication**: API key dan JWT support
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

- **[Client Usage Guide](client-usage.md)** - Panduan lengkap untuk menggunakan Atlas sebagai client
- **[Node Usage Guide](node-usage.md)** - Panduan untuk menjalankan compute nodes
- **[Architecture](architecture.md)** - Overview arsitektur sistem
- **[Developer Guide](developer-guide.md)** - Panduan untuk developers

### Component Documentation

- **[Chain](chain/README.md)** - Blockchain modules dan implementation
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

- **IPFS tidak connect**: Check `ipfs id` dan pastikan IPFS daemon running
- **Chain gRPC error**: Verify chain gRPC URL dan port (default: 9090)
- **Network timeout**: Check firewall rules dan network connectivity

### Job Failures

- **Job stuck**: Check node status dan resource availability
- **Model download failed**: Verify model CID di IPFS network
- **Training error**: Check task logs di work directory

### Node Issues

- **Node tidak receive tasks**: Verify node registration dan status "online"
- **Resource exhaustion**: Check CPU/GPU/Memory utilization
- **Heartbeat timeout**: Verify network connectivity ke blockchain

## 🤝 Contributing

Contributions are welcome! Please read our contributing guidelines and submit pull requests.

## 📄 License

MIT License - see LICENSE file for details

## 🔗 Links

- **Repository**: [GitHub](https://github.com/atlas/atlas)
- **Documentation**: [Docs](https://docs.atlas.ai)
- **Discord**: [Community](https://discord.gg/atlas)

## 🙏 Acknowledgments

- Cosmos SDK untuk blockchain framework
- IPFS untuk decentralized storage
- PyTorch dan TensorFlow untuk ML frameworks

---

**Made with ❤️ by the Atlas Team**
