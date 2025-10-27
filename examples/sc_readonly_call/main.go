package main

import (
	"fmt"

	"github.com/nafsilabs/massa-go/client"
	srl "github.com/nafsilabs/massa-go/sc/serialisation"
)

func main() {
	smartContractAddress := "AS12cdcRczrDe3TxeGqQU6TFWVuYnVN4SeSvJdQNvEvHZ2YwMafFa"
	callerAddress := "AU1y3oYTgK8RGzLWFVAGL3JxHLdTVKxBPHrwbo2Kj8a2CSbeMug"

	c := client.NewClient(false)

	// Build the argument using Args (writes a 4-byte little-endian length then data)
	name := "alice"
	args := srl.NewArgs(nil)
	args.AddString(name)
	param := client.JSONableSlice(args.Serialise())

	// ReadOnlyCallSC expects a []byte parameter; convert from JSONableSlice.
	resp, err := client.ReadOnlyCallSC(smartContractAddress, "getAge", []byte(param), "1", "1", callerAddress, c)
	if err != nil {
		fmt.Printf("error calling function: %v\n", err)
		return
	}

	fmt.Printf("response: %+v\n", resp)

	//fmt.Printf("decoding response: %v\n", resp.Result.Ok)

	// Use helper to extract the returned payload bytes from the ReadOnlyResult.
	payload, err := client.ReadOnlyResultToBytes(resp)
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

	fmt.Printf("Age of %s is %d\n", name, age)
}
