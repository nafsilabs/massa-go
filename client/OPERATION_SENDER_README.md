# Massa Client - Operation Sender

This document describes how to use the `OperationSender` to send different types of operations to the Massa blockchain network.

## Overview

The `OperationSender` provides a high-level abstraction for sending operations to the Massa network. It handles:
- Operation serialization
- Signing workflow
- gRPC streaming communication
- Error handling

## Supported Operations

1. **Transaction** - Transfer coins between addresses
2. **BuyRolls** - Purchase rolls for staking
3. **SellRolls** - Sell rolls
4. **ExecuteSC** - Deploy smart contract bytecode
5. **CallSC** - Call a smart contract function

## Basic Usage

### Initialize the Client

```go
import (
	"context"
	"time"
	
	"github.com/nafsilabs/massa-go/client"
)

// Create gRPC client
grpcClient, err := client.NewMassaGrpcClient(client.ClientConfig{
	Address:        "buildnet.massa.net:33037",
	UseTLS:         false,
	DefaultTimeout: 30 * time.Second,
})
if err != nil {
	log.Fatal(err)
}
defer grpcClient.Close()

// Create operation sender
opSender := client.NewOperationSender(grpcClient)
```

### Send a Transaction

```go
ctx := context.Background()

// You need to sign the operation first (using your wallet)
signature := "your_base64_signature"
publicKey := "your_base64_public_key"

operationID, err := opSender.SendTransaction(
	ctx,
	"AU12recipient address",         // recipient address
	1000000,                          // amount in NanoMassa
	100,                              // fee
	10,                               // expire period
	signature,
	publicKey,
)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Transaction sent! Operation ID: %s\n", operationID)
```

### Buy Rolls

```go
operationID, err := opSender.SendBuyRolls(
	ctx,
	1,                                // number of rolls to buy
	100,                              // fee
	10,                               // expire period
	signature,
	publicKey,
)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Rolls purchased! Operation ID: %s\n", operationID)
```

### Sell Rolls

```go
operationID, err := opSender.SendSellRolls(
	ctx,
	1,                                // number of rolls to sell
	100,                              // fee
	10,                               // expire period
	signature,
	publicKey,
)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Rolls sold! Operation ID: %s\n", operationID)
```

### Deploy Smart Contract (Execute SC)

```go
import (
	model "github.com/nafsilabs/massa-go/client/proto/massa/model/v1"
)

// Read compiled smart contract bytecode
bytecode, err := os.ReadFile("contract.wasm")
if err != nil {
	log.Fatal(err)
}

// Optional: prepare datastore entries
datastore := []*model.BytesMapFieldEntry{
	{
		Key:   []byte("config_key"),
		Value: []byte("config_value"),
	},
}

operationID, err := opSender.SendExecuteSC(
	ctx,
	bytecode,                         // smart contract bytecode
	1000000,                          // max gas
	100,                              // max coins
	datastore,                        // datastore entries (can be nil)
	100,                              // fee
	10,                               // expire period
	signature,
	publicKey,
)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Smart contract deployed! Operation ID: %s\n", operationID)
```

### Call Smart Contract

```go
// Prepare function parameters (ABI encoded)
parameters := []byte{/* ABI encoded parameters */}

operationID, err := opSender.SendCallSC(
	ctx,
	"AS12contract_address",           // target smart contract address
	"transfer",                       // function name to call
	parameters,                       // function parameters
	100000,                           // max gas
	0,                                // coins to send with call
	100,                              // fee
	10,                               // expire period
	signature,
	publicKey,
)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Smart contract called! Operation ID: %s\n", operationID)
```

## Advanced: Send Multiple Operations

```go
operations := []struct {
	Op        client.MassaOperation
	Signature string
	PublicKey string
}{
	{
		Op: &client.TransactionOp{
			RecipientAddress: "AU12recipient1",
			Amount:           1000000,
		},
		Signature: "signature1",
		PublicKey: "pubkey1",
	},
	{
		Op: &client.BuyRollsOp{
			RollCount: 1,
		},
		Signature: "signature2",
		PublicKey: "pubkey2",
	},
}

operationIDs, err := opSender.SendOperations(
	ctx,
	operations,
	100,                              // fee for all operations
	10,                               // expire period for all
)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Sent %d operations: %v\n", len(operationIDs), operationIDs)
```

## Custom Operations

You can create custom operations by implementing the `MassaOperation` interface:

```go
type MassaOperation interface {
	ToOperationType() *model.OperationType
}
```

