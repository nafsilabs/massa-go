package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nafsilabs/massa-go/client"
	"github.com/nafsilabs/massa-go/examples"
	"github.com/nafsilabs/massa-go/wallet"
)

func main() {
	fmt.Println("Massa Wallet-Integrated Operation Sender Example")
	fmt.Println("=================================================\n")

	// Step 1: Create or load a wallet
	fmt.Println("Step 1: Setting up wallet...")
	accountManager := wallet.NewAccountManager()

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
	account, err := w.AccountManager.GetAccount(nickname)
	if err != nil {
		account, err = w.ImportAccountFromPrivateKey(nickname, examples.Secret, password)
		if err != nil {
			fmt.Printf("failed to import account from secret: %v\n", err)
			os.Exit(1)
		}
	}

	// Prefer the known calling address from constants for clarity, but fall back
	// to the account's derived address if absent.
	callerAddr := examples.CallingAddress
	if callerAddr == "" && account.Address != nil {
		callerAddr = account.Address.String()
	}

	accountManager.ImportAccountFromPrivateKey(nickname, examples.Secret, password)
	// For this example, we'll create a new account
	// account, err := accountManager.CreateAccount("demo-account", "demo-password")
	// if err != nil {
	// 	log.Fatal("Failed to create account:", err)
	// }
	fmt.Printf("✓ Created account: %s\n", account.Nickname)
	fmt.Printf("  Address: %s\n\n", account.Address.String())

	// Step 2: Connect to Massa network
	fmt.Println("Step 2: Connecting to Massa buildnet...")
	grpcClient, err := client.NewMassaGrpcClient(client.ClientConfig{
		Address:        "buildnet.massa.net:33037",
		UseTLS:         false,
		DefaultTimeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatal("Failed to create gRPC client:", err)
	}
	defer grpcClient.Close()
	fmt.Println("✓ Connected successfully!\n")

	// Step 3: Create wallet-integrated operation sender
	fmt.Println("Step 3: Creating wallet operation sender...")
	walletSender := client.NewWalletOperationSender(grpcClient, accountManager)
	fmt.Println("✓ Wallet operation sender ready!\n")

	// Step 4: Demonstrate operations
	fmt.Println("Step 4: Example Operations (Demo Mode)")
	fmt.Println("---------------------------------------")
	fmt.Println("\nThe following examples show how to use the wallet-integrated sender:")
	fmt.Println("(Note: These are code examples only - not executing actual transactions)\n")

	// Example 1
	fmt.Println("1. Send Transaction")
	fmt.Println("   Transfer coins from your account to a recipient")
	fee := uint64(10_000_000)
	coins := uint64(1_000_000_000)
	expiry, _ := client.NextSlot(c)

	// Send transaction (automatic signing!)
	opID, err := walletSender.SendTransactionWithWallet(
		context.Background(), nickname, password,
		examples.RecipientAddress, coins, fee, expiry,
	)
	if err != nil {
		fmt.Printf("failed to send transaction operation: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Transaction operation sent successfully!")
	fmt.Println("   Operation ID:", opID)

	// // Example 2
	// fmt.Println("2. Buy Rolls")
	// fmt.Println("   Purchase rolls for staking")
	// fmt.Println()

	// // Example 3
	// fmt.Println("3. Call Smart Contract")
	// fmt.Println("   Call a function on a deployed smart contract")
	// fmt.Println()

	// // Example 4
	// fmt.Println("4. Deploy Smart Contract")
	// fmt.Println("   Deploy a new smart contract")
	// fmt.Println()

	// // Summary
	// fmt.Println(strings.Repeat("=", 60))
	// fmt.Println("Summary")
	// fmt.Println(strings.Repeat("=", 60))
	// fmt.Println("\nKey Advantages of WalletOperationSender:")
	// fmt.Println("1. ✓ Automatic signing using wallet accounts")
	// fmt.Println("2. ✓ No need to manually handle signatures and public keys")
	// fmt.Println("3. ✓ Secure password-protected key management")
	// fmt.Println("4. ✓ Simple, intuitive API")
	// fmt.Println("5. ✓ All operation types supported")

	// fmt.Println("\nTo actually send operations:")
	// fmt.Println("- Ensure your account has sufficient balance")
	// fmt.Println("- Use a valid recipient address")
	// fmt.Println("- Set appropriate fee and expiry period")
	// fmt.Println("- Handle the returned operation ID")

	// fmt.Println("\nComparison:")
	// fmt.Println()
	// fmt.Println("OLD WAY (Manual Signing):")
	// fmt.Println("  1. Create operation")
	// fmt.Println("  2. Sign with external wallet")
	// fmt.Println("  3. Get signature and public key")
	// fmt.Println("  4. Call SendOperation with signature")
	// fmt.Println()
	// fmt.Println("NEW WAY (Wallet-Integrated):")
	// fmt.Println("  1. Call SendXxxWithWallet(nickname, password, ...)")
	// fmt.Println("  2. Done! ✓")
}
