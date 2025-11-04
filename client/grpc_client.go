// other gRPC client methods
package client

import (
	"context"
	"fmt"
	"io"
	"time"

	pb "github.com/nafsilabs/massa-go/client/proto/massa/api/v1"

	"google.golang.org/grpc"
)

// GetContext returns a context with the default timeout
func (c *MassaClient) GetContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.defaultTimeout)
}

// GetContextWithTimeout returns a context with a custom timeout
func (c *MassaClient) GetContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// StreamHandler is a generic handler for processing stream responses
type StreamHandler[T any] func(msg *T) error

// NewOperationsStream creates a bidirectional stream for new operations
func (c *MassaClient) NewOperationsStream(ctx context.Context) (grpc.BidiStreamingClient[pb.NewOperationsRequest, pb.NewOperationsResponse], error) {
	return c.PublicClient.NewOperations(ctx)
}

// NewBlocksStream creates a bidirectional stream for new blocks
func (c *MassaClient) NewBlocksStream(ctx context.Context) (grpc.BidiStreamingClient[pb.NewBlocksRequest, pb.NewBlocksResponse], error) {
	return c.PublicClient.NewBlocks(ctx)
}

// NewFilledBlocksStream creates a bidirectional stream for new filled blocks
func (c *MassaClient) NewFilledBlocksStream(ctx context.Context) (grpc.BidiStreamingClient[pb.NewFilledBlocksRequest, pb.NewFilledBlocksResponse], error) {
	return c.PublicClient.NewFilledBlocks(ctx)
}

// NewEndorsementsStream creates a bidirectional stream for new endorsements
func (c *MassaClient) NewEndorsementsStream(ctx context.Context) (grpc.BidiStreamingClient[pb.NewEndorsementsRequest, pb.NewEndorsementsResponse], error) {
	return c.PublicClient.NewEndorsements(ctx)
}

// NewSlotExecutionOutputsStream creates a bidirectional stream for slot execution outputs
func (c *MassaClient) NewSlotExecutionOutputsStream(ctx context.Context) (grpc.BidiStreamingClient[pb.NewSlotExecutionOutputsRequest, pb.NewSlotExecutionOutputsResponse], error) {
	return c.PublicClient.NewSlotExecutionOutputs(ctx)
}

// NewSlotABICallStacksStream creates a bidirectional stream for slot ABI call stacks
func (c *MassaClient) NewSlotABICallStacksStream(ctx context.Context) (grpc.BidiStreamingClient[pb.NewSlotABICallStacksRequest, pb.NewSlotABICallStacksResponse], error) {
	return c.PublicClient.NewSlotABICallStacks(ctx)
}

// NewSlotTransfersStream creates a bidirectional stream for slot transfers
func (c *MassaClient) NewSlotTransfersStream(ctx context.Context) (grpc.BidiStreamingClient[pb.NewSlotTransfersRequest, pb.NewSlotTransfersResponse], error) {
	return c.PublicClient.NewSlotTransfers(ctx)
}

// SendOperationsStream creates a bidirectional stream for sending operations
func (c *MassaClient) SendOperationsStream(ctx context.Context) (grpc.BidiStreamingClient[pb.SendOperationsRequest, pb.SendOperationsResponse], error) {
	return c.PublicClient.SendOperations(ctx)
}

// SendBlocksStream creates a bidirectional stream for sending blocks
func (c *MassaClient) SendBlocksStream(ctx context.Context) (grpc.BidiStreamingClient[pb.SendBlocksRequest, pb.SendBlocksResponse], error) {
	return c.PublicClient.SendBlocks(ctx)
}

// SendEndorsementsStream creates a bidirectional stream for sending endorsements
func (c *MassaClient) SendEndorsementsStream(ctx context.Context) (grpc.BidiStreamingClient[pb.SendEndorsementsRequest, pb.SendEndorsementsResponse], error) {
	return c.PublicClient.SendEndorsements(ctx)
}