Example:

```go
type MyCustomOp struct {
	Data string
}

func (op *MyCustomOp) ToOperationType() *model.OperationType {
	// Return the appropriate OperationType
	return &model.OperationType{
		Type: &model.OperationType_Transaction{
			Transaction: &model.Transaction{
				// ... your operation data
			},
		},
	}
}

// Use it
customOp := &MyCustomOp{Data: "test"}
opID, err := opSender.SendOperation(ctx, customOp, fee, expirePeriod, signature, publicKey)
```

## Signing Operations

**Important:** The `OperationSender` does NOT handle signing. You must sign operations yourself before sending them.

The signature and public key parameters should be:
- **Signature**: Base58 or hex encoded signature of the serialized operation
- **PublicKey**: Base58 or hex encoded public key of the signer

For signing, you'll typically use the wallet package (when available) or your own key management solution.

## Helper Functions

### NewNativeAmount

Converts a uint64 amount to the protobuf NativeAmount structure:

```go
amount := client.NewNativeAmount(1000000)
// Returns: &model.NativeAmount{Mantissa: 1000000, Scale: 0}
```

### EncodeBase64 / DecodeBase64

```go
encoded := client.EncodeBase64([]byte("data"))
decoded, err := client.DecodeBase64(encoded)
```

## Operation Types Reference

### TransactionOp
```go
type TransactionOp struct {
	RecipientAddress string  // AU12...
	Amount           uint64  // Amount in NanoMassa
}
```

### BuyRollsOp
```go
type BuyRollsOp struct {
	RollCount uint64  // Number of rolls to purchase
}
```

### SellRollsOp
```go
type SellRollsOp struct {
	RollCount uint64  // Number of rolls to sell
}
```

### ExecuteSCOp
```go
type ExecuteSCOp struct {
	Bytecode  []byte                    // Compiled smart contract bytecode
	MaxGas    uint64                    // Maximum gas allowed
	MaxCoins  uint64                    // Maximum coins that can be spent
	Datastore []*model.BytesMapFieldEntry  // Key-value datastore (optional)
}
```

### CallSCOp
```go
type CallSCOp struct {
	TargetAddress  string  // AS12... (smart contract address)
	TargetFunction string  // Function name to call
	Parameter      []byte  // ABI-encoded parameters
	MaxGas         uint64  // Maximum gas allowed
	Coins          uint64  // Coins to send with the call
}
```

## Error Handling

The `OperationSender` returns detailed errors:

```go
operationID, err := opSender.SendTransaction(...)
if err != nil {
	if strings.Contains(err.Error(), "operation error:") {
		// Massa network rejected the operation
		fmt.Println("Operation rejected by network:", err)
	} else if strings.Contains(err.Error(), "failed to open stream") {
		// Network/connection error
		fmt.Println("Connection error:", err)
	} else {
		// Other errors (serialization, etc.)
		fmt.Println("Error:", err)
	}
	return
}
```

## Complete Example

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nafsilabs/massa-go/client"
)

func main() {
	// Initialize gRPC client
	grpcClient, err := client.NewMassaGrpcClient(client.ClientConfig{
		Address:        "buildnet.massa.net:33037",
		UseTLS:         false,
		DefaultTimeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer grpcClient.Close()

	// Create operation sender
	opSender := client.NewOperationSender(grpcClient)

	// Send a transaction
	// NOTE: In production, you would sign this operation with your private key
	signature := "your_signature_here"
	publicKey := "your_public_key_here"

	ctx := context.Background()
	operationID, err := opSender.SendTransaction(
		ctx,
		"AU12recipient_address_here",
		1000000,  // 1 Massa (1,000,000 NanoMassa)
		100,      // Fee
		10,       // Expire in 10 periods
		signature,
		publicKey,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Transaction sent successfully!\n")
	fmt.Printf("Operation ID: %s\n", operationID)
}
```

## Notes

1. **NativeAmount**: Amounts are represented with mantissa and scale. The helper `NewNativeAmount(amount)` creates an amount with scale 0.

2. **Expire Period**: Operations expire after the specified number of periods. Choose appropriately based on network conditions.

3. **Gas Limits**: For smart contract operations, ensure MaxGas is sufficient for your operation to complete.

4. **Bidirectional Streaming**: Under the hood, `SendOperations` uses gRPC bidirectional streaming for efficient communication with the Massa node.

5. **Thread Safety**: The `OperationSender` methods are not thread-safe. Use separate instances or add your own synchronization if calling from multiple goroutines.
