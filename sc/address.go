package sc

import (
	"errors"
	"fmt"
	"strings"
)

// Address represents a Massa blockchain address
type Address struct {
	value string
}

// NewAddress creates a new Address from a string
func NewAddress(addr string) (*Address, error) {
	if !ValidateAddress(addr) {
		return nil, errors.New("invalid address format")
	}
	return &Address{value: addr}, nil
}

// MustNewAddress creates a new Address from a string, panicking if invalid
func MustNewAddress(addr string) *Address {
	address, err := NewAddress(addr)
	if err != nil {
		panic(err)
	}
	return address
}

// String returns the string representation of the address
func (a *Address) String() string {
	return a.value
}

// Bytes returns the address as bytes for serialization
func (a *Address) Bytes() []byte {
	return []byte(a.value)
}

// Equal checks if two addresses are equal
func (a *Address) Equal(other *Address) bool {
	if a == nil || other == nil {
		return a == other
	}
	return a.value == other.value
}

// ValidateAddress checks if the given address string is valid
func ValidateAddress(address string) bool {
	if address == "" {
		return false
	}

	// Check if it's a Massa address (starts with AU and has correct length)
	if strings.HasPrefix(address, "AU") && len(address) == 49 {
		// Basic format validation - in a real implementation we'd do full base58 validation
		// For now, just check that it contains valid base58 characters
		validChars := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
		for _, char := range address[2:] {
			if !strings.ContainsRune(validChars, char) {
				return false
			}
		}
		return true
	}

	// For WebAssembly runtime, we'll use the native validation
	if isWasmRuntime() {
		ptr := newWasmString(address)
		result := wasmValidateAddress(ptr)
		return result != 0
	}

	return false
} // PublicKeyToAddress converts a public key to an address
func PublicKeyToAddress(pubKey string) (*Address, error) {
	if pubKey == "" {
		return nil, errors.New("public key cannot be empty")
	}

	if isWasmRuntime() {
		ptr := newWasmString(pubKey)
		resultPtr := wasmPublicKeyToAddress(ptr)
		addressStr := wasmStringFromPtr(resultPtr)
		return NewAddress(addressStr)
	}

	// For non-WASM environments (testing), return a mock address
	return NewAddress("AU12UBnqTHDQALpocVBnkPNy7y5CndUJQTLutaVDDFgMJcq5kQiKq")
}

// IsAddressEoa checks if the given address is an Externally Owned Account (EOA)
func IsAddressEoa(address string) bool {
	if !ValidateAddress(address) {
		return false
	}

	if isWasmRuntime() {
		ptr := newWasmString(address)
		result := wasmIsAddressEoa(ptr)
		return result != 0
	}

	// For non-WASM environments, assume it's an EOA if it's valid
	return true
}

// AddressFromBytes creates an Address from a byte slice
func AddressFromBytes(data []byte) (*Address, error) {
	return NewAddress(string(data))
}

// Serialize returns the address as a byte slice for storage/transmission
func (a *Address) Serialize() []byte {
	return []byte(a.value)
}

// Clone creates a copy of the address
func (a *Address) Clone() *Address {
	if a == nil {
		return nil
	}
	return &Address{value: a.value}
}

// Format implements custom formatting for the address
func (a *Address) Format(f fmt.State, verb rune) {
	switch verb {
	case 's', 'v':
		fmt.Fprint(f, a.value)
	case 'q':
		fmt.Fprintf(f, "%q", a.value)
	default:
		fmt.Fprintf(f, "%%!%c(Address=%s)", verb, a.value)
	}
}
