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

	fee := uint64(10_000_000)
	coins := uint64(1_000_000_000)
	desc := "example send coins"
	expiryDelta := uint64(100)

	op, err := transaction.New(examples.RecipientAddress, coins)
	if err != nil {
		fmt.Printf("error creating transaction operation: %v\n", err)
	}

	res, err := client.Call(c, expiryDelta, fee, op, acc, password, desc)
	if err != nil {
		fmt.Printf("send coin failed: %v\n", err)
		os.Exit(1)
	}
	pp.Printf("operation response: %v\n", res)
}

//dart
//Operation data: [128, 173, 226, 4, 100, 0, 0, 0, 2, 113, 251, 209, 172, 201, 199, 121, 170, 2, 40, 170, 110, 88, 109, 62, 193, 8, 217, 118, 52, 99, 28, 216, 174, 84, 7, 76, 35, 129, 140, 125, 128, 148, 235, 220, 3]
// Signature data: [0, 0, 0, 0, 4, 160, 248, 254, 0, 43, 119, 91, 15, 227, 124, 200, 197, 150, 132, 6, 97, 181, 136, 139, 205, 32, 131, 17, 64, 242, 30, 27, 77, 149, 9, 60, 139, 7, 70, 16, 101, 128, 173, 226, 4, 100, 0, 0, 0, 2, 113, 251, 209, 172, 201, 199, 121, 170, 2, 40, 170, 110, 88, 109, 62, 193, 8, 217, 118, 52, 99, 28, 216, 174, 84, 7, 76, 35, 129, 140, 125, 128, 148, 235, 220, 3]

//go
//operation data: [128 173 226 4 100 0 0 0 2 113 251 209 172 201 199 121 170 2 40 170 110 88 109 62 193 8 217 118 52 99 28 216 174 84 7 76 35 129 140 125 128 148 235 220 3]
// Signature data: [0 0 0 0 4 160 248 254 0 43 119 91 15 227 124 200 197 150 132 6 97 181 136 139 205 32 131 17 64 242 30 27 77 149 9 60 139 7 70 16 101 128 173 226 4 100 0 0 0 2 113 251 209 172 201 199 121 170 2 40 170 110 88 109 62 193 8 217 118 52 99 28 216 174 84 7 76 35 129 140 125 128 148 235 220 3]
