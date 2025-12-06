![massa-go](https://raw.githubusercontent.com/jwmdev/massa-go/master/media/massa-go.svg)

# Massa Go SDK
A comprehensive Go toolkit for the Massa blockchain ecosystem, providing smart contract development capabilities, client libraries, and wallet functionality.

## 🚀 Features

- **Smart Contract SDK** (`sc/`): Pure Go SDK for developing WebAssembly-compatible smart contracts
- **Blockchain Client** (`client/`): Interface to communicate with the Massa blockchain
- **Wallet Management** (`wallet/`): Account creation and transaction signing utilities
- **Complete Examples**: End-to-end examples of smart contract development and deployment


## Requirements

- GO go 1.25+

## Installation

```bash
go get github.com/nafsilabs/massa-go
```


## 📦 Components

### 🔧 Smart Contract SDK (`sc/`)

A pure Go SDK for developing smart contracts on the Massa blockchain that compile to WebAssembly. This SDK provides a comprehensive set of functions for interacting with the Massa blockchain, similar to the official [massa-as-sdk](https://github.com/massalabs/massa-as-sdk) but written in native Go.

**Key Features:**
- **Pure Go Implementation**: Write smart contracts in native Go syntax
- **WebAssembly Compilation**: Compiles to WASM compatible with Massa blockchain
- **Complete API Coverage**: Implements all major Massa blockchain functions
- **Type Safety**: Full Go type system with compile-time safety
- **Memory Management**: Automatic WebAssembly memory handling
- **Testing Support**: Mock implementations for local development and testing

**Quick Start:**
```go
package main

import massa "github.com/nafsilabs/massa-go/sc"

func main() {
    // Generate an event
    massa.GenerateEvent("Hello, Massa!")
    
    // Store data
    massa.SetString("greeting", "Hello from Go!")
    
    // Get caller information
    caller := massa.Caller()
    massa.Printf("Called by: %s", caller.String())
}
```

### 🌐 Blockchain Client (`client/`)

A comprehensive node client that enables you to communicate with the Massa blockchain. It offers an interface to retrieve data directly from the blockchain, interact with smart contracts, acquire and monitor events, and perform additional actions.

**Features:**
- Blockchain data retrieval
- Smart contract deployment
- Smart contract interaction
- Event monitoring
- Transaction broadcasting
- Connection management with automatic reconnection
- WebSocket and HTTP API support
- Multi-client connection pooling

**Quick Start:**
```go
package main

import "github.com/nafsilabs/massa-go/client"

func main() {
    // Connect to Massa testnet
    massaClient, err := client.NewTestnetClient()
    if err != nil {
        panic(err)
    }
    defer massaClient.Close()
    
    // Get node status
    status, err := massaClient.GetNodeStatus()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Connected to node: %s (slot %d)\n", 
        status.NodeID, status.CurrentSlot)
    
    // Subscribe to events
    eventHandler := func(event *client.Event) {
        fmt.Printf("Event: %s\n", event.Data)
    }
    
    subscriptionID, err := massaClient.SubscribeToEvents(
        &client.EventFilter{}, eventHandler)
    if err != nil {
        panic(err)
    }
    
    // Deploy a smart contract
    deployResult, err := massaClient.DeployContract(
        &client.ContractDeploymentRequest{
            Bytecode: contractBytecode,
            MaxGas:   1000000,
            Coins:    "0",
        }, 
        account, 
        "password",
    )
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Contract deployed: %s\n", 
        deployResult.ContractAddress)
}
```

*Source code inspired by: https://github.com/massalabs/station/tree/main/pkg* and https://github.com/massalabs/massa-web3
*Smart contract deployment inspired by: https://github.com/massalabs/massa-standards/blob/main/smart-contracts/assembly/contracts/deployer/deployer.ts*

### 💼 Wallet Management (`wallet/`)

Wallet functionality for creating accounts with public and private keys used to sign Massa blockchain transactions.

**Features:**
- Account generation
- Key management
- Transaction signing
- Address derivation

*Some implementation ideas from: https://github.com/massalabs/station/tree/main/pkg*

### 📚 Examples (`examples/`)

Complete examples demonstrating:
- Smart contract development
- Wallet account creation
- Client usage for deployment
- Smart contract interaction

## 🏃‍♂️ Quick Start

### 1. Initialize your project

```bash
go mod init your-massa-project
go get github.com/nafsilabs/massa-go
```

### 2. Create a Smart Contract

```go
// main.go
package main

import massa "github.com/nafsilabs/massa-go/sc"

func main() {
    massa.GenerateEvent("Hello, Massa!")
    massa.SetString("greeting", "Hello from Go!")
}
```

### 3. Build for WebAssembly

```bash
export GOOS=js
export GOARCH=wasm
go build -ldflags="-s -w" -o contract.wasm main.go
```

### 4. Deploy and Interact

```go
package main

import (
    "context"
    "log"

    "github.com/nafsilabs/massa-go"
)

func main() {
    // Create a Massa gRPC client
    client, err := massa.NewMassaClient(&massa.ClientConfig{
        Address: "testnet.massa.net:33035",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    ctx := context.Background()

    // Use wallet and client together (see examples/ for full flows)
    // ...

    _ = ctx
}
```

## 🏗️ Project Structure

```
massa-go/
├── massa.go               # High-level SDK entrypoint
├── sc/                    # Smart Contract SDK
│   ├── sc.go              # WebAssembly imports
│   ├── address.go         # Address operations
│   ├── storage.go         # Storage functions
│   ├── context.go         # Execution context
│   ├── contract.go        # Contract interactions
│   ├── coins.go           # Balance and transfers
│   ├── events.go          # Events and logging
│   ├── crypto.go          # Cryptographic functions
│   └── op_datastore.go    # Operation datastore
├── client/                # Blockchain client
├── wallet/                # Wallet management
├── examples/
│   └── smart_contract/    # Example smart contract
│       ├── main.go        # Contract implementation
│       ├── build.sh       # Build script
│       └── README.md      # Example documentation
└── README.md             # This file
```

## 📖 Documentation

### Smart Contract SDK

The `sc/` package provides comprehensive smart contract development capabilities:

- **Address Operations**: Validation, creation, and utilities
- **Storage Management**: Persistent key-value storage
- **Context Access**: Execution environment information
- **Contract Interactions**: Calling other contracts and bytecode operations
- **Coin Operations**: Balance queries and transfers
- **Event System**: Structured event generation and logging
- **Cryptography**: Hashing, signatures, and utilities
- **Operation Datastore**: Access to operation-level data

### API Modules

| Module | Description |
|--------|-------------|
| `Address` | Address validation, creation, and utilities |
| `Storage` | Persistent key-value storage operations |
| `Context` | Execution context (caller, timestamp, gas, etc.) |
| `Contract` | Contract calls, bytecode operations, messaging |
| `Coins` | Balance queries and coin transfers |
| `Events` | Event generation and structured logging |
| `Crypto` | Cryptographic functions (hashing, signatures) |
| `OpDatastore` | Operation-level data access |

## 🧪 Testing

The SDK provides mock implementations for testing:

```go
func TestContract(t *testing.T) {
    // The SDK automatically uses mocks when not in WASM runtime
    caller := massa.Caller() // Returns mock address
    massa.SetString("test", "value")
    value := massa.GetString("test")
    // Test your contract logic...
}
```

## 📦 Publishing

Create an annotated tag and push it to the remote:

```bash
# create an annotated tag
git tag -a vX.Y.Z -m "Release vX.Y.Z"

# push the specific tag to origin
git push origin vX.Y.Z

# (optional) push all local tags
git push --tags
```

## 🔐 Security Best Practices

1. **Input Validation**: Always validate addresses, amounts, and parameters
2. **Access Control**: Implement proper permission checks
3. **Error Handling**: Handle errors gracefully and don't expose internals
4. **Gas Management**: Be mindful of gas costs for operations
5. **State Consistency**: Use atomic operations for multi-step changes

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 References

- [massa-as-sdk](https://github.com/massalabs/massa-as-sdk) - Official AssemblyScript SDK
- [massa-as-sdk Documentation](https://as-sdk.docs.massa.net/index.html) - SDK documentation
- [Massa](https://massa.net/) - Massa blockchain
- [Massa Documentation](https://docs.massa.net/) - Official documentation

## 📞 Support

- [GitHub Issues](https://github.com/nafsilabs/massa-go/issues)
- [Massa Discord](https://discord.gg/massa)
- [Documentation](https://docs.massa.net/)

---

Built with ❤️ for the Massa ecosystem

