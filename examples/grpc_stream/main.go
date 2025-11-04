package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/nafsilabs/massa-go/client"
// 	pb "github.com/nafsilabs/massa-go/client/proto/massa/api/v1"
// )

// func main() {
// 	// Create a new Massa gRPC client
// 	config := client.ClientConfig{
// 		Address:        "buildnet.massa.net:33037",
// 		UseTLS:         false,
// 		DefaultTimeout: 30 * time.Second,
// 	}

// 	grpcClient, err := client.NewMassaGrpcClient(config)
// 	if err != nil {
// 		log.Fatalf("Failed to create gRPC client: %v", err)
// 	}
// 	defer grpcClient.Close()

// 	// Example 1: Get Status (Unary call)
// 	fmt.Println("=== Getting Status ===")
// 	ctx, cancel := grpcClient.GetContext()
// 	defer cancel()

// 	statusResp, err := grpcClient.PublicClient.GetStatus(ctx, &pb.GetStatusRequest{})
// 	if err != nil {
// 		log.Fatalf("GetStatus failed: %v", err)
// 	}
// 	fmt.Printf("Status: %+v\n\n", statusResp)

// 	// Example 2: Server Streaming - NewOperationsServer
// 	fmt.Println("=== Streaming New Operations (Server Stream) ===")
// 	streamCtx, streamCancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer streamCancel()

// 	operationCount := 0
// 	err = grpcClient.NewOperationsServerStream(streamCtx, &pb.NewOperationsServerRequest{
// 		Filters: []*pb.NewOperationsFilter{},
// 	}, func(msg *pb.NewOperationsServerResponse) error {
// 		operationCount++
// 		fmt.Printf("Received operation: %+v\n", msg)
// 		return nil
// 	})

// 	if err != nil && streamCtx.Err() != context.DeadlineExceeded {
// 		log.Printf("NewOperationsServerStream error: %v", err)
// 	}
// 	fmt.Printf("Total operations received: %d\n\n", operationCount)

// 	// Example 3: Server Streaming - NewBlocksServer
// 	fmt.Println("=== Streaming New Blocks (Server Stream) ===")
// 	blockCtx, blockCancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer blockCancel()

// 	blockCount := 0
// 	err = grpcClient.NewBlocksServerStream(blockCtx, &pb.NewBlocksServerRequest{
// 		Filters: []*pb.NewBlocksFilter{},
// 	}, func(msg *pb.NewBlocksServerResponse) error {
// 		blockCount++
// 		if msg.SignedBlock != nil {
// 			fmt.Printf("Received block: %+v\n", msg.SignedBlock)
// 		}
// 		return nil
// 	})

// 	if err != nil && blockCtx.Err() != context.DeadlineExceeded {
// 		log.Printf("NewBlocksServerStream error: %v", err)
// 	}
// 	fmt.Printf("Total blocks received: %d\n\n", blockCount)

// 	// Example 4: Bidirectional Streaming - NewOperations
// 	fmt.Println("=== Bidirectional Stream Example ===")
// 	bidiCtx, bidiCancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer bidiCancel()

// 	stream, err := grpcClient.NewOperationsStream(bidiCtx)
// 	if err != nil {
// 		log.Printf("NewOperationsStream failed: %v", err)
// 		return
// 	}

// 	// Send a subscription request
// 	if err := stream.Send(&pb.NewOperationsRequest{
// 		Filters: []*pb.NewOperationsFilter{},
// 	}); err != nil {
// 		log.Printf("Failed to send request: %v", err)
// 		return
// 	}

// 	// Receive responses
// 	go func() {
// 		for {
// 			resp, err := stream.Recv()
// 			if err != nil {
// 				log.Printf("Receive error: %v", err)
// 				return
// 			}
// 			fmt.Printf("Bidirectional stream received: %+v\n", resp)
// 		}
// 	}()

// 	// Wait for context timeout
// 	<-bidiCtx.Done()
// 	stream.CloseSend()

// 	fmt.Println("Done!")
// }
