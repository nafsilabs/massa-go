package sc

import "fmt"

// Cryptographic functions for Massa smart contracts

// Hash functions

// Blake3 computes the BLAKE3 hash of the given data
func Blake3(data []byte) []byte {
	if isWasmRuntime() {
		dataPtr := newWasmBytes(data)
		resultPtr := wasmBlake3(dataPtr)
		return wasmBytesFromPtr(resultPtr)
	}
	// In non-WASM mode, return a mock hash for testing
	mockHash := make([]byte, 32) // BLAKE3 produces 32-byte hashes
	for i := range mockHash {
		mockHash[i] = byte(i % 256)
	}
	return mockHash
}

// Sha256 computes the SHA-256 hash of the given data
func Sha256(data []byte) []byte {
	if isWasmRuntime() {
		dataPtr := newWasmBytes(data)
		resultPtr := wasmSha256(dataPtr)
		return wasmBytesFromPtr(resultPtr)
	}
	// In non-WASM mode, return a mock hash for testing
	mockHash := make([]byte, 32) // SHA-256 produces 32-byte hashes
	for i := range mockHash {
		mockHash[i] = byte((i + 1) % 256)
	}
	return mockHash
}

// Keccak256 computes the Keccak-256 hash of the given data
func Keccak256(data []byte) []byte {
	if isWasmRuntime() {
		dataPtr := newWasmBytes(data)
		resultPtr := wasmKeccak256(dataPtr)
		return wasmBytesFromPtr(resultPtr)
	}
	// In non-WASM mode, return a mock hash for testing
	mockHash := make([]byte, 32) // Keccak-256 produces 32-byte hashes
	for i := range mockHash {
		mockHash[i] = byte((i + 2) % 256)
	}
	return mockHash
}

// Mimc computes the MiMC hash of the given data
func Mimc(data []byte) []byte {
	if isWasmRuntime() {
		dataPtr := newWasmBytes(data)
		resultPtr := wasmMimc(dataPtr)
		return wasmBytesFromPtr(resultPtr)
	}
	// In non-WASM mode, return a mock hash for testing
	mockHash := make([]byte, 32) // MiMC produces 32-byte hashes
	for i := range mockHash {
		mockHash[i] = byte((i + 3) % 256)
	}
	return mockHash
}

// Signature verification functions

// IsSignatureValid verifies a signature against data and a public key
func IsSignatureValid(data, signature, publicKey string) bool {
	if isWasmRuntime() {
		dataPtr := newWasmString(data)
		signaturePtr := newWasmString(signature)
		publicKeyPtr := newWasmString(publicKey)
		result := wasmSignatureVerify(dataPtr, signaturePtr, publicKeyPtr)
		return result != 0
	}
	// In non-WASM mode, return true for testing (always valid)
	return true
}

// IsEvmSignatureValid verifies an Ethereum-style signature
func IsEvmSignatureValid(data, signature, publicKey []byte) bool {
	if isWasmRuntime() {
		dataPtr := newWasmBytes(data)
		signaturePtr := newWasmBytes(signature)
		publicKeyPtr := newWasmBytes(publicKey)
		result := wasmEvmSignatureVerify(dataPtr, signaturePtr, publicKeyPtr)
		return result != 0
	}
	// In non-WASM mode, return true for testing (always valid)
	return true
}

// EvmGetAddressFromPubkey derives an Ethereum address from a public key
func EvmGetAddressFromPubkey(publicKey []byte) []byte {
	if isWasmRuntime() {
		publicKeyPtr := newWasmBytes(publicKey)
		resultPtr := wasmEvmGetAddressFromPubkey(publicKeyPtr)
		return wasmBytesFromPtr(resultPtr)
	}
	// In non-WASM mode, return a mock address for testing
	mockAddress := make([]byte, 20) // Ethereum addresses are 20 bytes
	for i := range mockAddress {
		mockAddress[i] = byte(i % 256)
	}
	return mockAddress
}

// EvmGetPubkeyFromSignature recovers a public key from an Ethereum signature
func EvmGetPubkeyFromSignature(hash, signature []byte) []byte {
	if isWasmRuntime() {
		hashPtr := newWasmBytes(hash)
		signaturePtr := newWasmBytes(signature)
		resultPtr := wasmEvmGetPubkeyFromSignature(hashPtr, signaturePtr)
		return wasmBytesFromPtr(resultPtr)
	}
	// In non-WASM mode, return a mock public key for testing
	mockPubkey := make([]byte, 64) // Uncompressed public keys are 64 bytes
	for i := range mockPubkey {
		mockPubkey[i] = byte((i + 1) % 256)
	}
	return mockPubkey
}

// Convenience functions for common cryptographic operations

// HashString computes the BLAKE3 hash of a string
func HashString(s string) []byte {
	return Blake3([]byte(s))
}

// HashStringHex computes the BLAKE3 hash of a string and returns it as hex
func HashStringHex(s string) string {
	hash := HashString(s)
	return BytesToHex(hash)
}

// Sha256String computes the SHA-256 hash of a string
func Sha256String(s string) []byte {
	return Sha256([]byte(s))
}

