package sc

import (
	"testing"
)

// TestBasicOperations demonstrates basic SDK functionality
func TestBasicOperations(t *testing.T) {
	// Test address operations
	addr, err := NewAddress("AU12UBnqTHDQALpocVBnkPNy7y5CndUJQTLutaVDDFgMJcq5kQiKq")
	if err != nil {
		t.Fatalf("Failed to create address: %v", err)
	}

	if addr.String() != "AU12UBnqTHDQALpocVBnkPNy7y5CndUJQTLutaVDDFgMJcq5kQiKq" {
		t.Errorf("Address string mismatch")
	}

	// Test storage operations
	SetString("test_key", "test_value")
	value := GetString("test_key")
	if value != "test_value" {
		t.Errorf("Storage operation failed: expected 'test_value', got '%s'", value)
	}

	// Test context functions (these will use mock implementations)
	caller := Caller()
	if caller == nil {
		t.Error("Caller should not be nil")
	}

	// Test balance operations
	balance := Balance()
	if balance < 0 {
		t.Error("Balance should not be negative")
	}

	// Test event generation (should not panic)
	GenerateEvent("Test event")
	Print("Test print message")

	// Test cryptographic functions
	data := []byte("test data")
	hash := Blake3(data)
	if len(hash) != 32 {
		t.Errorf("BLAKE3 hash should be 32 bytes, got %d", len(hash))
	}

	// Test hex conversion
	hexStr := BytesToHex(hash)
	if len(hexStr) != 64 {
		t.Errorf("Hex string should be 64 characters, got %d", len(hexStr))
	}

	converted, err := HexToBytes(hexStr)
	if err != nil {
		t.Fatalf("Failed to convert hex back to bytes: %v", err)
	}

	if !BytesEqual(hash, converted) {
		t.Error("Hex conversion round-trip failed")
	}
}

// TestAddressValidation tests address validation functionality
func TestAddressValidation(t *testing.T) {
	validAddr := "AU12UBnqTHDQALpocVBnkPNy7y5CndUJQTLutaVDDFgMJcq5kQiKq"
	invalidAddr := "invalid_address"

	if !ValidateAddress(validAddr) {
		t.Error("Valid address should pass validation")
	}

	if ValidateAddress(invalidAddr) {
		t.Error("Invalid address should fail validation")
	}

	// Test address creation
	_, err := NewAddress(validAddr)
	if err != nil {
		t.Errorf("Should be able to create valid address: %v", err)
	}

	_, err = NewAddress(invalidAddr)
	if err == nil {
		t.Error("Should not be able to create invalid address")
	}
}

// TestStorageOperations tests storage functionality
func TestStorageOperations(t *testing.T) {
	key := "test_storage_key"
	value := "test_storage_value"

	// Test basic string storage
	SetString(key, value)
	retrieved := GetString(key)
	if retrieved != value {
		t.Errorf("Storage failed: expected '%s', got '%s'", value, retrieved)
	}

	// Test existence check
	if !HasString(key) {
		t.Error("Key should exist after setting")
	}

	// Test JSON storage
	data := map[string]interface{}{
		"name":    "test",
		"version": 1,
		"active":  true,
	}

	err := SetJSON("json_key", data)
	if err != nil {
		t.Fatalf("Failed to set JSON: %v", err)
	}

	var retrieved_data map[string]interface{}
	err = GetJSON("json_key", &retrieved_data)
	if err != nil {
		t.Fatalf("Failed to get JSON: %v", err)
	}

	if retrieved_data["name"] != "test" {
		t.Error("JSON storage failed for name field")
	}

	// Test deletion
	DeleteString(key)
	if HasString(key) {
		t.Error("Key should not exist after deletion")
	}
}

// TestCryptographicFunctions tests crypto functionality
func TestCryptographicFunctions(t *testing.T) {
	data := []byte("test data for hashing")

	// Test different hash functions
	blake3Hash := Blake3(data)
	sha256Hash := Sha256(data)
	keccakHash := Keccak256(data)
	mimcHash := Mimc(data)

	// All should produce 32-byte hashes
	if len(blake3Hash) != 32 {
		t.Errorf("BLAKE3 hash should be 32 bytes, got %d", len(blake3Hash))
	}
	if len(sha256Hash) != 32 {
		t.Errorf("SHA-256 hash should be 32 bytes, got %d", len(sha256Hash))
	}
	if len(keccakHash) != 32 {
		t.Errorf("Keccak-256 hash should be 32 bytes, got %d", len(keccakHash))
	}
	if len(mimcHash) != 32 {
		t.Errorf("MiMC hash should be 32 bytes, got %d", len(mimcHash))
	}

	// Test string hashing
	stringHash := HashString("test string")
	if len(stringHash) != 32 {
		t.Errorf("String hash should be 32 bytes, got %d", len(stringHash))
	}

	// Test hex conversion
	hexHash := HashStringHex("test string")
	if len(hexHash) != 64 {
		t.Errorf("Hex hash should be 64 characters, got %d", len(hexHash))
	}

	// Test bytes equality
	if !BytesEqual(blake3Hash, blake3Hash) {
		t.Error("Identical hashes should be equal")
	}

	if BytesEqual(blake3Hash, sha256Hash) {
		t.Error("Different hashes should not be equal")
	}
}

// TestCoinOperations tests coin-related functionality
func TestCoinOperations(t *testing.T) {
	// Test balance formatting
	balance := uint64(1000000) // 1 MAS
	formatted := FormatBalance(balance)
	if formatted != "1.000000 MAS" {
		t.Errorf("Balance formatting failed: expected '1.000000 MAS', got '%s'", formatted)
	}

	// Test amount parsing
	amount, err := ParseMassa("1.5 MAS")
	if err != nil {
		t.Fatalf("Failed to parse amount: %v", err)
	}
	if amount != 1500000 {
		t.Errorf("Amount parsing failed: expected 1500000, got %d", amount)
	}

	// Test unit conversions
	microMassa := ToMicroMassa(1.0)
	if microMassa != Massa {
		t.Errorf("Unit conversion failed: expected %d, got %d", Massa, microMassa)
	}

	massa := ToMassa(Massa)
	if massa != 1.0 {
		t.Errorf("Unit conversion failed: expected 1.0, got %f", massa)
	}
}

// TestEventGeneration tests event and logging functionality
func TestEventGeneration(t *testing.T) {
	// These should not panic
	GenerateEvent("Test event")
	Print("Test print")
	Printf("Test printf: %d", 42)

	// Test structured events
	err := CreateEvent("test_event", map[string]interface{}{
		"key":   "value",
		"count": 42,
	})
	if err != nil {
		t.Errorf("Failed to create structured event: %v", err)
	}

	// Test logging functions
	LogInfo("Test info message")
	LogError("Test error message")
	LogDebug("Test debug message")

	// Test custom events
	EmitCustomEvent("custom_test", map[string]interface{}{
		"test": true,
	})
}

// BenchmarkHashFunction benchmarks the BLAKE3 hash function
func BenchmarkHashFunction(b *testing.B) {
	data := []byte("benchmark data for hashing performance testing")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Blake3(data)
	}
}

// BenchmarkStorageOperations benchmarks storage operations
func BenchmarkStorageOperations(b *testing.B) {
	key := "benchmark_key"
	value := "benchmark_value"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SetString(key, value)
		GetString(key)
	}
}
