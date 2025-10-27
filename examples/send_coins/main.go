package main

import (
	"fmt"
	"os"

	"github.com/k0kubun/pp"
	"github.com/nafsilabs/massa-go/client"
	"github.com/nafsilabs/massa-go/client/sendoperation/transaction"
	examples "github.com/nafsilabs/massa-go/examples"
	"github.com/nafsilabs/massa-go/wallet"
)

func main() {
	// Configuration — adapt these to your environment
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

	fee := uint64(1000)
	coins := uint64(10000)
	desc := "example send coins"
	expiryDelta, err := client.NextSlot(c)
	if err != nil {
		fmt.Printf("failed to get next slot: %v\n", err)
		os.Exit(1)
	}

	op, err := transaction.New(examples.RecipientAddress, coins)
	if err != nil {
		fmt.Printf("error creating transaction operation: %v\n", err)
	}

	res, err := client.Call(c, expiryDelta, fee, op, acc, password, desc)
	if err != nil {
		fmt.Printf("contract call failed: %v\n", err)
		os.Exit(1)
	}
	pp.Printf("operation response: %v\n", res)
}