// Sha256StringHex computes the SHA-256 hash of a string and returns it as hex
func Sha256StringHex(s string) string {
	hash := Sha256String(s)
	return BytesToHex(hash)
}

// Keccak256String computes the Keccak-256 hash of a string
func Keccak256String(s string) []byte {
	return Keccak256([]byte(s))
}

// Keccak256StringHex computes the Keccak-256 hash of a string and returns it as hex
func Keccak256StringHex(s string) string {
	hash := Keccak256String(s)
	return BytesToHex(hash)
}

// Utility functions for encoding/decoding

// BytesToHex converts bytes to a hexadecimal string
func BytesToHex(data []byte) string {
	const hexChars = "0123456789abcdef"
	result := make([]byte, len(data)*2)
	for i, b := range data {
		result[i*2] = hexChars[b>>4]
		result[i*2+1] = hexChars[b&0x0f]
	}
	return string(result)
}

// HexToBytes converts a hexadecimal string to bytes
func HexToBytes(hexStr string) ([]byte, error) {
	// Remove 0x prefix if present
	if len(hexStr) >= 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}

	// Ensure even length
	if len(hexStr)%2 != 0 {
		return nil, fmt.Errorf("hex string must have even length")
	}

	result := make([]byte, len(hexStr)/2)
	for i := 0; i < len(result); i++ {
		high := hexCharToNibble(hexStr[i*2])
		low := hexCharToNibble(hexStr[i*2+1])
		if high < 0 || low < 0 {
			return nil, fmt.Errorf("invalid hex character")
		}
		result[i] = byte(high<<4 | low)
	}

	return result, nil
}

// hexCharToNibble converts a hex character to its numeric value
func hexCharToNibble(c byte) int {
	switch {
	case '0' <= c && c <= '9':
		return int(c - '0')
	case 'a' <= c && c <= 'f':
		return int(c - 'a' + 10)
	case 'A' <= c && c <= 'F':
		return int(c - 'A' + 10)
	default:
		return -1
	}
}

// Cryptographic utilities for smart contracts

// GenerateKeyPair generates a new key pair (mock implementation for testing)
type KeyPair struct {
	PublicKey  []byte
	PrivateKey []byte
}

// GenerateRandomBytes generates random bytes using the unsafe random function
// WARNING: This is not cryptographically secure!
func GenerateRandomBytes(length int) []byte {
	result := make([]byte, length)

	// Use the unsafe random function as a seed
	seed := UnsafeRandom()

	// Simple linear congruential generator for demo purposes
	// This is NOT secure and should not be used for real cryptographic purposes
	state := uint64(seed)
	for i := 0; i < length; i++ {
		state = (state*1103515245 + 12345) % (1 << 32)
		result[i] = byte(state & 0xff)
	}

	return result
}

// HashChain computes a chain of hashes
func HashChain(data []byte, iterations int) []byte {
	result := data
	for i := 0; i < iterations; i++ {
		result = Blake3(result)
	}
	return result
}

// VerifyHashChain verifies a hash chain
func VerifyHashChain(originalData, finalHash []byte, iterations int) bool {
	computed := HashChain(originalData, iterations)
	return BytesEqual(computed, finalHash)
}

// BytesEqual compares two byte slices for equality
func BytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// CombineHashes combines multiple hashes into a single hash
func CombineHashes(hashes ...[]byte) []byte {
	var combined []byte
	for _, hash := range hashes {
		combined = append(combined, hash...)
	}
	return Blake3(combined)
}

// Merkle tree utilities (simplified implementation)

// ComputeMerkleRoot computes the Merkle root of a list of data items
func ComputeMerkleRoot(data [][]byte) []byte {
	if len(data) == 0 {
		return Blake3([]byte{})
	}

	if len(data) == 1 {
		return Blake3(data[0])
	}

	// Hash all leaf nodes
	hashes := make([][]byte, len(data))
	for i, item := range data {
		hashes[i] = Blake3(item)
	}

	// Build the tree bottom-up
	for len(hashes) > 1 {
		var nextLevel [][]byte

		for i := 0; i < len(hashes); i += 2 {
			if i+1 < len(hashes) {
				// Combine two hashes
				combined := append(hashes[i], hashes[i+1]...)
				nextLevel = append(nextLevel, Blake3(combined))
			} else {
				// Odd number of hashes, promote the last one
				nextLevel = append(nextLevel, hashes[i])
			}
		}

		hashes = nextLevel
	}

	return hashes[0]
}

// VerifyMerkleProof verifies a Merkle proof
func VerifyMerkleProof(data []byte, proof [][]byte, root []byte, index int) bool {
	hash := Blake3(data)

	for _, proofHash := range proof {
		if index%2 == 0 {
			// Current node is left child
			combined := append(hash, proofHash...)
			hash = Blake3(combined)
		} else {
			// Current node is right child
			combined := append(proofHash, hash...)
			hash = Blake3(combined)
		}
		index /= 2
	}

	return BytesEqual(hash, root)
}
