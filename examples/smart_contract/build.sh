#!/bin/bash

# Build script to compile the Go smart contract to WebAssembly for Massa
# Usage: ./build.sh

set -e

echo "Building Massa smart contract from Go to WebAssembly..."

# Ensure we're in the correct directory
cd "$(dirname "$0")"

# Set environment variables for WebAssembly compilation
export GOOS=js
export GOARCH=wasm

# Build flags for WebAssembly optimization
BUILD_FLAGS="-ldflags=-s -ldflags=-w"

# Output file
OUTPUT_FILE="smart_contract.wasm"

# Build the smart contract
echo "Compiling Go code to WebAssembly..."
go build -ldflags="-s -w" -o $OUTPUT_FILE main.go

if [ $? -eq 0 ]; then
    echo "✅ Smart contract compiled successfully to $OUTPUT_FILE"
    echo "📊 File size: $(wc -c < $OUTPUT_FILE) bytes"
    
    # Check if wasm-opt is available for further optimization
    if command -v wasm-opt &> /dev/null; then
        echo "🔧 Optimizing WebAssembly with wasm-opt..."
        wasm-opt -Oz --enable-bulk-memory $OUTPUT_FILE -o ${OUTPUT_FILE}.optimized
        mv ${OUTPUT_FILE}.optimized $OUTPUT_FILE
        echo "✅ WebAssembly optimized"
        echo "📊 Optimized file size: $(wc -c < $OUTPUT_FILE) bytes"
    else
        echo "ℹ️  wasm-opt not found. Install it with 'npm install -g wasm-opt' for smaller WASM files"
    fi
    
    # Validate the WASM file
    echo "🔍 Validating WebAssembly file..."
    if command -v wasm-validate &> /dev/null; then
        if wasm-validate $OUTPUT_FILE; then
            echo "✅ WebAssembly file is valid"
        else
            echo "❌ WebAssembly file validation failed"
            exit 1
        fi
    else
        echo "ℹ️  wasm-validate not found. Install wabt tools for validation"
    fi
    
    echo ""
    echo "🚀 Smart contract is ready for deployment!"
    echo "📝 To deploy this contract to Massa:"
    echo "   1. Use the Massa client or wallet"
    echo "   2. Deploy the contract using: $OUTPUT_FILE"
    echo "   3. The main function will be called automatically"
    echo ""
    echo "🔧 Available functions in this contract:"
    echo "   - main() - Entry point, called automatically"
    echo "   - GetContractInfo() - Returns contract information"
    echo "   - UpdateGreeting(string) - Updates greeting message (deployer only)"
    echo "   - GetGreeting() - Returns current greeting"
    echo "   - TransferFunds(string, uint64) - Transfers funds to address"
    
else
    echo "❌ Build failed"
    exit 1
fi