package sc

import (
	"encoding/json"
	"fmt"
)

// Context provides functions for interacting with the execution context
// of a smart contract on the Massa blockchain.

// Caller returns the address of the caller of the current smart contract
func Caller() *Address {
	if isWasmRuntime() {
		// In the AssemblyScript version, this comes from the call stack
		// For now, we'll implement a simplified version
		return MustNewAddress("AU12UBnqTHDQALpocVBnkPNy7y5CndUJQTLutaVDDFgMJcq5kQiKq")
	}
	// In non-WASM mode, return a mock address for testing
	return MustNewAddress("AU12UBnqTHDQALpocVBnkPNy7y5CndUJQTLutaVDDFgMJcq5kQiKq")
}

// Callee returns the address of the current smart contract
func Callee() *Address {
	if isWasmRuntime() {
		// In the AssemblyScript version, this comes from the call stack
		// For now, we'll implement a simplified version
		return MustNewAddress("AU12UBnqTHDQALpocVBnkPNy7y5CndUJQTLutaVDDFgMJcq5kQiKr")
	}
	// In non-WASM mode, return a mock address for testing
	return MustNewAddress("AU12UBnqTHDQALpocVBnkPNy7y5CndUJQTLutaVDDFgMJcq5kQiKr")
}

// TransactionCreator returns the address of the transaction creator
func TransactionCreator() *Address {
	// In Massa, the transaction creator is typically the first address in the call stack
	stack := AddressStack()
	if len(stack) > 0 {
		return stack[0]
	}
	// Fallback to caller if no stack available
	return Caller()
}

// TransferredCoins returns the amount of coins transferred in the current call
func TransferredCoins() uint64 {
	if isWasmRuntime() {
		return wasmGetCallCoins()
	}
	// In non-WASM mode, return 0 for testing
	return 0
}

// Timestamp returns the current timestamp of the blockchain
func Timestamp() uint64 {
	if isWasmRuntime() {
		return wasmGetTime()
	}
	// In non-WASM mode, return a mock timestamp for testing
	return 1640995200000 // Jan 1, 2022 timestamp in milliseconds
}

// RemainingGas returns the remaining gas available for execution
func RemainingGas() uint64 {
	if isWasmRuntime() {
		return wasmRemainingGas()
	}
	// In non-WASM mode, return a large number for testing
	return 1000000000000000000
}

// CurrentPeriod returns the current period of the network
func CurrentPeriod() uint64 {
	if isWasmRuntime() {
		return wasmGetCurrentPeriod()
	}
	// In non-WASM mode, return 0 for testing
	return 0
}

// CurrentThread returns the current thread of the execution context
func CurrentThread() uint8 {
	if isWasmRuntime() {
		return uint8(wasmGetCurrentThread())
	}
	// In non-WASM mode, return thread 1 for testing
	return 1
}

// ChainId returns the current chain ID
func ChainId() uint64 {
	if isWasmRuntime() {
		return wasmChainId()
	}
	// In non-WASM mode, return mainnet chain ID for testing
	return 77658366
}

// OwnedAddresses returns the list of addresses owned by the current execution context
func OwnedAddresses() []*Address {
	if isWasmRuntime() {
		resultPtr := wasmGetOwnedAddresses()
		jsonStr := wasmStringFromPtr(resultPtr)
		var addressStrings []string
		if err := json.Unmarshal([]byte(jsonStr), &addressStrings); err != nil {
			return nil
		}

		addresses := make([]*Address, len(addressStrings))
		for i, addrStr := range addressStrings {
			addr, err := NewAddress(addrStr)
			if err == nil {
				addresses[i] = addr
			}
		}
		return addresses
	}

	// In non-WASM mode, return mock addresses for testing
	addr1, _ := NewAddress("AU12UBnqTHDQALpocVBnkPNy7y5CndUJQTLutaVDDFgMJcq5kQiKq")
	return []*Address{addr1}
}

// AddressStack returns the current call stack as a list of addresses
func AddressStack() []*Address {
	if isWasmRuntime() {
		resultPtr := wasmGetCallStack()
		jsonStr := wasmStringFromPtr(resultPtr)
		var addressStrings []string
		if err := json.Unmarshal([]byte(jsonStr), &addressStrings); err != nil {
			return nil
		}

		addresses := make([]*Address, len(addressStrings))
		for i, addrStr := range addressStrings {
			addr, err := NewAddress(addrStr)
			if err == nil {
				addresses[i] = addr
			}
		}
		return addresses
	}

	// In non-WASM mode, return mock call stack for testing
	caller, _ := NewAddress("AU12UBnqTHDQALpocVBnkPNy7y5CndUJQTLutaVDDFgMJcq5kQiKq")
	callee, _ := NewAddress("AU12UBnqTHDQALpocVBnkPNy7y5CndUJQTLutaVDDFgMJcq5kQiKr")
	return []*Address{caller, callee}
}

// IsDeployingContract returns true if the contract is being deployed
func IsDeployingContract() bool {
	// In Massa, a contract is being deployed if the call stack has only one address
	// and it's the same as the current contract address
	stack := AddressStack()
	if len(stack) == 1 {
		callee := Callee()
		return stack[0].Equal(callee)
	}
	return false
}

// UnsafeRandom returns a pseudo-random number
// WARNING: This function is not cryptographically secure and should not be used
// for security-critical applications
func UnsafeRandom() int64 {
	if isWasmRuntime() {
		return wasmUnsafeRandom()
	}
	// In non-WASM mode, return a deterministic value for testing
	return 42
}

// GetOriginOperationId returns the ID of the operation that started this execution
func GetOriginOperationId() string {
	if isWasmRuntime() {
		resultPtr := wasmGetOriginOperationId()
		return wasmStringFromPtr(resultPtr)
	}
	// In non-WASM mode, return a mock operation ID for testing
	return "O1234567890abcdef"
}

// Slot represents a specific time slot in the Massa blockchain
type Slot struct {
	Period uint64
	Thread uint8
}

// NewSlot creates a new Slot
func NewSlot(period uint64, thread uint8) Slot {
	return Slot{
		Period: period,
		Thread: thread,
	}
}

// CurrentSlot returns the current slot
func CurrentSlot() Slot {
	return Slot{
		Period: CurrentPeriod(),
		Thread: CurrentThread(),
	}
}

// String returns a string representation of the slot
func (s Slot) String() string {
	return fmt.Sprintf("Slot{Period: %d, Thread: %d}", s.Period, s.Thread)
}

// Equal checks if two slots are equal
func (s Slot) Equal(other Slot) bool {
	return s.Period == other.Period && s.Thread == other.Thread
}

// Before checks if this slot comes before another slot
func (s Slot) Before(other Slot) bool {
	if s.Period < other.Period {
		return true
	}
	if s.Period == other.Period {
		return s.Thread < other.Thread
	}
	return false
}

// After checks if this slot comes after another slot
func (s Slot) After(other Slot) bool {
	return other.Before(s)
}
