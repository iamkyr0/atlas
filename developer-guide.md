# Developer Guide

Guide for developers contributing to the Atlas platform.

## Project Structure

```
atlas/
├── chain/                    # Cosmos SDK blockchain
│   ├── x/                    # Custom modules
│   │   ├── compute/          # Compute node management
│   │   ├── training/         # Training job coordination
│   │   ├── storage/          # Storage node management
│   │   ├── reward/           # Reward distribution
│   │   ├── model/            # Model registry
│   │   ├── health/           # Health monitoring
│   │   ├── recovery/         # Recovery mechanisms
│   │   ├── sharding/         # Data sharding
│   │   └── validation/       # Validation logic
│   └── cmd/atlasd/           # Chain daemon
├── node/                     # Compute node
│   ├── executor/             # Task execution
│   ├── resource/             # Resource management
│   ├── health/               # Health monitoring
│   ├── recovery/             # Checkpoint recovery
│   ├── proof/                # Proof generation
│   ├── validator/            # Assignment validation
│   └── network/              # Network utilities
├── federated-learning/       # FL implementation
│   ├── client/               # Client-side training
│   ├── aggregator/           # Gradient aggregation
│   ├── protocols/            # Communication protocols
│   ├── sharding/             # Shard coordination
│   └── validation/           # Gradient validation
├── lora/                     # LoRA fine-tuning
│   ├── adapters/             # LoRA adapters
│   └── training/             # Training logic
├── storage/                  # Storage layer
│   ├── manager/              # IPFS manager
│   ├── sharding/             # Data sharding
│   ├── validation/          # Content validation
│   └── pubsub/               # Pub/sub messaging
├── sdk/                      # Client SDK
│   └── python/               # Python SDK
│       ├── atlas/            # SDK package
│       ├── tests/            # Unit tests
│       └── setup.py          # Package setup
└── docs/                     # Documentation
```

## Development Setup

### Prerequisites

- Go 1.21+
- Python 3.10+
- Docker
- IPFS node
- Cosmos SDK

### Build Chain

```bash
cd chain
go mod download
go build -o atlasd ./cmd/atlasd
```

### Build Node

```bash
cd node
go mod download
go build -o atlas-node ./cmd/node/main.go
```

### Setup Python SDK

```bash
cd sdk/python
pip install -r requirements.txt
pip install -e .
```

## Testing

### Chain Tests

```bash
cd chain
go test ./... -v
```

### Node Tests

```bash
cd node
go test ./... -v
```

### FL Tests

```bash
cd federated-learning
go test ./... -v
```

### LoRA Tests

```bash
cd lora
go test ./... -v
```

### Storage Tests

```bash
cd storage
go test ./... -v
```

### Python SDK Tests

```bash
cd sdk/python
pytest tests/ -v
```

## Code Style

### Go

- Use `gofmt` for formatting
- Follow Go naming conventions
- Write unit tests for all functions
- Remove all comments (including godoc)

### Python

- Follow PEP 8
- Use type hints
- Write docstrings (but remove in final code)
- Write unit tests

## Adding New Features

### Chain Module

1. Define message types in `types/`
2. Implement keeper logic in `keeper/`
3. Add message handlers in `keeper/msg_server.go`
4. Add query handlers in `keeper/grpc_query.go`
5. Write unit tests

### Node Component

1. Create package in `node/`
2. Implement functionality
3. Add to main.go if needed
4. Write unit tests

### SDK Feature

1. Add to appropriate module in `sdk/python/atlas/`
2. Update CLI if needed
3. Write unit tests
4. Update documentation

## Debugging

### Chain Debugging

```bash
atlasd start --log-level debug
```

### Node Debugging

```bash
atlas-node start --log-level debug
```

### IPFS Debugging

```bash
ipfs daemon --debug
```

## Contributing

1. Fork the repository
2. Create feature branch
3. Make changes
4. Write tests
5. Run all tests
6. Submit pull request

## Documentation

- Update README.md files in each folder
- Update main README.md if needed
- Add examples for new features
- Update architecture.md for architectural changes

## Release Process

1. Update version numbers
2. Run all tests
3. Update documentation
4. Create release tag
5. Build binaries
6. Publish release

