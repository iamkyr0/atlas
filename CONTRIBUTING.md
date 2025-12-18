# Contributing to Atlas

First off, thank you for considering contributing to Atlas! 🎉

This document provides guidelines and instructions for contributing to the Atlas decentralized AI platform. Following these guidelines helps communicate that you respect the time of the developers managing and developing this open source project.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct:
- Be respectful and inclusive
- Welcome newcomers and help them learn
- Focus on constructive feedback
- Respect different viewpoints and experiences

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce** the behavior
- **Expected behavior**
- **Actual behavior**
- **Environment details** (OS, Go/Python version, etc.)
- **Screenshots/logs** if applicable

**Bug Report Template:**
```markdown
## Bug Description
Brief description of the bug

## Steps to Reproduce
1. Step one
2. Step two
3. See error

## Expected Behavior
What should happen

## Actual Behavior
What actually happens

## Environment
- OS: [e.g., macOS 14.0]
- Go version: [e.g., 1.21]
- Python version: [e.g., 3.10]
- Atlas version: [e.g., 0.1.0]

## Additional Context
Any other relevant information
```

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, include:

- **Clear description** of the enhancement
- **Use case** - why is this useful?
- **Proposed solution** (if you have one)
- **Alternatives considered**

### Pull Requests

Pull requests are welcome! Please follow these steps:

1. **Fork the repository**
2. **Create a feature branch** (`git checkout -b feature/amazing-feature`)
3. **Make your changes**
4. **Add tests** for new functionality
5. **Ensure all tests pass**
6. **Update documentation** if needed
7. **Commit your changes** (`git commit -m 'Add amazing feature'`)
8. **Push to the branch** (`git push origin feature/amazing-feature`)
9. **Open a Pull Request**

## Development Setup

### Prerequisites

- **Go** 1.21 or higher
- **Python** 3.10 or higher
- **Docker** (optional, for containerized development)
- **IPFS node** (or use Infura/IPFS gateway)
- **Git**

### Getting Started

1. **Fork and clone the repository:**
```bash
git clone https://github.com/iamkyr0/atlas.git
cd atlas
```

2. **Set up Go dependencies:**
```bash
cd chain
go mod download
```

3. **Set up Python SDK:**
```bash
cd ../sdk/python
pip install -r requirements.txt
pip install -e .  # Install in development mode
```

4. **Set up IPFS (optional):**
```bash
# Install IPFS
# Or use Infura gateway: https://ipfs.infura.io:5001
ipfs daemon
```

5. **Run tests:**
```bash
# Go tests
cd chain
go test ./x/... -v

# Python tests
cd ../sdk/python
pytest -v
```

## Coding Standards

### Go Code Style

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Use `golint` or `golangci-lint` for linting
- Write meaningful comments for exported functions
- Keep functions small and focused

**Example:**
```go
// RegisterNode registers a new compute node on the blockchain
func (k Keeper) RegisterNode(ctx sdk.Context, node types.Node) error {
    // Implementation
}
```

### Python Code Style

