# Massa Go Client

A comprehensive Go client library for interacting with the Massa blockchain. This client provides full access to blockchain data, smart contract deployment and interaction, transaction broadcasting, and real-time event monitoring.

## 🚀 Features

- **Blockchain Data Retrieval**: Access blocks, transactions, balances, and node information
- **Smart Contract Deployment**: Deploy single or multiple contracts with constructor support
- **Smart Contract Interaction**: Call contract functions and read contract state
- **Transaction Broadcasting**: Send transfers, buy/sell rolls, and other operations
- **Event Monitoring**: Real-time event subscription via WebSocket and polling
- **Connection Management**: Automatic reconnection, retry logic, and connection pooling
- **Network Support**: Mainnet, testnet, and custom network configurations

## 📦 Installation

```bash
go get github.com/nafsilabs/massa-go/client
```

## 🔧 Dependencies

- `github.com/gorilla/websocket` - WebSocket client
- `github.com/nafsilabs/massa-go/wallet` - Wallet functionality
- `golang.org/x/crypto` - Cryptographic functions

## 📖 Quick Start

### Basic Client Setup

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/nafsilabs/massa-go/client"
)

func main() {
    // Create a client for testnet
    massaClient, err := client.NewTestnetClient()
    if err != nil {
        log.Fatal(err)
    }
    defer massaClient.Close()
    
    // Get node status
    status, err := massaClient.GetNodeStatus()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Connected to Massa node: %s\n", status.NodeID)
    fmt.Printf("Current slot: %d\n", status.CurrentSlot)
}
```

### Client Builder Pattern

```go
// Create a custom client with specific configuration
client, err := client.NewClientBuilder().
    WithTestnet().
    WithTimeout(60 * time.Second).
    WithRetryCount(5).
    WithAPIKey("your-api-key").
    BuildAndConnect()
```

## 🌐 Network Configurations

### Predefined Networks

```go
// Mainnet
client, err := client.NewMainnetClient()

// Testnet  
client, err := client.NewTestnetClient()

// Local development node
client, err := client.NewLocalnetClient()
```

### Custom Networks

```go
customNetwork := &client.NetworkConfig{
    Name:      "custom",
    ChainID:   12345,
    RPCURL:    "https://your-node.com/api/v2",
    WSURL:     "wss://your-node.com/ws",
    PublicAPI: "https://your-node.com/api/v2",
}

client := client.NewClient(client.DefaultClientConfig(customNetwork))
```

## 🔍 Blockchain Data Retrieval

### Node Information

```go
// Get node status
status, err := client.GetNodeStatus()
fmt.Printf("Current slot: %d\n", status.CurrentSlot)

// Get node configuration
config, err := client.GetNodeInfo()
fmt.Printf("Chain ID: %d\n", config.ChainID)
```

### Blocks and Operations

```go
// Get block by ID
block, err := client.GetBlock("block_id_here")

// Get multiple blocks
blocks, err := client.GetBlocks([]string{"id1", "id2"})

// Get blocks by slot numbers
blocks, err := client.GetBlocksBySlots([]uint64{100, 101, 102})

// Get operation information
operation, err := client.GetOperation("operation_id_here")
```

### Account Information

```go
// Get account balance
balance, err := client.GetBalance("AU1234...")
fmt.Printf("Final balance: %s MAS\n", balance.Final)

// Get multiple balances
addresses := []string{"AU1234...", "AU5678..."}
balances, err := client.GetBalances(addresses)
```

## 💼 Smart Contract Operations

### Contract Deployment

```go
// Deploy a single contract
deployRequest := &client.ContractDeploymentRequest{
    Bytecode: contractBytecode,
    MaxGas:   1000000,
    Coins:    "0",
    Datastore: map[string][]byte{
        "initial_value": []byte("Hello, Massa!"),
    },
    Constructor: &client.ConstructorCall{
        Function:  "constructor",
        Parameter: []byte(`{"owner": "AU1234..."}`),
        Coins:     "0",
    },
}

result, err := client.DeployContract(deployRequest, account, "password")
fmt.Printf("Contract deployed at: %s\n", result.ContractAddress)
```

### Multiple Contract Deployment

```go
// Deploy multiple contracts in one transaction
multiRequest := &client.MultiContractDeploymentRequest{
    Contracts: []client.ContractDeploymentRequest{
        {Bytecode: contract1Bytecode, MaxGas: 500000, Coins: "0"},
        {Bytecode: contract2Bytecode, MaxGas: 500000, Coins: "0"},
    },
    MaxGas: 1000000,
}

