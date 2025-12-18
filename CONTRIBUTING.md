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

## Areas We Need Help

### High Priority

- **Frontend/UI**: Web interface for the platform
- **Testing**: More comprehensive test coverage
- **Documentation**: Tutorials, guides, and examples
- **Performance**: Optimizing resource usage
- **Error Handling**: Better error messages and recovery

### Medium Priority

- **Security**: Security audits and improvements
- **DevOps**: CI/CD pipelines and automation
- **Monitoring**: Better observability and metrics
- **Internationalization**: Multi-language support
- **Accessibility**: Improve accessibility features

### Good First Issues

Look for issues labeled `good first issue` - these are perfect for newcomers:

- Small bug fixes
- Documentation improvements
- Test additions
- Code cleanup
- Minor feature additions

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
- **Discord**: [Coming soon]
- **Email**: [Your email]

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Thank You!

Thank you for taking the time to contribute to Atlas! Every contribution, no matter how small, makes a difference. 🚀

---

**Happy coding!** 💻