- Follow [PEP 8](https://pep8.org/) style guide
- Use `black` for formatting
- Use `pylint` or `flake8` for linting
- Write docstrings for all public functions
- Type hints are encouraged

**Example:**
```python
def submit_job(
    self,
    model_id: str,
    dataset_cid: str,
    config: Dict[str, Any],
) -> str:
    """
    Submit a training job to the Atlas platform.
    
    Args:
        model_id: Model identifier
        dataset_cid: IPFS CID of the dataset
        config: Training configuration
        
    Returns:
        Job ID
    """
    # Implementation
```

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/) format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(chain): add node reputation system
fix(sdk): handle IPFS connection timeout
docs(readme): update installation instructions
test(node): add resource manager tests
```

## Project Structure

```
atlas/
├── chain/              # Cosmos SDK blockchain
│   ├── x/             # Custom modules
│   └── cmd/           # CLI commands
├── node/              # Compute node
│   ├── executor/      # Task execution
│   ├── resource/      # Resource management
│   └── health/       # Health monitoring
├── federated-learning/ # FL implementation
├── lora/             # LoRA fine-tuning
├── storage/          # IPFS storage layer
├── sdk/              # Client SDK
│   └── python/       # Python SDK
└── tests/            # Integration tests
```

## Current Priorities

Based on our codebase analysis, here are the current high-priority tasks that need immediate attention:

### 🔴 CRITICAL Priority (Blocking Core Functionality)

#### 1. Protobuf Generation & gRPC Client
**Status**: Not Started | **Impact**: HIGH | **Effort**: Medium (1-2 days)

**Problem**: 
- Python SDK (`sdk/python/atlas/chain/client.py`) requires protobuf stubs that haven't been generated
- All methods throw `NotImplementedError` when stubs are missing
- Node blockchain client (`node/validator/blockchain_client.go`) is still a placeholder

**Files to Fix**:
- `sdk/python/atlas/chain/client.py` - All methods need protobuf stubs
- `node/validator/blockchain_client.go` - All methods are placeholders
- `node/cmd/node/main.go` - `registerNodeOnBlockchain()` function

**How to Help**:
```bash
# Generate protobuf stubs for Python SDK
cd chain
# Generate .proto files from Cosmos SDK modules
# Generate Python stubs
python -m grpc_tools.protoc -I. --python_out=../sdk/python/atlas/chain --grpc_python_out=../sdk/python/atlas/chain *.proto
```

#### 2. Node Blockchain Client Implementation
**Status**: Not Started | **Impact**: HIGH | **Effort**: High (2-3 days)

**Problem**:
- `HTTPBlockchainClient` in `node/validator/blockchain_client.go` is a placeholder
- Nodes cannot query blockchain for assignment validation
- Node registration still returns errors

**Files to Fix**:
- `node/validator/blockchain_client.go` - Implement all methods
- `node/cmd/node/main.go` - Fix `registerNodeOnBlockchain()`

**How to Help**:
- Implement gRPC client in Go for node
- Or implement HTTP REST client using Cosmos SDK REST API
- Implement `registerNodeOnBlockchain()` with gRPC/HTTP client

### 🟠 HIGH Priority (Important Features)

#### 3. Model Registration via gRPC
**Status**: Not Started | **Impact**: MEDIUM | **Effort**: Low (1 day)

**Problem**: 
- `register_model()` in `sdk/python/atlas/chain/client.py` line 222 still throws `NotImplementedError`

**Files to Fix**:
- `sdk/python/atlas/chain/client.py` - Implement `register_model()` method
- Add to protobuf message types
- Update keeper to handle model registration

#### 4. Model Sharding by Layers
**Status**: Not Started | **Impact**: MEDIUM | **Effort**: Medium (1-2 days)

**Problem**:
- `SplitModelByLayers()` in `storage/sharding/` is still a placeholder
- Strategy "layer" is not implemented

**Files to Fix**:
- `storage/sharding/splitter.go` - Implement `SplitModelByLayers()`
- `storage/sharding/manager.go` - Update shard manager

**How to Help**:
- Implement parsing of model layers (PyTorch/TensorFlow)
- Split model based on layer boundaries
- Update shard manager to support layer-based sharding

#### 5. Integration Tests
**Status**: Partial | **Impact**: HIGH | **Effort**: High (2-3 days)

**Problem**:
- `tests/integration/test_job_submission.py` - Many `NotImplementedError`
- `tests/integration/test_node_failure.py` - Many `NotImplementedError`
- End-to-end flow is not tested

**Files to Fix**:
- `tests/integration/test_job_submission.py` - Fix all tests
- `tests/integration/test_node_failure.py` - Fix all tests
- Add end-to-end test with local chain + IPFS

**How to Help**:
- Fix integration tests with mocks/stubs
- Implement end-to-end test with local chain + IPFS
- Test full job lifecycle

#### 6. Error Handling & Recovery
**Status**: Partial | **Impact**: MEDIUM | **Effort**: Medium (2 days)

**Problem**:
- Many `NotImplementedError` without fallback mechanisms
- Network failure handling is not robust
- Task failure recovery needs improvement

**Files to Fix**:
- All files with `NotImplementedError` - Add fallback mechanisms
- Improve retry logic with exponential backoff
- Better error messages and logging

#### 7. LoRA Training Script
**Status**: Partial | **Impact**: MEDIUM | **Effort**: Medium (1-2 days)

**Problem**:
- `lora/training/trainer.go` line 132-175 - Training script is still simulated
- No actual LoRA training logic
- Comment says: "In production, this would..."

**Files to Fix**:
- `lora/training/trainer.go` - Implement actual LoRA training

**How to Help**:
- Implement actual LoRA training with PyTorch
- Load base model and apply LoRA adapters
- Train only LoRA parameters
- Save updated weights

### 🟡 MEDIUM Priority (Important but Can Wait)

#### 8. Security Audit
**Status**: Not Started | **Impact**: MEDIUM | **Effort**: High (3-5 days)

- Input validation needs to be stricter
- API key management needs review
- Encryption key management needs improvement

#### 9. Test Coverage
**Status**: Partial | **Impact**: LOW | **Effort**: High (3-5 days)

- Some modules don't have unit tests
- Test coverage < 80% for some components
- Target: > 80% coverage for all modules

#### 10. Performance Optimization
**Status**: Partial | **Impact**: LOW | **Effort**: Medium (2-3 days)

- Resource allocation could be more efficient
- Model caching strategy needs improvement
- Network optimization for IPFS

#### 11. Documentation
**Status**: Partial | **Impact**: LOW | **Effort**: Low (1-2 days)

- API documentation is incomplete
- Deployment guides need more detail
- Troubleshooting guides need to be added

### 🟢 Good First Issues

Perfect for newcomers looking to get started:

- Small bug fixes
- Documentation improvements
- Test additions
- Code cleanup
- Minor feature additions

Look for issues labeled `good first issue` on GitHub.

## Getting Started with Priority Tasks

### Recommended Order of Implementation

We recommend tackling priority tasks in this order for maximum impact:

1. **Protobuf Generation** (1-2 days) - Unblocks many other features
   - Generate protobuf stubs for Python SDK
   - Fix all `NotImplementedError` in `chain/client.py`

2. **Node Blockchain Client** (2-3 days) - Critical for node functionality
   - Implement gRPC client in Go
   - Fix `HTTPBlockchainClient`
   - Fix `registerNodeOnBlockchain()`

3. **Integration Tests** (2-3 days) - Ensures system reliability
   - Fix all integration tests
   - Setup test environment with local chain

4. **LoRA Training Script** (1-2 days) - Core functionality
   - Implement actual LoRA training
   - Test with real model

5. **Model Registration & Sharding** (2-3 days) - Important features
   - Implement model registration
   - Implement layer-based sharding

**Total Estimated Time**: 8-13 days for critical + high priority items

### How to Claim a Task

1. Check if the task is already assigned in GitHub Issues
2. Comment on the issue to claim it
3. Create a branch: `git checkout -b fix/protobuf-generation`
4. Work on the task following our coding standards
5. Submit a PR with tests and documentation

### Questions?

If you're unsure about how to approach a priority task:

- Open a GitHub Discussion to ask questions
- Check existing issues for similar work
- Review the codebase to understand the context
- Reach out to maintainers for guidance

## Testing Guidelines

### Go Tests

- Write unit tests for all new functions
- Use table-driven tests when appropriate
- Mock external dependencies
- Aim for >80% code coverage

**Example:**
```go
func TestRegisterNode(t *testing.T) {
    tests := []struct {
        name    string
        node    types.Node
        wantErr bool
    }{
        {
            name: "valid node",
            node: types.Node{
                ID:      "node-1",
                Address: "cosmos1abc...",
            },
            wantErr: false,
        },
        // More test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Python Tests

- Use `pytest` for testing
- Write unit tests and integration tests
- Use fixtures for common setup
- Mock external API calls

**Example:**
```python
import pytest
from atlas import AtlasClient

@pytest.mark.asyncio
async def test_submit_job():
    async with AtlasClient() as client:
        job_id = await client.submit_job(
            model_id="model-123",
            dataset_cid="QmXXX...",
            config={"epochs": 10}
        )
        assert job_id is not None
```

## Documentation

### Code Documentation

- Write clear docstrings/comments
- Explain "why" not just "what"
- Update README.md for user-facing changes
- Add examples for complex functions

### User Documentation

- Update relevant `.md` files
- Add code examples
- Include screenshots for UI changes
- Keep documentation up to date

## Pull Request Process

1. **Update CHANGELOG.md** with your changes
2. **Update documentation** if needed
3. **Add tests** for new features
4. **Ensure all tests pass**
5. **Request review** from maintainers

### PR Checklist

- [ ] Code follows style guidelines
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] All tests pass
- [ ] No merge conflicts

### Review Process

- Maintainers will review within 48 hours
- Address review comments promptly
- Be open to feedback and suggestions
- Keep PRs focused and small when possible

## Recognition

Contributors will be:
- **Listed in CONTRIBUTORS.md**
- **Mentioned in release notes**
- **Given credit in documentation**
- **Invited to join core team** (for significant contributions)

## Getting Help

- **GitHub Discussions**: For questions and discussions
- **GitHub Issues**: For bug reports and feature requests
- **Repository**: [https://github.com/iamkyr0/atlas](https://github.com/iamkyr0/atlas)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Thank You!

Thank you for taking the time to contribute to Atlas! Every contribution, no matter how small, makes a difference. 🚀

---

**Happy coding!** 💻

