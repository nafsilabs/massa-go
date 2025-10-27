/*
Package sc provides a comprehensive Go SDK for developing Massa smart contracts
that compile to WebAssembly.

This package offers a complete set of functions for interacting with the Massa blockchain,
including storage operations, contract calls, event generation, cryptographic functions,
and more. It's designed to be similar to the official massa-as-sdk but written in pure Go.

# Key Features

• Pure Go implementation with full type safety
• WebAssembly compilation support
• Automatic memory management for WASM
• Mock implementations for testing
• Comprehensive API coverage
• Rich error handling and logging

# Core Modules

The SDK is organized into several key modules:

Address Operations:
  - Address validation and creation
  - Public key to address conversion
  - EOA (Externally Owned Account) detection

Storage Management:
  - Key-value storage operations
  - JSON serialization support
  - Cross-address storage access

Execution Context:
  - Caller and callee information
  - Timestamp and gas information
  - Call stack access
  - Chain and thread information

Contract Interactions:
  - Function calls to other contracts
  - Local execution capabilities
  - Bytecode operations
  - Message passing
  - Deferred calls

Coin Operations:
  - Balance queries
  - Coin transfers
  - Amount formatting and parsing
  - Batch operations

Events and Logging:
  - Event generation
  - Structured logging
  - Custom event types
  - Debug utilities

Cryptographic Functions:
  - Multiple hash algorithms (BLAKE3, SHA-256, Keccak-256, MiMC)
  - Signature verification
  - Ethereum compatibility
  - Merkle tree utilities

Operation Datastore:
  - Operation-level data access
  - Key enumeration and filtering
  - Categorized data access

# Quick Start

Basic smart contract example:

	package main

	import massa "github.com/nafsilabs/massa-go/sc"

	func main() {
		// Generate an event
		massa.GenerateEvent("Hello, Massa!")

		// Store some data
		massa.SetString("greeting", "Hello from Go!")

		// Get caller information
		caller := massa.Caller()
		massa.Printf("Called by: %s", caller.String())

		// Check balance
		balance := massa.Balance()
		massa.Printf("Balance: %s", massa.FormatBalance(balance))
	}

# Building for WebAssembly

To compile your smart contract for the Massa blockchain:

	export GOOS=js
	export GOARCH=wasm
	go build -ldflags="-s -w" -o contract.wasm main.go

# Testing

The SDK automatically provides mock implementations when not running in a
WebAssembly environment, making it easy to test your contracts:

	func TestContract(t *testing.T) {
		// These calls will use mock implementations
		caller := massa.Caller()
		massa.SetString("test", "value")
		value := massa.GetString("test")

		// Test your contract logic...
	}

# Architecture

The SDK is designed to work seamlessly in both WebAssembly and native Go environments:

• In WebAssembly: Uses native Massa blockchain functions via imports
• In Native Go: Uses mock implementations for testing and development

The isWasmRuntime() function automatically detects the environment and
routes function calls appropriately.

# Memory Management

WebAssembly memory management is handled automatically through helper functions:

• newWasmString() - Converts Go strings to WASM pointers
• newWasmBytes() - Converts Go byte slices to WASM pointers
• wasmStringFromPtr() - Converts WASM pointers to Go strings
• wasmBytesFromPtr() - Converts WASM pointers to Go byte slices

# Error Handling

The SDK provides comprehensive error handling with multiple approaches:

• Return errors from functions where appropriate
• Panic for programming errors (e.g., invalid addresses)
• Log errors with structured data
• Emit error events for blockchain visibility

# Performance Considerations

• Minimize storage operations and data size
• Use efficient algorithms and data structures
• Leverage WASM optimization tools (wasm-opt)
• Be mindful of gas costs for operations

# Security Best Practices

• Always validate input parameters
• Implement proper access controls
• Handle errors gracefully
• Use atomic operations for state changes
• Be aware of reentrancy issues

For more information and examples, see the project repository at:
https://github.com/nafsilabs/massa-go
*/
package sc
