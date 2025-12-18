# Atlas Python SDK

Decentralized AI Platform SDK for Python - P2P Architecture

## Installation

```bash
pip install -r requirements.txt
# Or install as package
pip install -e .
```

## Usage

### CLI

```bash
# Submit a job
atlas submit-job <model_id> <dataset_cid> --config config.json

# List jobs
atlas list-jobs

# Get job status
atlas get-job <job_id>

# Upload dataset
atlas upload-dataset <path> --encrypt

# Register model
atlas register-model <path> --name <name> --version <version>

# Download model
atlas download-model <model_id>

# Serve model for inference
atlas serve-model <model_id> --port 8000

# Start daemon mode
atlas daemon --port 8080
```

### Python API

```python
from atlas import AtlasClient

# Initialize client (P2P mode, no API Gateway)
client = AtlasClient(
    ipfs_api_url="/ip4/127.0.0.1/tcp/5001",
    chain_grpc_url="localhost:9090",
    creator="your_address"
)

# Submit job
job_id = await client.submit_job(
    model_id="model-123",
    dataset_cid="QmXXX...",
    config={"epochs": 10}
)

# Monitor job
async for update in client.subscribe_to_job_updates(job_id):
    print(update)
```

## Environment Variables

- `ATLAS_IPFS_API`: IPFS API URL (default: `/ip4/127.0.0.1/tcp/5001`)
- `ATLAS_CHAIN_GRPC`: Chain gRPC URL (default: `localhost:9090`)
- `ATLAS_CREATOR`: Creator address for transactions

## Architecture

The SDK uses a fully decentralized P2P architecture:

- **Direct Blockchain Access**: gRPC client for blockchain interaction
- **IPFS Pub/Sub**: Real-time communication
- **Node Discovery**: IPFS DHT with blockchain fallback
- **No API Gateway**: All communication is peer-to-peer