result, err := client.DeployMultipleContracts(multiRequest, account, "password")
fmt.Printf("Deployed %d contracts\n", len(result.ContractAddresses))
```

### Contract Interaction

```go
// Call a contract function
callResult, err := client.CallContract(
    "AS1234...",                    // contract address
    "transfer",                     // function name
    []byte(`{"to": "AU5678...", "amount": 1000}`), // parameters
    100000,                         // max gas
    "0",                           // coins to send
    account,                       // signing account
    "password",                    // account password
)

// Read contract state (no transaction required)
readResult, err := client.ReadContract(
    "AS1234...",                    // contract address
    "get_balance",                  // function name
    []byte(`{"address": "AU1234..."}`), // parameters
    "AU1234...",                   // caller address
)
fmt.Printf("Contract returned: %s\n", string(readResult.Result))
```

## 💸 Transaction Operations

### Coin Transfers

```go
// Transfer coins
result, err := client.TransferCoins(
    "AU1234...",    // from address
    "AU5678...",    // to address
    "1000000000",   // amount in nanoMAS (1 MAS)
    account,        // signing account
    "password",     // account password
)

fmt.Printf("Transfer transaction: %s\n", result.OperationID)
```

### Validator Operations

```go
// Buy validator rolls
result, err := client.BuyRolls(
    "AU1234...",    // from address
    5,              // number of rolls
    account,        // signing account
    "password",     // account password
)

// Sell validator rolls
result, err := client.SellRolls(
    "AU1234...",    // from address
    2,              // number of rolls
    account,        // signing account
    "password",     // account password
)
```

### Batch Transactions

```go
// Send multiple transactions
operations := []*client.Operation{
    {
        Creator: account.Address.String(),
        Fee: "1000000",
        ExpirePeriod: 5,
        Type: "Transfer",
        Content: &client.TransferOperation{
            Recipient: "AU5678...",
            Amount: "500000000",
        },
    },
    // ... more operations
}

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

results, err := client.BatchTransactions(ctx, operations, account, "password")
```

## 📡 Event Monitoring

### WebSocket Event Subscription

```go
// Subscribe to all events
eventHandler := func(event *client.Event) {
    fmt.Printf("Event from slot %d: %s\n", event.Context.Slot, event.Data)
}

subscriptionID, err := client.SubscribeToEvents(&client.EventFilter{}, eventHandler)

// Subscribe to contract-specific events
subscriptionID, err := client.SubscribeToContractEvents("AS1234...", eventHandler)

// Subscribe to address-specific events
subscriptionID, err := client.SubscribeToAddressEvents("AU1234...", eventHandler)

// Unsubscribe when done
err = client.UnsubscribeFromEvents(subscriptionID)
```

### Event Filtering

```go
// Filter events by slot range
startSlot := uint64(1000)
endSlot := uint64(2000)
filter := &client.EventFilter{
    Start:   &startSlot,
    End:     &endSlot,
    Emitter: &contractAddress,
}

events, err := client.GetEvents(filter)
```

### Channel-Based Event Streaming

```go
// Create an event stream
stream, err := client.NewEventStream(&client.EventFilter{})
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

// Listen to events
for {
    select {
    case event := <-stream.Events():
        fmt.Printf("Received event: %+v\n", event)
    case err := <-stream.Errors():
        fmt.Printf("Stream error: %v\n", err)
    case <-time.After(30 * time.Second):
        return // Timeout
    }
}
```

### Polling-Based Event Monitoring

```go
// For environments without WebSocket support
poller := client.NewEventPoller(
    &client.EventFilter{},
    5*time.Second,  // polling interval
    eventHandler,
)

poller.Start()
defer poller.Stop()
```

## 🔄 Connection Management

### Connection Pool

```go
// Create multiple client configurations for load balancing
configs := []*client.ClientConfig{
    client.DefaultClientConfig(client.MainnetConfig),
    client.DefaultClientConfig(client.MainnetConfig),
    client.DefaultClientConfig(client.MainnetConfig),
}

// Create client pool
pool, err := client.NewClientPool(configs)
if err != nil {
    log.Fatal(err)
}
defer pool.Close()

// Get a healthy client
ctx := context.Background()
healthyClient, err := pool.GetHealthyClient(ctx)
if err != nil {
    log.Fatal(err)
}

