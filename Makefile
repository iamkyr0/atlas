.PHONY: help setup build test clean docker-build

help:
	@echo "Atlas - Decentralized AI Platform"
	@echo ""
	@echo "Available targets:"
	@echo "  setup       - Setup development environment"
	@echo "  build       - Build all components"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  docker-build - Build Docker images"

setup:
	@echo "Setting up development environment..."
	cd chain && go mod download
	cd sdk/python && pip install -r requirements.txt
	@echo "Setup complete!"

build:
	@echo "Building Atlas..."
	cd chain && make build
	cd node && go build -o bin/atlas-node ./cmd/node
	cd api && go build -o bin/atlas-api ./cmd/api
	@echo "Build complete!"

test:
	@echo "Running tests..."
	cd chain && go test ./...
	cd node && go test ./...
	cd api && go test ./...
	cd sdk/python && pytest tests/
	@echo "Tests complete!"

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf dist/
	find . -type d -name "__pycache__" -exec rm -r {} +
	find . -type f -name "*.pyc" -delete
	@echo "Clean complete!"

docker-build:
	@echo "Building Docker images..."
	docker-compose build
	@echo "Docker build complete!"

