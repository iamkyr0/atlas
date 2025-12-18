# Atlas Test Suite

## Overview
Comprehensive test suite for Atlas decentralized AI platform.

## Test Structure

```
tests/
├── unit/              # Unit tests for individual components
├── integration/       # Integration tests
├── e2e/              # End-to-end tests
└── performance/      # Performance benchmarks
```

## Running Tests

### Unit Tests
```bash
# Chain tests
cd chain && go test ./...

# Node tests
cd node && go test ./...

# Python SDK tests
cd sdk/python && pytest tests/
```

### Integration Tests
```bash
make test-integration
```

### E2E Tests
```bash
make test-e2e
```

## Test Coverage Goals
- Unit tests: >80% coverage
- Integration tests: Cover all critical paths
- E2E tests: Cover main user workflows