// Use the client
status, err := healthyClient.GetNodeStatus()
```

### Health Monitoring

```go
// Check client health
if client.IsConnected() {
    fmt.Println("Client is connected")
}

// Force reconnection
err := client.Reconnect()

// Pool health check
healthResults := pool.HealthCheck(context.Background())
for i, err := range healthResults {
    if err != nil {
        fmt.Printf("Client %d is unhealthy: %v\n", i, err)
    }
}
```

## 🛠️ Advanced Usage

### Custom Request Handling

```go
// Use retry wrapper for critical operations
ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
defer cancel()

err := client.WithRetry(ctx, 5, func() error {
    _, err := client.GetNodeStatus()
    return err
})
```

### Configuration Management

```go
// Get current configuration
config := client.GetConfig()
fmt.Printf("Current network: %s\n", config.Network.Name)

// Update configuration
newConfig := client.DefaultClientConfig(client.MainnetConfig)
newConfig.Timeout = 60 * time.Second
err := client.SetConfig(newConfig)
```

## 📊 Monitoring and Debugging

### Transaction Monitoring

```go
// Wait for transaction completion
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result, err := client.WaitForTransaction(ctx, operationID)
if err != nil {
    log.Printf("Transaction failed: %v", err)
} else {
    fmt.Printf("Transaction completed with status: %s\n", result.Status)
}

// Get transaction receipt
receipt, err := client.GetTransactionReceipt(operationID)
fmt.Printf("Gas used: %d\n", receipt.GasUsed)
```

### Error Handling

```go
// Check for specific error types
if rpcErr, ok := err.(*client.RPCError); ok {
    fmt.Printf("RPC error code: %d, message: %s\n", rpcErr.Code, rpcErr.Message)
}

// Operation status checking
switch result.Status {
case client.OperationStatusSuccess:
    fmt.Println("Operation succeeded")
case client.OperationStatusFailed:
    fmt.Printf("Operation failed: %s\n", result.Error)
case client.OperationStatusPending:
    fmt.Println("Operation is pending")
}
```

## 🔐 Security Best Practices

1. **API Keys**: Store API keys securely and never commit them to version control
2. **Password Management**: Use secure password storage for wallet account passwords
3. **Network Validation**: Always validate network configurations in production
4. **Error Handling**: Implement proper error handling for all network operations
5. **Resource Cleanup**: Always close clients and streams when done

## 🧪 Testing

### Unit Testing

```go
func TestClientConnection(t *testing.T) {
    client, err := client.NewTestnetClient()
    require.NoError(t, err)
    defer client.Close()
    
    status, err := client.GetNodeStatus()
    require.NoError(t, err)
    require.NotEmpty(t, status.NodeID)
}
```

### Integration Testing

```go
func TestSmartContractDeployment(t *testing.T) {
    client, err := client.NewTestnetClient()
    require.NoError(t, err)
    defer client.Close()
    
    // Deploy test contract
    deployRequest := &client.ContractDeploymentRequest{
        Bytecode: loadTestContract(),
        MaxGas:   1000000,
        Coins:    "0",
    }
    
    result, err := client.DeployContract(deployRequest, testAccount, "password")
    require.NoError(t, err)
    require.NotEmpty(t, result.ContractAddress)
}
```

## 📝 Examples

See the [example](./example/) directory for complete working examples:

- Basic client usage
- Smart contract deployment and interaction  
- Event monitoring and filtering
- Transaction operations
- Connection pooling
- Error handling patterns

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.

## 🔗 Related Projects

- [massa-go/sc](../sc) - Smart Contract SDK
- [massa-go/wallet](../wallet) - Wallet Management
- [Massa](https://massa.net/) - Massa Blockchain

## 🐛 Troubleshooting

### Common Issues

**Connection Errors**
```bash
# Check if the node is running and accessible
curl -X POST -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"get_status","id":1}' \
  https://test.massa.net/api/v2
```

**WebSocket Issues**
```go
// Disable WebSocket if having connectivity issues
client := client.NewClientBuilder().
    WithTestnet().
    WithCustomWebSocket(""). // Empty WebSocket URL disables it
    Build()
```

**Timeout Issues**
```go
// Increase timeout for slow connections
client := client.NewClientBuilder().
    WithTimeout(120 * time.Second).
    WithRetryCount(10).
    Build()
```

---

Built with ❤️ for the Massa ecosystem