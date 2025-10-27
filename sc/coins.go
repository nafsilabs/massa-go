package sc

import (
	"fmt"
)

// Coin operations for managing balance and transfers in Massa smart contracts

// Balance returns the balance of the current address in the smallest unit (micro-Massa)
func Balance() uint64 {
	if isWasmRuntime() {
		return wasmGetBalance()
	}
	// In non-WASM mode, return a mock balance for testing
	return 1000000000 // 1000 Massa in micro-Massa
}

// BalanceOf returns the balance of the specified address in the smallest unit (micro-Massa)
func BalanceOf(address *Address) uint64 {
	if isWasmRuntime() {
		addressPtr := newWasmString(address.String())
		return wasmGetBalanceFor(addressPtr)
	}
	// In non-WASM mode, return a mock balance for testing
	return 1000000000 // 1000 Massa in micro-Massa
}

// TransferCoins transfers coins from the current address to the specified address
func TransferCoins(to *Address, amount uint64) {
	if isWasmRuntime() {
		toPtr := newWasmString(to.String())
		wasmTransferCoins(toPtr, amount)
	}
	// In non-WASM mode, this would transfer in mock system for testing
}

// TransferCoinsOf transfers coins from one address to another
// This can only be called if the current execution context has permission
func TransferCoinsOf(from *Address, to *Address, amount uint64) {
	if isWasmRuntime() {
		fromPtr := newWasmString(from.String())
		toPtr := newWasmString(to.String())
		wasmTransferCoinsFor(fromPtr, toPtr, amount)
	}
	// In non-WASM mode, this would transfer in mock system for testing
}

// Convenience functions for common operations

// TransferToString transfers coins to an address specified as a string
func TransferToString(to string, amount uint64) error {
	address, err := NewAddress(to)
	if err != nil {
		return err
	}
	TransferCoins(address, amount)
	return nil
}

// TransferFromString transfers coins from an address specified as a string
func TransferFromString(from string, to string, amount uint64) error {
	fromAddr, err := NewAddress(from)
	if err != nil {
		return err
	}
	toAddr, err := NewAddress(to)
	if err != nil {
		return err
	}
	TransferCoinsOf(fromAddr, toAddr, amount)
	return nil
}

// Unit conversion helpers

const (
	// MicroMassa is the smallest unit (1 micro-Massa)
	MicroMassa uint64 = 1

	// MilliMassa represents 1 milli-Massa (1000 micro-Massa)
	MilliMassa uint64 = 1000

	// Massa represents 1 Massa (1,000,000 micro-Massa)
	Massa uint64 = 1000000
)

// ToMicroMassa converts Massa to micro-Massa
func ToMicroMassa(massa float64) uint64 {
	return uint64(massa * float64(Massa))
}

// ToMilliMassa converts micro-Massa to milli-Massa
func ToMilliMassa(microMassa uint64) float64 {
	return float64(microMassa) / float64(MilliMassa)
}

// ToMassa converts micro-Massa to Massa
func ToMassa(microMassa uint64) float64 {
	return float64(microMassa) / float64(Massa)
}

// FormatBalance formats a balance in micro-Massa to a human-readable string
func FormatBalance(microMassa uint64) string {
	if microMassa >= Massa {
		massa := ToMassa(microMassa)
		return fmt.Sprintf("%.6f MAS", massa)
	} else if microMassa >= MilliMassa {
		milliMassa := ToMilliMassa(microMassa)
		return fmt.Sprintf("%.3f mMAS", milliMassa)
	} else {
		return fmt.Sprintf("%d µMAS", microMassa)
	}
}

// ParseMassa parses a string amount in Massa and returns micro-Massa
func ParseMassa(amount string) (uint64, error) {
	// This is a simple implementation - in a real SDK you'd want more robust parsing
	var massa float64
	var unit string

	// Try to parse as "123.456 MAS" format
	n, err := fmt.Sscanf(amount, "%f %s", &massa, &unit)
	if err != nil || n != 2 {
		// Try to parse as just a number (assume MAS)
		n, err = fmt.Sscanf(amount, "%f", &massa)
		if err != nil || n != 1 {
			return 0, fmt.Errorf("invalid amount format: %s", amount)
		}
		unit = "MAS"
	}

	switch unit {
	case "MAS", "MASSA":
		return ToMicroMassa(massa), nil
	case "mMAS", "mMASSA":
		return uint64(massa * float64(MilliMassa)), nil
	case "µMAS", "uMAS", "microMAS":
		return uint64(massa), nil
	default:
		return 0, fmt.Errorf("unknown unit: %s", unit)
	}
}

// Balance checking helpers

// HasSufficientBalance checks if an address has at least the specified amount
func HasSufficientBalance(address *Address, amount uint64) bool {
	return BalanceOf(address) >= amount
}

// HasSufficientBalanceString checks if an address (as string) has at least the specified amount
func HasSufficientBalanceString(address string, amount uint64) (bool, error) {
	addr, err := NewAddress(address)
	if err != nil {
		return false, err
	}
	return HasSufficientBalance(addr, amount), nil
}

// GetTransferableBalan返回可转账余额 (current balance minus any locked amounts)
// For now, this is the same as the regular balance, but could be extended
// to account for locked funds, staking, etc.
func GetTransferableBalance(address *Address) uint64 {
	return BalanceOf(address)
}

// Batch transfer operations

// TransferBatch performs multiple transfers in a single operation
type Transfer struct {
	To     *Address
	Amount uint64
}

// TransferBatch transfers coins to multiple addresses
func TransferBatch(transfers []Transfer) error {
	// Validate all transfers first
	totalAmount := uint64(0)
	for _, transfer := range transfers {
		if transfer.To == nil {
			return fmt.Errorf("invalid transfer: nil address")
		}
		if transfer.Amount == 0 {
			return fmt.Errorf("invalid transfer: zero amount")
		}
		totalAmount += transfer.Amount
	}

	// Check if we have sufficient balance
	currentBalance := Balance()
	if currentBalance < totalAmount {
		return fmt.Errorf("insufficient balance: have %d, need %d", currentBalance, totalAmount)
	}

	// Perform all transfers
	for _, transfer := range transfers {
		TransferCoins(transfer.To, transfer.Amount)
	}

	return nil
}

// Event types for coin operations
type CoinTransferEvent struct {
	From   *Address `json:"from"`
	To     *Address `json:"to"`
	Amount uint64   `json:"amount"`
}

// EmitTransferEvent emits a transfer event (this would typically be called internally)
func EmitTransferEvent(from, to *Address, amount uint64) {
	// This would use the event system to emit the transfer event
	// For now, we'll just log it if in debug mode
	if isDebugMode() {
		// Use the logging function from the events module
		message := fmt.Sprintf("Transfer: %s -> %s: %d µMAS",
			from.String(), to.String(), amount)
		// We'll implement Print in the events module
		_ = message // Prevent unused variable error for now
	}
} // isDebugMode checks if we're in debug mode (for testing purposes)
func isDebugMode() bool {
	// This could be set via environment variables or compile-time flags
	return false
}
