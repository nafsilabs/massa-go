# Operation Sender Example

This example demonstrates how to use the `OperationSender` to construct different types of Massa blockchain operations.

## What This Example Shows

1. **Transaction Operation** - Transfer coins between addresses
2. **Buy Rolls Operation** - Purchase rolls for staking
3. **Call Smart Contract** - Call a function on a deployed smart contract
4. **Execute Smart Contract** - Deploy a new smart contract

## Running the Example

```bash
cd examples/operation_sender
go build
./operation_sender
```

## Expected Output

The example will:
- Connect to the Massa buildnet
- Create sample operations (without signing)
- Display the operation details
- Show how to access the underlying gRPC client

## Important Notes

### This is a Demonstration Only

This example does NOT actually send operations to the network because:
- Operations must be signed with a private key
- The example uses placeholder addresses and data

### To Actually Send Operations

To send real operations, you need to:

1. **Sign the operation** with your private key:
```go
// You'll need a wallet implementation
signature, err := wallet.SignOperation(operation)
publicKey := wallet.GetPublicKey()
```

2. **Call the Send method**:
```go
opSender := client.NewOperationSender(grpcClient)

operationID, err := opSender.SendTransaction(
    ctx,
    recipientAddress,
    amount,
    fee,
    expirePeriod,
    signature,  // From step 1
    publicKey,  // From step 1
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Operation sent! ID: %s\n", operationID)
```

3. **Handle the response**:
```go
// The operation ID can be used to track the operation status
fmt.Printf("Operation ID: %s\n", operationID)
```

## Operation Types

### TransactionOp
Transfers coins from sender to recipient.

```go
op := &client.TransactionOp{
    RecipientAddress: "AU12...",
    Amount:           1000000, // NanoMassa
}
```

### BuyRollsOp
Purchases rolls for staking.

```go
op := &client.BuyRollsOp{
    RollCount: 1,
}
```

### CallSCOp
Calls a function on a smart contract.

```go
op := &client.CallSCOp{
    TargetAddress:  "AS12...",
    TargetFunction: "transfer",
    Parameter:      []byte{/* ABI encoded */},
    MaxGas:         100000,
    Coins:          0,
}
```

### ExecuteSCOp
Deploys a new smart contract.

```go
bytecode, err := os.ReadFile("contract.wasm")
op := &client.ExecuteSCOp{
    Bytecode:  bytecode,
    MaxGas:    1000000,
    MaxCoins:  100,
    Datastore: nil,
}
```

## Next Steps

1. See [OPERATION_SENDER_README.md](../../client/OPERATION_SENDER_README.md) for complete documentation
2. Implement wallet signing functionality
3. Test with small amounts on buildnet first
4. Monitor operation status after sending

## Related Examples

- **grpc_stream** - Demonstrates gRPC streaming with new blocks/operations
- **smart_contract** - Shows smart contract deployment and interaction
- **send_coins** - Basic coin transfer example

## Resources

- [Massa Documentation](https://docs.massa.net)
- [Massa gRPC API](https://github.com/massalabs/massa-proto)
- [Operation Sender README](../../client/OPERATION_SENDER_README.md)
