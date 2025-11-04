package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nafsilabs/massa-go/client"
	examples "github.com/nafsilabs/massa-go/examples"
	srl "github.com/nafsilabs/massa-go/sc/serialisation"
	"github.com/nafsilabs/massa-go/utils"
	"github.com/nafsilabs/massa-go/wallet"
)

func main() {
	// Configuration — adapt these to your environment
	name := "alice"
	age := uint32(21)

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

	// Prefer the known calling address from constants for clarity, but fall back
	// to the account's derived address if absent.
	callerAddr := examples.CallingAddress
	if callerAddr == "" && acc.Address != nil {
		callerAddr = acc.Address.String()
	}

	cfg := client.ClientConfig{
		Address:        "buildnet.massa.net:33037",
		UseTLS:         false,
		DefaultTimeout: 30 * time.Second,
		ChainID:        utils.BUILDNET,
		Account:        acc,
	}
	grpc, err := client.NewMassaClient(&cfg)
	if err != nil {
		log.Fatal("Failed to create gRPC client:", err)
	}
	defer grpc.Close()
	fmt.Println("Connected successfully!")

	// Build parameters using Args: AddString(name) followed by AddU32(age).
	args := srl.NewArgs(nil)
	args.AddString(name)
	args.AddU32(age)
	param := args.Serialise()

	// Call the smart contract function 'changeAge'. We'll use the client helper
	// `CallFunction`, which signs the operation using the provided account.
	// fee, maxGas, coins, expiryDelta, async
	fee := 0.01
	maxGas := 1.0 // 0 means the client will estimate
	coins := 2.0  // no coins attached

	resp, err := grpc.CallSC(context.Background(), examples.Nickname, examples.Password, examples.SmartContractAddress, "changeAge", param, maxGas, coins, fee)
	if err != nil {
		log.Fatalf("Failed to call smart contract: %v\n", err)
	}

	// Wait for operation to be executed

	fmt.Printf("Operation submitted. Operation ID: %s\n", resp)
}
