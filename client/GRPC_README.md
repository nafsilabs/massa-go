# Massa gRPC Client

A comprehensive Go client for the Massa blockchain gRPC API with full support for unary calls and streaming (both server-streaming and bidirectional streaming).

## Features

- **Public Service Client**: Access to all public blockchain data and operations
- **Private Service Client**: Administrative and node management operations
- **Streaming Support**:
  - Server Streaming (unidirectional from server)
  - Bidirectional Streaming (two-way communication)
- **Connection Management**: Easy connection setup with TLS support
- **Context Management**: Built-in timeout and cancellation support

## Installation

```bash
go get massa
```

## Quick Start

### Creating a Client

```go
package main

import (
    "time"
    "massa/client"
)

func main() {
    config := client.ClientConfig{
        Address:        "buildnet.massa.net:33037",
        UseTLS:         false,
        DefaultTimeout: 30 * time.Second,
    }

    grpcClient, err := client.NewMassaGrpcClient(config)
    if err != nil {
        panic(err)
    }
    defer grpcClient.Close()
}
```

### Unary Calls

Unary calls are simple request-response operations:

```go
ctx, cancel := grpcClient.GetContext()
defer cancel()

// Get blockchain status
status, err := grpcClient.PublicClient.GetStatus(ctx, &pb.GetStatusRequest{})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Status: %+v\n", status)

// Get blocks
blocks, err := grpcClient.PublicClient.GetBlocks(ctx, &pb.GetBlocksRequest{
    BlockIds: []string{"block-id-1", "block-id-2"},
})
if err != nil {
    log.Fatal(err)
}

// Execute read-only call
result, err := grpcClient.PublicClient.ExecuteReadOnlyCall(ctx, &pb.ExecuteReadOnlyCallRequest{
    MaxGas: 1000000,
    Call: &pb.ReadOnlyCall{
        // ... call parameters
    },
})
```

### Server Streaming

Server streaming allows you to subscribe to events and receive a continuous stream of updates:

```go
import (
    "context"
    "time"
    pb "massa/client/proto/massa/api/v1"
)

// Stream new operations
ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
defer cancel()

err := grpcClient.NewOperationsServerStream(ctx, 
    &pb.NewOperationsServerRequest{
        Filters: []*pb.NewOperationsServerRequest_OperationsFilter{},
    }, 
    func(msg *pb.NewOperationsServerResponse) error {
        fmt.Printf("New operation: %+v\n", msg)
        return nil
    },
)

// Stream new blocks
err = grpcClient.NewBlocksServerStream(ctx,
    &pb.NewBlocksServerRequest{
        Filters: []*pb.NewBlocksServerRequest_BlocksFilter{},
    },
    func(msg *pb.NewBlocksServerResponse) error {
        fmt.Printf("New block: %+v\n", msg)
        return nil
    },
)

// Stream new filled blocks (blocks with operations)
err = grpcClient.NewFilledBlocksServerStream(ctx,
    &pb.NewFilledBlocksServerRequest{
        Filters: []*pb.NewFilledBlocksServerRequest_FilledBlocksFilter{},
    },
    func(msg *pb.NewFilledBlocksServerResponse) error {
        fmt.Printf("New filled block: %+v\n", msg)
        return nil
    },
)

// Stream slot execution outputs
err = grpcClient.NewSlotExecutionOutputsServerStream(ctx,
    &pb.NewSlotExecutionOutputsServerRequest{
        Filters: []*pb.NewSlotExecutionOutputsServerRequest_SlotExecutionOutputsFilter{},
    },
    func(msg *pb.NewSlotExecutionOutputsServerResponse) error {
        fmt.Printf("Slot execution output: %+v\n", msg)
        return nil
    },
)

// Stream transactions throughput
err = grpcClient.TransactionsThroughputServerStream(ctx,
    &pb.TransactionsThroughputServerRequest{},
    func(msg *pb.TransactionsThroughputServerResponse) error {
        fmt.Printf("Throughput: %+v\n", msg)
        return nil
    },
)
```

### Bidirectional Streaming

Bidirectional streaming allows both client and server to send messages independently:

```go
ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
defer cancel()

// Create bidirectional stream
stream, err := grpcClient.NewOperationsStream(ctx)
if err != nil {
    log.Fatal(err)
}

// Send subscription request
err = stream.Send(&pb.NewOperationsRequest{
    Filters: []*pb.NewOperationsRequest_OperationsFilter{
        // Add your filters here
    },
})
if err != nil {
    log.Fatal(err)
}

// Receive messages in a goroutine
go func() {
    for {
        resp, err := stream.Recv()
        if err != nil {
            if err == io.EOF {
                return
            }
            log.Printf("Receive error: %v", err)
            return
        }
        fmt.Printf("Received: %+v\n", resp)
    }
}()

// Send more messages as needed
err = stream.Send(&pb.NewOperationsRequest{
    // Another request
})

// Close the send direction when done
stream.CloseSend()

// Wait for context to finish
<-ctx.Done()
```

### Available Bidirectional Streams

