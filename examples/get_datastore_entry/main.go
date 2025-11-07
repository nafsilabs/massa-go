package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nafsilabs/massa-go/client"
	"github.com/nafsilabs/massa-go/sc/serialisation"
	"github.com/nafsilabs/massa-go/utils"
)

func main() {
	// Replace with your deployed smart contract address
	smartContractAddress := "AS12cdcRczrDe3TxeGqQU6TFWVuYnVN4SeSvJdQNvEvHZ2YwMafFa"

	// key name we want to read from the contract's datastore
	name := "alice"

	ar := serialisation.NewArgs([]byte{})
	ar.AddString(name)
	key := ar.Serialise()

	cfg := client.ClientConfig{
		Address:        "buildnet.massa.net:33037",
		UseTLS:         false,
		DefaultTimeout: 30 * time.Second,
		ChainID:        utils.BUILDNET,
	}

	mc, err := client.NewMassaClient(&cfg)
	if err != nil {
		log.Fatal("failed to create gRPC client:", err)
	}
	defer mc.Close()

	ctx := context.Background()

	entries, err := mc.GetDatastoreEntries(ctx, smartContractAddress, [][]byte{key})
	if err != nil {
		log.Fatalf("error fetching datastore entries: %v", err)
	}

	if len(entries) == 0 {
		fmt.Printf("no datastore entry found for key '%s'\n", name)
		return
	}

	e := entries[0]
	final := e.GetFinalValue()
	candidate := e.GetCandidateValue()

	fmt.Printf("final bytes: %v\n", final)
	fmt.Printf("candidate bytes: %v\n", candidate)

	// Try to decode final value as a u32 (common pattern for simple examples).
	if len(final) > 0 {
		args := serialisation.NewArgs(final)
		if age, err := args.NextU32(); err == nil {
			fmt.Printf("decoded final value (u32): %d\n", age)
		} else {
			fmt.Printf("final raw (hex): %x\n", final)
		}
	}
}
