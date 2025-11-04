package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"time"

	pb "github.com/nafsilabs/massa-go/client/proto/massa/api/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// MassaGrpcClient wraps the gRPC connection and provides access to both public and private services
type MassaGrpcClient struct {
	conn           *grpc.ClientConn
	PublicClient   pb.PublicServiceClient
	PrivateClient  pb.PrivateServiceClient
	defaultTimeout time.Duration
}

// ClientConfig holds configuration for the gRPC client
type ClientConfig struct {
	Address        string
	UseTLS         bool
	TLSConfig      *tls.Config
	DefaultTimeout time.Duration
	DialOptions    []grpc.DialOption
}

// NewMassaGrpcClient creates a new Massa gRPC client
func NewMassaGrpcClient(config ClientConfig) (*MassaGrpcClient, error) {
	if config.DefaultTimeout == 0 {
		config.DefaultTimeout = 30 * time.Second
	}

	var opts []grpc.DialOption

	// Add transport credentials
	if config.UseTLS {
		if config.TLSConfig != nil {
			opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(config.TLSConfig)))
		} else {
			opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
		}
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Add any additional dial options
	opts = append(opts, config.DialOptions...)

	// Establish connection
	conn, err := grpc.Dial(config.Address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	return &MassaGrpcClient{
		conn:           conn,
		PublicClient:   pb.NewPublicServiceClient(conn),
		PrivateClient:  pb.NewPrivateServiceClient(conn),
		defaultTimeout: config.DefaultTimeout,
	}, nil
}

// Close closes the gRPC connection
func (c *MassaGrpcClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetContext returns a context with the default timeout
func (c *MassaGrpcClient) GetContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.defaultTimeout)
}

// GetContextWithTimeout returns a context with a custom timeout
func (c *MassaGrpcClient) GetContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// StreamHandler is a generic handler for processing stream responses
type StreamHandler[T any] func(msg *T) error

// NewOperationsStream creates a bidirectional stream for new operations
func (c *MassaGrpcClient) NewOperationsStream(ctx context.Context) (grpc.BidiStreamingClient[pb.NewOperationsRequest, pb.NewOperationsResponse], error) {
	return c.PublicClient.NewOperations(ctx)
}

// NewBlocksStream creates a bidirectional stream for new blocks
func (c *MassaGrpcClient) NewBlocksStream(ctx context.Context) (grpc.BidiStreamingClient[pb.NewBlocksRequest, pb.NewBlocksResponse], error) {
	return c.PublicClient.NewBlocks(ctx)
}

// NewFilledBlocksStream creates a bidirectional stream for new filled blocks
func (c *MassaGrpcClient) NewFilledBlocksStream(ctx context.Context) (grpc.BidiStreamingClient[pb.NewFilledBlocksRequest, pb.NewFilledBlocksResponse], error) {
	return c.PublicClient.NewFilledBlocks(ctx)
}

// NewEndorsementsStream creates a bidirectional stream for new endorsements
func (c *MassaGrpcClient) NewEndorsementsStream(ctx context.Context) (grpc.BidiStreamingClient[pb.NewEndorsementsRequest, pb.NewEndorsementsResponse], error) {
	return c.PublicClient.NewEndorsements(ctx)
}

// NewSlotExecutionOutputsStream creates a bidirectional stream for slot execution outputs
func (c *MassaGrpcClient) NewSlotExecutionOutputsStream(ctx context.Context) (grpc.BidiStreamingClient[pb.NewSlotExecutionOutputsRequest, pb.NewSlotExecutionOutputsResponse], error) {
	return c.PublicClient.NewSlotExecutionOutputs(ctx)
}

// NewSlotABICallStacksStream creates a bidirectional stream for slot ABI call stacks
func (c *MassaGrpcClient) NewSlotABICallStacksStream(ctx context.Context) (grpc.BidiStreamingClient[pb.NewSlotABICallStacksRequest, pb.NewSlotABICallStacksResponse], error) {
	return c.PublicClient.NewSlotABICallStacks(ctx)
}

// NewSlotTransfersStream creates a bidirectional stream for slot transfers
func (c *MassaGrpcClient) NewSlotTransfersStream(ctx context.Context) (grpc.BidiStreamingClient[pb.NewSlotTransfersRequest, pb.NewSlotTransfersResponse], error) {
	return c.PublicClient.NewSlotTransfers(ctx)
}

// SendOperationsStream creates a bidirectional stream for sending operations
func (c *MassaGrpcClient) SendOperationsStream(ctx context.Context) (grpc.BidiStreamingClient[pb.SendOperationsRequest, pb.SendOperationsResponse], error) {
	return c.PublicClient.SendOperations(ctx)
}

// SendBlocksStream creates a bidirectional stream for sending blocks
func (c *MassaGrpcClient) SendBlocksStream(ctx context.Context) (grpc.BidiStreamingClient[pb.SendBlocksRequest, pb.SendBlocksResponse], error) {
	return c.PublicClient.SendBlocks(ctx)
}

// SendEndorsementsStream creates a bidirectional stream for sending endorsements
func (c *MassaGrpcClient) SendEndorsementsStream(ctx context.Context) (grpc.BidiStreamingClient[pb.SendEndorsementsRequest, pb.SendEndorsementsResponse], error) {
	return c.PublicClient.SendEndorsements(ctx)
}

// TransactionsThroughputStream creates a bidirectional stream for transactions throughput
func (c *MassaGrpcClient) TransactionsThroughputStream(ctx context.Context) (grpc.BidiStreamingClient[pb.TransactionsThroughputRequest, pb.TransactionsThroughputResponse], error) {
	return c.PublicClient.TransactionsThroughput(ctx)
}

// Server streaming methods

// NewOperationsServerStream creates a server streaming connection for new operations
func (c *MassaGrpcClient) NewOperationsServerStream(ctx context.Context, req *pb.NewOperationsServerRequest, handler StreamHandler[pb.NewOperationsServerResponse]) error {
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
func (c *MassaGrpcClient) NewBlocksServerStream(ctx context.Context, req *pb.NewBlocksServerRequest, handler StreamHandler[pb.NewBlocksServerResponse]) error {
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
func (c *MassaGrpcClient) NewFilledBlocksServerStream(ctx context.Context, req *pb.NewFilledBlocksServerRequest, handler StreamHandler[pb.NewFilledBlocksServerResponse]) error {
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
func (c *MassaGrpcClient) NewEndorsementsServerStream(ctx context.Context, req *pb.NewEndorsementsServerRequest, handler StreamHandler[pb.NewEndorsementsServerResponse]) error {
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
func (c *MassaGrpcClient) NewSlotExecutionOutputsServerStream(ctx context.Context, req *pb.NewSlotExecutionOutputsServerRequest, handler StreamHandler[pb.NewSlotExecutionOutputsServerResponse]) error {
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
func (c *MassaGrpcClient) TransactionsThroughputServerStream(ctx context.Context, req *pb.TransactionsThroughputServerRequest, handler StreamHandler[pb.TransactionsThroughputServerResponse]) error {
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
func (c *MassaGrpcClient) NewTransfersInfoServerStream(ctx context.Context, req *pb.NewTransfersInfoServerRequest, handler StreamHandler[pb.NewTransfersInfoServerResponse]) error {
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