```go
// Operations
stream, err := grpcClient.NewOperationsStream(ctx)
stream, err := grpcClient.SendOperationsStream(ctx)

// Blocks
stream, err := grpcClient.NewBlocksStream(ctx)
stream, err := grpcClient.SendBlocksStream(ctx)
stream, err := grpcClient.NewFilledBlocksStream(ctx)

// Endorsements
stream, err := grpcClient.NewEndorsementsStream(ctx)
stream, err := grpcClient.SendEndorsementsStream(ctx)

// Execution & Transfers
stream, err := grpcClient.NewSlotExecutionOutputsStream(ctx)
stream, err := grpcClient.NewSlotABICallStacksStream(ctx)
stream, err := grpcClient.NewSlotTransfersStream(ctx)

// Throughput
stream, err := grpcClient.TransactionsThroughputStream(ctx)
```

## Public Service Methods

### Unary Methods

- `ExecuteReadOnlyCall`: Execute a read-only smart contract call
- `GetBlocks`: Get blocks by IDs
- `GetDatastoreEntries`: Get datastore entries
- `GetEndorsements`: Get endorsements by IDs
- `GetNextBlockBestParents`: Get next block best parents
- `GetOperations`: Get operations by IDs
- `GetScExecutionEvents`: Get smart contract execution events
- `GetSelectorDraws`: Get selector draws
- `GetStakers`: Get stakers information
- `GetStatus`: Get node status
- `GetTransactionsThroughput`: Get transactions throughput
- `QueryState`: Query blockchain state
- `SearchBlocks`: Search for blocks
- `SearchEndorsements`: Search for endorsements
- `SearchOperations`: Search for operations
- `GetOperationABICallStacks`: Get ABI call stack of an operation
- `GetSlotABICallStacks`: Get ABI call stacks for a slot
- `GetSlotTransfers`: Get all transfers for a given slot

### Server Streaming Methods

- `NewBlocksServerStream`: Subscribe to new blocks
- `NewEndorsementsServerStream`: Subscribe to new endorsements
- `NewFilledBlocksServerStream`: Subscribe to new blocks with operations
- `NewOperationsServerStream`: Subscribe to new operations
- `NewSlotExecutionOutputsServerStream`: Subscribe to slot execution outputs
- `TransactionsThroughputServerStream`: Subscribe to transaction throughput updates
- `NewTransfersInfoServerStream`: Subscribe to transfer information

### Bidirectional Streaming Methods

- `NewBlocksStream`: Bidirectional stream for blocks
- `NewEndorsementsStream`: Bidirectional stream for endorsements
- `NewFilledBlocksStream`: Bidirectional stream for filled blocks
- `NewOperationsStream`: Bidirectional stream for operations
- `NewSlotExecutionOutputsStream`: Bidirectional stream for slot execution outputs
- `NewSlotABICallStacksStream`: Bidirectional stream for ABI call stacks
- `NewSlotTransfersStream`: Bidirectional stream for slot transfers
- `SendBlocksStream`: Bidirectional stream for sending blocks
- `SendEndorsementsStream`: Bidirectional stream for sending endorsements
- `SendOperationsStream`: Bidirectional stream for sending operations
- `TransactionsThroughputStream`: Bidirectional stream for throughput monitoring

## Private Service Methods

Access administrative and node management operations:

```go
// Get node status
status, err := grpcClient.PrivateClient.GetNodeStatus(ctx, &pb.GetNodeStatusRequest{})

// Add staking keys
resp, err := grpcClient.PrivateClient.AddStakingSecretKeys(ctx, &pb.AddStakingSecretKeysRequest{
    SecretKeys: []string{"secret-key-1"},
})

// Ban nodes
resp, err := grpcClient.PrivateClient.BanNodesByIds(ctx, &pb.BanNodesByIdsRequest{
    NodeIds: []string{"node-id-1"},
})

// Get MIP status
mipStatus, err := grpcClient.PrivateClient.GetMipStatus(ctx, &pb.GetMipStatusRequest{})
```

## Configuration Options

```go
type ClientConfig struct {
    Address        string            // gRPC server address (e.g., "buildnet.massa.net:33037")
    UseTLS         bool              // Enable TLS encryption
    TLSConfig      *tls.Config       // Custom TLS configuration (optional)
    DefaultTimeout time.Duration     // Default timeout for operations
    DialOptions    []grpc.DialOption // Additional gRPC dial options
}
```

## Network Endpoints

- **Mainnet**: `mainnet.massa.net:33037`
- **Buildnet**: `buildnet.massa.net:33037`

## Error Handling

All methods return errors that should be checked:

```go
resp, err := grpcClient.PublicClient.GetStatus(ctx, &pb.GetStatusRequest{})
if err != nil {
    // Handle error
    log.Printf("Error: %v", err)
    return
}
```

For streaming methods, handle both connection errors and message processing errors:

```go
err := grpcClient.NewBlocksServerStream(ctx, req, func(msg *pb.NewBlocksServerResponse) error {
    // Handle message
    if someCondition {
        return fmt.Errorf("processing error")
    }
    return nil
})

if err != nil && ctx.Err() != context.DeadlineExceeded {
    log.Printf("Stream error: %v", err)
}
```

## Context Management

The client provides helper methods for context creation:

```go
// Use default timeout
ctx, cancel := grpcClient.GetContext()
defer cancel()

// Use custom timeout
ctx, cancel := grpcClient.GetContextWithTimeout(5 * time.Minute)
defer cancel()

// Use context with deadline
ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour))
defer cancel()
```

## Examples

See the `examples/grpc_stream/` directory for complete working examples:

- Basic unary calls
- Server streaming subscriptions
- Bidirectional streaming communication
- Error handling patterns

## License

This client is part of the massa-go project.
