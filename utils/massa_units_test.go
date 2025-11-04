package utils

import (
	"math/big"
	"testing"
)

func TestFromMAS(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		expected int64
	}{
		{"1 Massa", 1.0, 1_000_000_000},
		{"0.5 Massa", 0.5, 500_000_000},
		{"1.5 Massa", 1.5, 1_500_000_000},
		{"0.000000001 Massa", 0.000000001, 1},
		{"100 Massa", 100.0, 100_000_000_000},
		{"0 Massa", 0.0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromMAS(tt.amount)
			expected := big.NewInt(tt.expected)

			if result.Cmp(expected) != 0 {
				t.Errorf("FromMAS(%f) = %s, want %s", tt.amount, result.String(), expected.String())
			}
		})
	}
}

func TestToMAS(t *testing.T) {
	tests := []struct {
		name     string
		amount   int64
		expected float64
	}{
		{"1 Massa", 1_000_000_000, 1.0},
		{"0.5 Massa", 500_000_000, 0.5},
		{"1.5 Massa", 1_500_000_000, 1.5},
		{"1 nanoMassa", 1, 0.000000001},
		{"100 Massa", 100_000_000_000, 100.0},
		{"0 Massa", 0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amount := big.NewInt(tt.amount)
			result := ToMAS(amount)

			// Use a small epsilon for floating point comparison
			epsilon := 0.000000001
			if abs(result-tt.expected) > epsilon {
				t.Errorf("ToMAS(%s) = %f, want %f", amount.String(), result, tt.expected)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	tests := []float64{1.0, 0.5, 1.5, 100.0, 0.123456789}

	for _, original := range tests {
		t.Run("RoundTrip", func(t *testing.T) {
			// Convert to nanoMassa and back
			nanoMassa := FromMAS(original)
			result := ToMAS(nanoMassa)

			// Allow small floating point error
			epsilon := 0.000000001
			if abs(result-original) > epsilon {
				t.Errorf("Round trip failed: %f -> %s -> %f", original, nanoMassa.String(), result)
			}
		})
	}
}

func TestFromNanoMAS(t *testing.T) {
	tests := []struct {
		name     string
		nanoMAS  uint64
		expected int64
	}{
		{"1 billion nanoMassa", 1_000_000_000, 1_000_000_000},
		{"500 million nanoMassa", 500_000_000, 500_000_000},
		{"1 nanoMassa", 1, 1},
		{"0 nanoMassa", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromNanoMAS(tt.nanoMAS)
			expected := big.NewInt(tt.expected)

			if result.Cmp(expected) != 0 {
				t.Errorf("FromNanoMAS(%d) = %s, want %s", tt.nanoMAS, result.String(), expected.String())
			}
		})
	}
}

func TestToNanoMAS(t *testing.T) {
	tests := []struct {
		name     string
		amount   int64
		expected uint64
	}{
		{"1 billion", 1_000_000_000, 1_000_000_000},
		{"500 million", 500_000_000, 500_000_000},
		{"1", 1, 1},
		{"0", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amount := big.NewInt(tt.amount)
			result := ToNanoMAS(amount)

			if result != tt.expected {
				t.Errorf("ToNanoMAS(%s) = %d, want %d", amount.String(), result, tt.expected)
			}
		})
	}
}

func TestToNanoMASWithNil(t *testing.T) {
	result := ToNanoMAS(nil)
	if result != 0 {
		t.Errorf("ToNanoMAS(nil) = %d, want 0", result)
	}
}

func TestToMASWithNil(t *testing.T) {
	result := ToMAS(nil)
	if result != 0.0 {
		t.Errorf("ToMAS(nil) = %f, want 0.0", result)
	}
}

func TestParseMassa(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		expected uint64
	}{
		{"1 Massa", 1.0, 1_000_000_000},
		{"0.5 Massa", 0.5, 500_000_000},
		{"1.5 Massa", 1.5, 1_500_000_000},
		{"100 Massa", 100.0, 100_000_000_000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseMassa(tt.amount)

			if result != tt.expected {
				t.Errorf("ParseMassa(%f) = %d, want %d", tt.amount, result, tt.expected)
			}
		})
	}
}

// Helper function for floating point comparison
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
