package main

import (
	massa "github.com/nafsilabs/massa-go/sc"
)

// main function is called automatically when the smart contract is executed by the blockchain.
// The function argument is unused and can be safely ignored.
func main() {
	// Generate a "Hello, World!" event
	massa.GenerateEvent("Hello, World!")

	// Also print for debugging purposes
	massa.Print("Hello World smart contract executed successfully!")

	// Demonstrate various SDK functionalities
	demonstrateSDK()
}

// demonstrateSDK shows various features of the Massa Go SDK
func demonstrateSDK() {
	// Context information
	caller := massa.Caller()
	callee := massa.Callee()
	timestamp := massa.Timestamp()

	massa.Printf("Caller: %s", caller.String())
	massa.Printf("Callee: %s", callee.String())
	massa.Printf("Timestamp: %d", timestamp)

	// Storage operations
	massa.SetString("greeting", "Hello from Go!")
	greeting := massa.GetString("greeting")
	massa.Printf("Stored greeting: %s", greeting)

	// Balance information
	balance := massa.Balance()
	massa.Printf("Current balance: %s", massa.FormatBalance(balance))

	// Create structured events
	err := massa.CreateEvent("contract_execution", map[string]interface{}{
		"caller":    caller.String(),
		"callee":    callee.String(),
		"timestamp": timestamp,
		"balance":   balance,
		"message":   "Contract executed successfully",
	})
	if err != nil {
		massa.LogError("Failed to create event", err)
	}

	// Demonstrate cryptographic functions
	data := []byte("Hello, cryptography!")
	hash := massa.Blake3(data)
	massa.Printf("BLAKE3 hash: %s", massa.BytesToHex(hash))

	// Demonstrate address validation
	if massa.ValidateAddress(caller.String()) {
		massa.LogInfo("Caller address is valid")
	}

	// Store some structured data
	contractInfo := map[string]interface{}{
		"name":        "Hello World Contract",
		"version":     "1.0.0",
		"deployed_at": timestamp,
		"deployer":    caller.String(),
	}

	err = massa.SetJSON("contract_info", contractInfo)
	if err != nil {
		massa.LogError("Failed to store contract info", err)
	} else {
		massa.LogInfo("Contract info stored successfully")
	}

	// Emit a custom event with the contract information
	massa.EmitCustomEvent("contract_info_stored", map[string]interface{}{
		"contract_address": callee.String(),
		"info":             contractInfo,
	})
}

// Additional functions to demonstrate contract capabilities

// GetContractInfo returns information about this contract
func GetContractInfo() map[string]interface{} {
	var info map[string]interface{}
	err := massa.GetJSON("contract_info", &info)
	if err != nil {
		massa.LogError("Failed to retrieve contract info", err)
		return nil
	}
	return info
}

// UpdateGreeting updates the stored greeting message
func UpdateGreeting(newGreeting string) {
	// Only allow the deployer to update the greeting
	caller := massa.Caller()

	var contractInfo map[string]interface{}
	err := massa.GetJSON("contract_info", &contractInfo)
	if err != nil {
		massa.LogError("Failed to get contract info", err)
		return
	}

	deployer, ok := contractInfo["deployer"].(string)
	if !ok {
		massa.LogError("Invalid deployer information in contract info", nil)
		return
	}

	if caller.String() != deployer {
		massa.LogError("Only deployer can update greeting", map[string]interface{}{
			"caller":   caller.String(),
			"deployer": deployer,
		})
		return
	}

	// Update the greeting
	massa.SetString("greeting", newGreeting)
	massa.LogInfo("Greeting updated", map[string]interface{}{
		"old_greeting": massa.GetString("greeting"),
		"new_greeting": newGreeting,
		"updated_by":   caller.String(),
	})

	// Emit an event
	massa.EmitCustomEvent("greeting_updated", map[string]interface{}{
		"new_greeting": newGreeting,
		"updated_by":   caller.String(),
		"timestamp":    massa.Timestamp(),
	})
}

// GetGreeting returns the current greeting message
func GetGreeting() string {
	greeting := massa.GetString("greeting")
	if greeting == "" {
		return "No greeting set"
	}
	return greeting
}

// TransferFunds transfers funds to another address (example of coin operations)
func TransferFunds(to string, amount uint64) {
	caller := massa.Caller()

	// Validate the target address
	if !massa.ValidateAddress(to) {
		massa.LogError("Invalid target address", map[string]interface{}{
			"address": to,
		})
		return
	}

	// Check if we have sufficient balance
	balance := massa.Balance()
	if balance < amount {
		massa.LogError("Insufficient balance", map[string]interface{}{
			"balance":  balance,
			"required": amount,
		})
		return
	}

	// Create address object
	toAddr, err := massa.NewAddress(to)
	if err != nil {
		massa.LogError("Failed to create address object", err)
		return
	}

	// Perform transfer
	massa.TransferCoins(toAddr, amount)

	// Log the transfer
	massa.LogInfo("Transfer completed", map[string]interface{}{
		"from":   caller.String(),
		"to":     to,
		"amount": amount,
	})

	// Emit transfer event
	massa.EmitCustomEvent("funds_transferred", map[string]interface{}{
		"from":      caller.String(),
		"to":        to,
		"amount":    amount,
		"timestamp": massa.Timestamp(),
	})
}
