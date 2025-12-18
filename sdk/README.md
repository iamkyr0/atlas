# Atlas SDK

Client SDK for the Atlas Decentralized AI Platform. Provides Python SDK for interacting with the platform in a fully decentralized P2P manner.

## Components

### Python SDK (`python/`)
Python SDK with CLI and daemon mode support.

**Features:**
- P2P communication via IPFS
- Direct blockchain interaction via gRPC
- Job submission and monitoring
- Model registration and serving
- Dataset upload with encryption support
- Real-time updates via IPFS pub/sub
- CLI interface
- Daemon mode for application integration

**See:** [python/README.md](python/README.md) for detailed usage.

## Architecture

The SDK uses a fully decentralized architecture:

- **No API Gateway**: All communication is peer-to-peer
- **IPFS Pub/Sub**: Real-time messaging and updates
- **Blockchain gRPC**: Direct chain state queries and transactions
- **IPFS DHT**: Node discovery with blockchain fallback
- **Private Networks**: Support for encrypted private IPFS networks

## Quick Start

```bash
cd python
pip install -r requirements.txt

# Use CLI
atlas submit-job model-123 QmXXX... --config config.json

# Or use Python API
python -c "from atlas import AtlasClient; ..."
```

## Documentation

- [Python SDK Documentation](python/README.md)
- [Client Usage Guide](../client-usage.md)
- [Architecture Overview](../architecture.md)

