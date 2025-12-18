.PHONY: build install test clean

build:
	@echo "Building atlasd..."
	go build -o bin/atlasd ./cmd/atlasd

install:
	@echo "Installing atlasd..."
	go install ./cmd/atlasd

test:
	@echo "Running tests..."
	go test ./...

clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean

