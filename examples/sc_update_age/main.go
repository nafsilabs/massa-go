package main

import (
	"fmt"
	"os"

	"github.com/nafsilabs/massa-go/client"
	examples "github.com/nafsilabs/massa-go/examples"
	srl "github.com/nafsilabs/massa-go/sc/serialisation"
	"github.com/nafsilabs/massa-go/wallet"
)

func main() {
	// Configuration — adapt these to your environment
	name := "alice"
	age := uint32(21)
	walletPath := "example_wallet.json"
	nickname := "example-account"
	password := "example-password"

	// Create client (RPC params come from defaults inside client.NewClient)
	c := client.NewClient(false)

	// Create a lightweight wallet. If you already have an account you can
	// LoadWallet and use that instead. Here we import the example private key
	// defined in examples/constants.go so the example is repeatable.
	w := wallet.NewWallet(&wallet.WalletConfig{WalletPath: walletPath})

	// Try to reuse existing account, otherwise import from the example private key.
	acc, err := w.AccountManager.GetAccount(nickname)
	if err != nil {
		acc, err = w.ImportAccountFromPrivateKey(nickname, examples.Secret, password)
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

	// Build parameters using Args: AddString(name) followed by AddU32(age).
	args := srl.NewArgs(nil)
	args.AddString(name)
	args.AddU32(age)
	param := args.Serialise()

	// Call the smart contract function 'changeAge'. We'll use the client helper
	// `CallFunction`, which signs the operation using the provided account.
	// fee, maxGas, coins, expiryDelta, async
	fee := uint64(1000)
	maxGas := uint64(10000) // 0 means the client will estimate
	coins := uint64(0)      // no coins attached
	async := false
	desc := "example changeAge"
	expiryDelta, err := client.NextSlot(c)
	if err != nil {
		fmt.Printf("failed to get next slot: %v\n", err)
		os.Exit(1)
	}

	res, err := client.CallFunction(c, examples.SmartContractAddress, callerAddr, "changeAge", param, fee, maxGas, coins, expiryDelta, async, acc, password, desc)
	if err != nil {
		fmt.Printf("contract call failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Operation submitted. Event: %s\n", res.Event)
	fmt.Printf("Operation ID: %s\n", res.OperationResponse.OperationID)
}
