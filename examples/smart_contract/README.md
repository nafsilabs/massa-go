# Massa Go Smart Contract Example

This example demonstrates how to create a smart contract for the Massa blockchain using pure Go and the Massa Go SDK.

## Features

This smart contract showcases:

- **Event Generation**: Emitting events to the blockchain
- **Storage Operations**: Storing and retrieving data
- **Context Information**: Accessing caller, callee, timestamp, etc.
- **Balance Operations**: Checking and transferring coins
- **Cryptographic Functions**: Hashing data with BLAKE3
- **Address Validation**: Validating Massa addresses
- **Structured Data**: JSON serialization/deserialization
- **Access Control**: Restricting functions to specific users
- **Error Handling**: Proper error logging and event emission

## Contract Functions

### Main Functions

- `main()` - Entry point, called automatically when contract is executed
- `GetContractInfo()` - Returns information about the contract
- `GetGreeting()` - Returns the current greeting message
- `UpdateGreeting(string)` - Updates the greeting (deployer only)
- `TransferFunds(string, uint64)` - Transfers funds to an address

### SDK Demonstration

The contract demonstrates various SDK features:

1. **Context Functions**:
   - `massa.Caller()` - Get the address that called the contract
   - `massa.Callee()` - Get the current contract address
   - `massa.Timestamp()` - Get current blockchain timestamp

2. **Storage Operations**:
   - `massa.SetString()` / `massa.GetString()` - Store/retrieve strings
   - `massa.SetJSON()` / `massa.GetJSON()` - Store/retrieve JSON data

3. **Balance & Transfers**:
   - `massa.Balance()` - Get current balance
   - `massa.TransferCoins()` - Transfer coins to another address
   - `massa.FormatBalance()` - Format balance for display

4. **Events & Logging**:
   - `massa.GenerateEvent()` - Generate blockchain events
   - `massa.Print()` / `massa.Printf()` - Debug logging
   - `massa.LogInfo()` / `massa.LogError()` - Structured logging
   - `massa.CreateEvent()` - Create structured events

5. **Cryptography**:
   - `massa.Blake3()` - BLAKE3 hashing
   - `massa.BytesToHex()` - Convert bytes to hex string

6. **Address Operations**:
   - `massa.ValidateAddress()` - Validate address format
   - `massa.NewAddress()` - Create address objects

## Building the Contract

### Prerequisites

- Go 1.21 or later
- Make sure `GOOS=js` and `GOARCH=wasm` are supported

### Build Commands

```bash
# Make the build script executable
chmod +x build.sh

# Build the contract
./build.sh
```

Or manually:

```bash
export GOOS=js
export GOARCH=wasm
go build -ldflags="-s -w" -o smart_contract.wasm main.go
```

### Optimization (Optional)

For smaller WASM files, install optimization tools:

```bash
# Install wasm-opt for size optimization
npm install -g wasm-opt

# Install wabt for validation
# On macOS: brew install wabt
# On Ubuntu: apt install wabt
```

## Deploying to Massa

1. **Build the Contract**: Run `./build.sh` to generate `smart_contract.wasm`

2. **Deploy with Massa Client**:
   ```bash
   # Using massa client
   deploy_contract smart_contract.wasm
   ```

3. **Call Contract Functions**:
   ```bash
   # Call the main function (automatically called on deployment)
   call_contract <contract_address> main

   # Get greeting
   call_contract <contract_address> GetGreeting

   # Update greeting (deployer only)
   call_contract <contract_address> UpdateGreeting "Hello from Massa!"

   # Transfer funds
   call_contract <contract_address> TransferFunds "AU12..." 1000000
   ```

## Contract Events

The contract emits several types of events:

- `Hello, World!` - Basic greeting event
- `contract_execution` - Details about contract execution
- `contract_info_stored` - When contract info is stored
- `greeting_updated` - When greeting is changed
- `funds_transferred` - When funds are transferred

## Error Handling

The contract includes comprehensive error handling:

- Invalid addresses are rejected
- Insufficient balance prevents transfers
- Access control prevents unauthorized operations
- All errors are logged with context information

## Testing

The SDK includes mock implementations for testing:

```go
// The SDK automatically detects non-WASM environments
// and provides mock implementations for testing
func isWasmRuntime() bool {
    // Returns false in test environments
    return false
}
```

## Development Notes

### Memory Management

The Go SDK handles WebAssembly memory management automatically:

- String and byte slice conversion to/from WASM pointers
- Automatic cleanup of temporary allocations
- Safe handling of null pointers

### Performance Considerations

- Use `massa.SetJSON()` for complex data structures
- Batch storage operations when possible
- Consider gas costs for large data operations
- Use events for off-chain data indexing

### Security Best Practices

1. **Validate Inputs**: Always validate addresses and amounts
2. **Access Control**: Implement proper permission checks
3. **Error Handling**: Don't expose internal errors to callers
4. **State Consistency**: Use transactions for multi-step operations

## SDK Documentation

For complete SDK documentation, see the main package files:

- `address.go` - Address operations and validation
- `storage.go` - Persistent storage operations
- `context.go` - Execution context information
- `contract.go` - Contract calls and bytecode operations
- `coins.go` - Balance and transfer operations
- `events.go` - Event generation and logging
- `crypto.go` - Cryptographic functions
- `op_datastore.go` - Operation-level data access

## License

This example is part of the Massa Go SDK and follows the same license terms.