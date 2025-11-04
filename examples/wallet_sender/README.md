# Wallet-Integrated Operation Sender Example

This example demonstrates how to use the `WalletOperationSender` to send operations with automatic wallet signing.

## What's New?

The `WalletOperationSender` provides **automatic signing** using wallet accounts, eliminating the need to manually handle signatures and public keys.

## Running the Example

```bash
cd examples/wallet_sender
go build
./wallet_sender
```

## Key Features

### Automatic Signing
No need to manually sign operations - the sender handles it automatically:

```go
// Create wallet operation sender
walletSender := client.NewWalletOperationSender(grpcClient, accountManager)

// Send transaction with automatic signing
operationID, err := walletSender.SendTransactionWithWallet(
    ctx,
    "my-account",            // Account nickname
    "my-password",           // Account password
    "AU12recipient",         // Recipient address
    1000000,                 // Amount
    100,                     // Fee
    currentPeriod + 3,       // Expiry
)
```

### vs Manual Signing

**Old Way (OperationSender):**
```go
// Step 1: Create operation
op := &client.TransactionOp{...}

// Step 2: Sign manually (external process)
signature, publicKey := wallet.SignOperation(op)

// Step 3: Send with signature
operationID, err := sender.SendOperation(ctx, op, fee, expiry, signature, publicKey)
```

**New Way (WalletOperationSender):**
```go
// One step - automatic signing!
operationID, err := walletSender.SendTransactionWithWallet(
    ctx, "my-account", "password", recipientAddr, amount, fee, expiry,
)
```

## Usage Examples

### 1. Setup Wallet and Client

```go
// Create account manager
accountManager := wallet.NewAccountManager()

// Create or import an account
account, err := accountManager.CreateAccount("my-account", "my-password")
// OR
account, err := accountManager.ImportAccountFromPrivateKey("my-account", privateKey, "password")

// Connect to Massa network
grpcClient, err := client.NewMassaGrpcClient(client.ClientConfig{
    Address: "buildnet.massa.net:33037",
    UseTLS:  false,
})
defer grpcClient.Close()

// Create wallet operation sender
walletSender := client.NewWalletOperationSender(grpcClient, accountManager)
```

### 2. Send Transaction

```go
operationID, err := walletSender.SendTransactionWithWallet(
    ctx,
    "my-account",              // Your account nickname
    "my-password",             // Your account password
    "AU12recipient_address",   // Recipient
    1000000,                   // 1 Massa = 1,000,000 NanoMassa
    100,                       // Fee
    getCurrentPeriod() + 3,    // Expiry period
)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Transaction sent! ID: %s\n", operationID)
```

### 3. Buy Rolls

```go
operationID, err := walletSender.SendBuyRollsWithWallet(
    ctx,
    "my-account",
    "my-password",
    1,                         // Number of rolls
    100,                       // Fee
    getCurrentPeriod() + 3,
)
```

### 4. Call Smart Contract

```go
// Prepare parameters (ABI encoded)
parameters := []byte{/* your encoded params */}

operationID, err := walletSender.SendCallSCWithWallet(
    ctx,
    "my-account",
    "my-password",
    "AS12contract_address",    // Smart contract address
    "transfer",                // Function name
    parameters,                // Function parameters
    100000,                    // Max gas
    0,                         // Coins to send
    100,                       // Fee
    getCurrentPeriod() + 3,
)
```

### 5. Deploy Smart Contract

```go
bytecode, err := os.ReadFile("contract.wasm")
if err != nil {
    log.Fatal(err)
}

operationID, err := walletSender.SendExecuteSCWithWallet(
    ctx,
    "my-account",
    "my-password",
    bytecode,                  // Contract bytecode
    1000000,                   // Max gas
    100,                       // Max coins
    nil,                       // Datastore (optional)
    100,                       // Fee
    getCurrentPeriod() + 3,
)
```

## Available Methods

### Core Methods

- `SendOperationWithWallet(ctx, nickname, password, op, fee, expiry)` - Send any operation type
- `SendTransactionWithWallet(...)` - Send coin transfer
- `SendBuyRollsWithWallet(...)` - Buy rolls for staking
- `SendSellRollsWithWallet(...)` - Sell rolls
- `SendExecuteSCWithWallet(...)` - Deploy smart contract
- `SendCallSCWithWallet(...)` - Call smart contract function

## Security Considerations

### Password Management
- Passwords are used to decrypt private keys temporarily for signing
- Private keys are never stored unencrypted in memory longer than necessary
- Use strong passwords for accounts

### Best Practices
1. **Never hardcode passwords** - read from environment or secure input
2. **Use secure password storage** - consider using OS keychain
3. **Limit account access** - use separate accounts for different purposes
4. **Backup private keys** - export and store securely offline

## Integration with Existing Wallet

If you already have a wallet file:

```go
// Load existing wallet
wallet, err := wallet.LoadWallet("path/to/wallet.json")
if err != nil {
    log.Fatal(err)
}

// Use the wallet's account manager
walletSender := client.NewWalletOperationSender(grpcClient, wallet.AccountManager)

// Send operations
operationID, err := walletSender.SendTransactionWithWallet(...)
```

## Error Handling

```go
operationID, err := walletSender.SendTransactionWithWallet(...)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "failed to unlock account"):
        fmt.Println("Incorrect password or account not found")
    case strings.Contains(err.Error(), "failed to sign operation"):
        fmt.Println("Signing error - check account status")
    case strings.Contains(err.Error(), "operation error"):
        fmt.Println("Network rejected operation:", err)
    case strings.Contains(err.Error(), "failed to open stream"):
        fmt.Println("Network connection error:", err)
    default:
        fmt.Println("Unexpected error:", err)
    }
    return
}

fmt.Printf("Success! Operation ID: %s\n", operationID)
```

## Performance Tips

1. **Reuse WalletOperationSender** - create once, use multiple times
2. **Keep AccountManager** - don't recreate for each operation
3. **Connection pooling** - reuse gRPC client connection
4. **Batch operations** - if sending multiple operations, consider batching

## Comparison Table

| Feature | OperationSender | WalletOperationSender |
|---------|-----------------|----------------------|
| Signing | Manual (external) | Automatic (built-in) |
| Wallet Integration | No | Yes |
| Password Protection | No | Yes |
| Key Management | External | Built-in |
| Convenience Methods | Yes | Yes |
| All Operation Types | Yes | Yes |
| Learning Curve | Higher | Lower |

## When to Use Each

### Use `OperationSender` when:
- You have your own signing mechanism
- Using hardware wallets or external signers
- Need fine-grained control over signing
- Integrating with existing wallet systems

### Use `WalletOperationSender` when:
- Building a wallet application
- Want automatic signing
- Need built-in key management
- Prefer simplicity over customization

## Next Steps

1. See [massa_client.go](../../client/massa_client.go) for implementation details
2. Check [OPERATION_SENDER_README.md](../../client/OPERATION_SENDER_README.md) for manual signing
3. Explore [wallet package](../../wallet/) for advanced wallet features
4. Test on buildnet before using on mainnet

## Related Examples

- **operation_sender** - Manual signing approach
- **grpc_stream** - Streaming operations and blocks
- **send_coins** - Basic transaction example

## Resources

- [Massa Documentation](https://docs.massa.net)
- [Massa Wallet Specification](https://github.com/massalabs/massa-wallet)
- [Operation Signing Guide](https://docs.massa.net/en/latest/web3-dev/smart-contracts.html)
