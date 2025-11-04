package utils

import (
	"math"
	"math/big"
)

// Constants for Massa unit conversions
const (
	// DecimalScale is the number of decimal places in Massa (9 for nanoMassa)
	DecimalScale = 9

	// DecimalFactor represents the maximum scale factor
	DecimalFactor = math.MaxFloat64

	// NanoMassaPerMassa is the number of nanoMassa in 1 Massa (10^9)
	NanoMassaPerMassa = 1_000_000_000
)

// FromMAS converts Massa (as float64) to nanoMassa (as *big.Int)
// Example: FromMAS(1.5) returns 1500000000 nanoMassa
func FromMAS(amount float64) *big.Int {
	return DoubleToMassaInt(amount)
}

// ToMAS converts nanoMassa (as *big.Int) to Massa (as float64)
// Example: ToMAS(1500000000) returns 1.5 Massa
func ToMAS(amount *big.Int) float64 {
	if amount == nil {
		return 0.0
	}

	// Scale factor: 10^9
	scaleFactor := new(big.Int).Exp(big.NewInt(10), big.NewInt(DecimalScale), nil)

	// Convert amount to float64 and divide by scale factor
	amountFloat := new(big.Float).SetInt(amount)
	scaleFloat := new(big.Float).SetInt(scaleFactor)

	result := new(big.Float).Quo(amountFloat, scaleFloat)

	massaAmount, _ := result.Float64()
	return massaAmount
}

// DoubleToMassaInt converts a float64 Massa amount to nanoMassa as *big.Int
// This multiplies the amount by 10^9 to convert from Massa to nanoMassa
func DoubleToMassaInt(amount float64) *big.Int {
	// Multiply by 10^9 to convert Massa to nanoMassa
	nanoMassa := amount * NanoMassaPerMassa

	// Convert to big.Int
	return big.NewInt(int64(nanoMassa))
}

// MassaIntToDouble converts nanoMassa (*big.Int) to Massa as float64
// This is an alias for ToMAS for consistency with naming conventions
func MassaIntToDouble(amount *big.Int) float64 {
	return ToMAS(amount)
}

// FromNanoMAS converts nanoMassa (as uint64) to *big.Int
// This is useful when working with uint64 amounts
func FromNanoMAS(nanoMassa uint64) *big.Int {
	return big.NewInt(int64(nanoMassa))
}

// ToNanoMAS converts *big.Int to nanoMassa (as uint64)
// Returns 0 if amount is nil or exceeds uint64 range
func ToNanoMAS(amount *big.Int) uint64 {
	if amount == nil {
		return 0
	}

	// Check if it fits in uint64
	if !amount.IsUint64() {
		return 0
	}

	return amount.Uint64()
}

// FormatMassa formats a Massa amount (float64) as a string with specified decimal places
func FormatMassa(amount float64, decimals int) string {
	if decimals < 0 {
		decimals = DecimalScale
	}

	format := "%." + string(rune(decimals+'0')) + "f"
	return string(rune(format[0])) + " MAS"
}

// ParseMassa parses a float64 amount and returns nanoMassa as uint64
func ParseMassa(amount float64) uint64 {
	nanoMassa := amount * NanoMassaPerMassa
	return uint64(nanoMassa)
}
