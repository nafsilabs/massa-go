package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nafsilabs/massa-go/client"
	"github.com/nafsilabs/massa-go/examples"
	"github.com/nafsilabs/massa-go/utils"
	"github.com/nafsilabs/massa-go/wallet"
)

func main() {

	// Create a lightweight wallet. If you already have an account you can
	// LoadWallet and use that instead. Here we import the example private key
	// defined in examples/constants.go so the example is repeatable.
	w := wallet.NewWallet(&wallet.WalletConfig{WalletPath: examples.WalletPath})

	// Try to reuse existing account, otherwise import from the example private key.
	acc, err := w.AccountManager.GetAccount(examples.Nickname)
	if err != nil {
		acc, err = w.ImportAccountFromPrivateKey(examples.Nickname, examples.Secret, examples.Password)
		if err != nil {
			fmt.Printf("failed to import account from secret: %v\n", err)
			os.Exit(1)
		}
	}

	// Initialize gRPC client
	fmt.Println("Connecting to Massa buildnet...")
	cfg := client.ClientConfig{
		Address:        "buildnet.massa.net:33037",
		UseTLS:         false,
		DefaultTimeout: 30 * time.Second,
		ChainID:        utils.BUILDNET,
		Account:        acc,
	}
	if err != nil {
		log.Fatal("Failed to create gRPC client:", err)
	}
	// Create operation sender
	grpc, err := client.NewMassaClient(&cfg)
	if err != nil {
		log.Fatal("Failed to create massa client:", err)
	}
	defer grpc.Close()
	fmt.Println("Connected successfully!")
	fmt.Println("Buying roles...")

	opID, err := grpc.BuyRolls(context.Background(), examples.Nickname, examples.Password, 2, 0.01)
	if err != nil {
		fmt.Printf("failed to buy rolls: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Buy Rolls Operation Created: %v\n", opID)
	fmt.Printf("  Roll Count: %d\n", 2)
	fmt.Printf("  Fee: 100 NanoMassa\n")
	fmt.Printf("  Expire Period: 10\n\n")

}
