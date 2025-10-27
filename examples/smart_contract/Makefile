.PHONY: build run clean test wasm

# Default target
all: build

# Run the smart contract locally (with stub functions)
run:
	go run main.go

# Build for WebAssembly (Massa blockchain)
wasm:
	GOOS=js GOARCH=wasm go build -o smart_contract.wasm main.go
	@echo "WebAssembly build complete: smart_contract.wasm"
	@echo "Size: $$(wc -c < smart_contract.wasm) bytes"

# Alias for wasm target
build: wasm

# Test the smart contract
test:
	go test ./...

# Clean build artifacts
clean:
	rm -f smart_contract.wasm

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Show help
help:
	@echo "Available targets:"
	@echo "  run     - Run the smart contract locally"
	@echo "  build   - Build the smart contract to WebAssembly"
	@echo "  wasm    - Build the smart contract to WebAssembly"
	@echo "  test    - Run tests"
	@echo "  clean   - Clean build artifacts"
	@echo "  deps    - Install dependencies"
	@echo "  fmt     - Format code"
	@echo "  help    - Show this help message"