// TransactionsThroughputStream creates a bidirectional stream for transactions throughput
func (c *MassaClient) TransactionsThroughputStream(ctx context.Context) (grpc.BidiStreamingClient[pb.TransactionsThroughputRequest, pb.TransactionsThroughputResponse], error) {
	return c.PublicClient.TransactionsThroughput(ctx)
}

// Server streaming methods

// NewOperationsServerStream creates a server streaming connection for new operations
func (c *MassaClient) NewOperationsServerStream(ctx context.Context, req *pb.NewOperationsServerRequest, handler StreamHandler[pb.NewOperationsServerResponse]) error {
	stream, err := c.PublicClient.NewOperationsServer(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create new operations server stream: %w", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error receiving from stream: %w", err)
		}
		if err := handler(msg); err != nil {
			return fmt.Errorf("error handling message: %w", err)
		}
	}
}

// NewBlocksServerStream creates a server streaming connection for new blocks
func (c *MassaClient) NewBlocksServerStream(ctx context.Context, req *pb.NewBlocksServerRequest, handler StreamHandler[pb.NewBlocksServerResponse]) error {
	stream, err := c.PublicClient.NewBlocksServer(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create new blocks server stream: %w", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error receiving from stream: %w", err)
		}
		if err := handler(msg); err != nil {
			return fmt.Errorf("error handling message: %w", err)
		}
	}
}

// NewFilledBlocksServerStream creates a server streaming connection for new filled blocks
func (c *MassaClient) NewFilledBlocksServerStream(ctx context.Context, req *pb.NewFilledBlocksServerRequest, handler StreamHandler[pb.NewFilledBlocksServerResponse]) error {
	stream, err := c.PublicClient.NewFilledBlocksServer(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create new filled blocks server stream: %w", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error receiving from stream: %w", err)
		}
		if err := handler(msg); err != nil {
			return fmt.Errorf("error handling message: %w", err)
		}
	}
}

// NewEndorsementsServerStream creates a server streaming connection for new endorsements
func (c *MassaClient) NewEndorsementsServerStream(ctx context.Context, req *pb.NewEndorsementsServerRequest, handler StreamHandler[pb.NewEndorsementsServerResponse]) error {
	stream, err := c.PublicClient.NewEndorsementsServer(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create new endorsements server stream: %w", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error receiving from stream: %w", err)
		}
		if err := handler(msg); err != nil {
			return fmt.Errorf("error handling message: %w", err)
		}
	}
}

// NewSlotExecutionOutputsServerStream creates a server streaming connection for slot execution outputs
func (c *MassaClient) NewSlotExecutionOutputsServerStream(ctx context.Context, req *pb.NewSlotExecutionOutputsServerRequest, handler StreamHandler[pb.NewSlotExecutionOutputsServerResponse]) error {
	stream, err := c.PublicClient.NewSlotExecutionOutputsServer(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create slot execution outputs server stream: %w", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error receiving from stream: %w", err)
		}
		if err := handler(msg); err != nil {
			return fmt.Errorf("error handling message: %w", err)
		}
	}
}

// TransactionsThroughputServerStream creates a server streaming connection for transactions throughput
func (c *MassaClient) TransactionsThroughputServerStream(ctx context.Context, req *pb.TransactionsThroughputServerRequest, handler StreamHandler[pb.TransactionsThroughputServerResponse]) error {
	stream, err := c.PublicClient.TransactionsThroughputServer(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create transactions throughput server stream: %w", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error receiving from stream: %w", err)
		}
		if err := handler(msg); err != nil {
			return fmt.Errorf("error handling message: %w", err)
		}
	}
}

// NewTransfersInfoServerStream creates a server streaming connection for transfers info
func (c *MassaClient) NewTransfersInfoServerStream(ctx context.Context, req *pb.NewTransfersInfoServerRequest, handler StreamHandler[pb.NewTransfersInfoServerResponse]) error {
	stream, err := c.PublicClient.NewTransfersInfoServer(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create transfers info server stream: %w", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error receiving from stream: %w", err)
		}
		if err := handler(msg); err != nil {
			return fmt.Errorf("error handling message: %w", err)
		}
	}
}
