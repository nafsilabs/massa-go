package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nafsilabs/massa-go/client"
	srl "github.com/nafsilabs/massa-go/sc/serialisation"
	"github.com/nafsilabs/massa-go/utils"
)

func main() {
	smartContractAddress := "AS12cdcRczrDe3TxeGqQU6TFWVuYnVN4SeSvJdQNvEvHZ2YwMafFa"
	callerAddress := "AU1y3oYTgK8RGzLWFVAGL3JxHLdTVKxBPHrwbo2Kj8a2CSbeMug"

	// Build the argument using Args (writes a 4-byte little-endian length then data)
	name := "alice"
	args := srl.NewArgs(nil)
	args.AddString(name)
	fee := 0.01
	coins := 1.0
	params := args.Serialise()
	targetFunction := "getAge"

	cfg := client.ClientConfig{
		Address:        "buildnet.massa.net:33037",
		UseTLS:         false,
		DefaultTimeout: 30 * time.Second,
		ChainID:        utils.BUILDNET,
	}
	grpc, err := client.NewMassaClient(&cfg)
	if err != nil {
		log.Fatal("Failed to create gRPC client:", err)
	}
	defer grpc.Close()
	fmt.Println("Connected successfully!")

	resp, err := grpc.ReadOnlyCallSC(context.Background(), smartContractAddress, targetFunction, params, coins, fee, callerAddress)

	fmt.Printf("decoding response: %v\n", resp.GetCallResult())

	// Use helper to extract the returned payload bytes from the ReadOnlyResult.
	payload := resp.GetCallResult()
	if err != nil {
		fmt.Printf("error extracting payload: %v\n", err)
		return
	}

	argsResp := srl.NewArgs(payload)
	age, err := argsResp.NextU32()
	if err != nil {
		fmt.Printf("error deserializing response: %v\n", err)
		return
	}

	fmt.Printf("Age of %s is %d years\n", name, age)
}
