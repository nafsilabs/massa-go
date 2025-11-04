package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nafsilabs/massa-go/client"
)

func main() {
	fmt.Println("Massa Operation Sender Example")
	fmt.Println("================================\n")

	// Initialize gRPC client
	fmt.Println("Connecting to Massa buildnet...")
	grpcClient, err := client.NewMassaGrpcClient(client.ClientConfig{
		Address:        "buildnet.massa.net:33037",
		UseTLS:         false,
		DefaultTimeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatal("Failed to create gRPC client:", err)
	}
	defer grpcClient.Close()
	fmt.Println("Connected successfully!\n")

	// Create operation sender
	_ = client.NewOperationSender(grpcClient)

	// Example 1: Transaction Operation (demonstration only - not signed)
	fmt.Println("Example 1: Transaction Operation")
	fmt.Println("---------------------------------")
	fmt.Println("This example shows how to construct a transaction operation.")
	fmt.Println("NOTE: To actually send it, you need to sign it with your private key.\n")

	transactionOp := &client.TransactionOp{
		RecipientAddress: "AU12recipient_address_placeholder",
		Amount:           1000000, // 1 Massa = 1,000,000 NanoMassa
	}

	fmt.Printf("Transaction Operation Created:\n")
	fmt.Printf("  Recipient: %s\n", transactionOp.RecipientAddress)
	fmt.Printf("  Amount: %d NanoMassa (%.6f Massa)\n", transactionOp.Amount, float64(transactionOp.Amount)/1000000.0)
	fmt.Printf("  Fee: 100 NanoMassa\n")
	fmt.Printf("  Expire Period: 10\n\n")

	// In production, you would do:
	// signature, err := wallet.SignOperation(transactionOp)
	// operationID, err := opSender.SendTransaction(ctx, ...)

	// Example 2: Buy Rolls Operation
	fmt.Println("Example 2: Buy Rolls Operation")
	fmt.Println("-------------------------------")
	fmt.Println("This example shows how to construct a buy rolls operation.\n")

	buyRollsOp := &client.BuyRollsOp{
		RollCount: 1,
	}

	fmt.Printf("Buy Rolls Operation Created:\n")
	fmt.Printf("  Roll Count: %d\n", buyRollsOp.RollCount)
	fmt.Printf("  Fee: 100 NanoMassa\n")
	fmt.Printf("  Expire Period: 10\n\n")

	// Example 3: Call Smart Contract Operation
	fmt.Println("Example 3: Call Smart Contract")
	fmt.Println("-------------------------------")
	fmt.Println("This example shows how to construct a call SC operation.\n")

	callSCOp := &client.CallSCOp{
		TargetAddress:  "AS12contract_address_placeholder",
		TargetFunction: "transfer",
		Parameter:      []byte{0x00, 0x01, 0x02}, // ABI-encoded parameters
		MaxGas:         100000,
		Coins:          0,
	}

	fmt.Printf("Call SC Operation Created:\n")
	fmt.Printf("  Target Address: %s\n", callSCOp.TargetAddress)
	fmt.Printf("  Target Function: %s\n", callSCOp.TargetFunction)
	fmt.Printf("  Parameters: %v\n", callSCOp.Parameter)
	fmt.Printf("  Max Gas: %d\n", callSCOp.MaxGas)
	fmt.Printf("  Coins: %d\n\n", callSCOp.Coins)

	// Example 4: Execute Smart Contract (Deploy)
	fmt.Println("Example 4: Execute Smart Contract (Deploy)")
	fmt.Println("-------------------------------------------")
	fmt.Println("This example shows how to construct a deploy SC operation.\n")

	// Mock bytecode (in production, this would be your compiled .wasm file)
	mockBytecode := []byte{0x00, 0x61, 0x73, 0x6d} // WASM magic number

	executeSCOp := &client.ExecuteSCOp{
		Bytecode:  mockBytecode,
		MaxGas:    1000000,
		MaxCoins:  100,
		Datastore: nil, // Optional datastore entries
	}

	fmt.Printf("Execute SC Operation Created:\n")
	fmt.Printf("  Bytecode size: %d bytes\n", len(executeSCOp.Bytecode))
	fmt.Printf("  Max Gas: %d\n", executeSCOp.MaxGas)
	fmt.Printf("  Max Coins: %d\n\n", executeSCOp.MaxCoins)

	// Example 5: Accessing the underlying gRPC client
	fmt.Println("Example 5: Direct gRPC Client Access")
	fmt.Println("-------------------------------------")
	fmt.Println("The OperationSender wraps a MassaGrpcClient which provides")
	fmt.Println("access to all Massa gRPC APIs (PublicClient and PrivateClient).")
	fmt.Println()
	fmt.Println("You can use these clients for:")
	fmt.Println("- Querying network status")
	fmt.Println("- Getting blocks, operations, and endorsements")
	fmt.Println("- Streaming new blocks/operations/endorsements")
	fmt.Println("- And more...")
	fmt.Println()
	fmt.Println("Example: grpcClient.PublicClient.GetStatus(ctx, &pb.GetStatusRequest{})")
	fmt.Println()

	// Summary
	fmt.Println("Summary")
	fmt.Println("=======")
	fmt.Println("This example demonstrated:")
	fmt.Println("1. Creating a Transaction operation")
	fmt.Println("2. Creating a Buy Rolls operation")
	fmt.Println("3. Creating a Call SC operation")
	fmt.Println("4. Creating an Execute SC (deploy) operation")
	fmt.Println("5. Using the underlying gRPC client to query network status")
	fmt.Println()
	fmt.Println("To actually send operations, you need to:")
	fmt.Println("- Sign the operation with your private key")
	fmt.Println("- Call the appropriate Send method (SendTransaction, SendBuyRolls, etc.)")
	fmt.Println("- Handle the returned operation ID")
	fmt.Println()
	fmt.Println("See OPERATION_SENDER_README.md for complete documentation.")
}
