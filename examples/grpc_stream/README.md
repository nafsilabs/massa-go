# Massa gRPC Streaming Example

This example demonstrates how to use the Massa gRPC client with streaming support.

## What This Example Shows

1. **Unary Calls**: Simple request-response (GetStatus)
2. **Server Streaming**: Subscribe to new operations and blocks
3. **Bidirectional Streaming**: Two-way communication with the server

## Running the Example

```bash
cd examples/grpc_stream
go run .
```

Or build and run:

```bash
go build .
./grpc_stream
```

## What Happens

The example will:

1. Connect to the Massa buildnet gRPC server
2. Get the current node status
3. Stream new operations for 10 seconds
4. Stream new blocks for 10 seconds  
5. Open a bidirectional stream for operations for 10 seconds

## Code Structure

- **Unary call example**: `GetStatus` - fetches node status
- **Server streaming example**: `NewOperationsServerStream` and `NewBlocksServerStream` - receives continuous updates
- **Bidirectional streaming example**: `NewOperationsStream` - allows sending and receiving messages

## Customization

You can modify:
- `config.Address` - Change to mainnet or another network
- Stream timeouts - Adjust the context timeout durations
- Filters - Add filters to the stream requests to only receive specific data

## Output

You'll see:
- Node status information
- Count of operations received
- Count of blocks received
- Real-time stream messages (if any data is available)

## Notes

- The streams timeout after the specified duration (10 seconds in the example)
- If no blocks or operations are produced during the timeout, you'll see zero counts
- Connection errors will be displayed if the server is unavailable
