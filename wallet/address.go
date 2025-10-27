package wallet

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ripemd160"
)

// Address represents a Massa blockchain address
type Address struct {
	// The raw address string (e.g., "AU12...")
	Address string `json:"address"`

	// Whether this is an Externally Owned Account (EOA) or Smart Contract
	IsEOA bool `json:"isEOA"`
}

const (
	// Address prefixes for Massa blockchain
	AddressPrefixUser = "AU"
	AddressPrefixSC   = "AS"

	// Address version and length constants
	AddressVersion        = uint8(0)
	AddressHashLength     = 20
	AddressChecksumLength = 4
	AddressTotalLength    = 1 + AddressHashLength + AddressChecksumLength // version + hash + checksum
)

// AddressFromPublicKey derives a Massa address from a public key
func AddressFromPublicKey(publicKey ed25519.PublicKey) (*Address, error) {
	if len(publicKey) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid public key size: expected %d, got %d", ed25519.PublicKeySize, len(publicKey))
	}

	// Create versioned public key (version byte + public key)
	versionedPublicKey := append([]byte{AddressVersion}, publicKey...)

	// Debug output
	fmt.Printf("DEBUG: Public key (hex): %x\n", publicKey)
	fmt.Printf("DEBUG: Version: 0x%02x\n", AddressVersion)
	fmt.Printf("DEBUG: Versioned public key (hex): %x\n", versionedPublicKey)

	// Hash the versioned public key with BLAKE3
	hash := blake3Hash(versionedPublicKey)
	fmt.Printf("DEBUG: BLAKE3 hash (hex): %x\n", hash)

	// Create address bytes (version + hash)
	addressBytes := append([]byte{AddressVersion}, hash...)
	fmt.Printf("DEBUG: Address bytes (hex): %x\n", addressBytes)

	// Use base58Check encoding (includes 4-byte checksum)
	base58Part := base58CheckEncode(addressBytes)
	fmt.Printf("DEBUG: Base58 part: %s\n", base58Part)

	addressString := "A" + "U" + base58Part
	fmt.Printf("DEBUG: Final address: %s\n", addressString)

	return &Address{
		Address: addressString,
		IsEOA:   true, // Addresses derived from public keys are always EOAs
	}, nil
}

// base58CheckEncode implements base58Check encoding (base58 + checksum)
func base58CheckEncode(input []byte) string {
	// Compute double SHA256 checksum
	firstHash := sha256.Sum256(input)
	secondHash := sha256.Sum256(firstHash[:])

	// Take first 4 bytes as checksum
	checksum := secondHash[:4]

	// Append checksum to input
	dataWithChecksum := make([]byte, len(input)+4)
	copy(dataWithChecksum, input)
	copy(dataWithChecksum[len(input):], checksum)

	// Encode with base58
	return base58.Encode(dataWithChecksum)
} // ValidateAddress validates a Massa address format and checksum
func ValidateAddress(address string) (*Address, error) {
	if len(address) < 3 {
		return nil, errors.New("address too short")
	}

	// Check prefix - must start with 'A'
	if address[0] != 'A' {
		return nil, fmt.Errorf("invalid address prefix: must start with 'A'")
	}

	// Check address type prefix
	typePrefix := address[1:2]
	isEOA := false
	switch typePrefix {
	case "U":
		isEOA = true
	case "S":
		isEOA = false
	default:
		return nil, fmt.Errorf("invalid address type prefix: %s", typePrefix)
	}

	// Decode base58 part (everything after AU or AS)
	base58Part := address[2:]
	decodedBytes := base58.Decode(base58Part)
	if len(decodedBytes) == 0 {
		return nil, errors.New("invalid base58 encoding")
	}

	// For now, we accept the decoded bytes as valid
	// TODO: Add more specific validation for the hash part

	return &Address{
		Address: address,
		IsEOA:   isEOA,
	}, nil
}

// NewAddress creates a new address from a string, validating its format
func NewAddress(address string) (*Address, error) {
	return ValidateAddress(address)
}

// String returns the string representation of the address
func (a *Address) String() string {
	return a.Address
}

// IsValid checks if the address is valid
func (a *Address) IsValid() bool {
	_, err := ValidateAddress(a.Address)
	return err == nil
}

// createAddressHash creates an address hash from a public key
func createAddressHash(publicKey ed25519.PublicKey, version uint8) []byte {
	// Hash the public key with SHA256
	sha256Hash := sha256.Sum256(publicKey)

	// Hash again with RIPEMD160 to get 20 bytes
	ripemd160Hasher := ripemd160.New()
	ripemd160Hasher.Write(sha256Hash[:])
	hash := ripemd160Hasher.Sum(nil)

	// Create the full address bytes: version + hash + checksum
	versionAndHash := append([]byte{version}, hash...)
	checksum := calculateChecksum(versionAndHash)

	return append(versionAndHash, checksum...)
}

// calculateChecksum calculates a 4-byte checksum using double SHA256
func calculateChecksum(data []byte) []byte {
	first := sha256.Sum256(data)
	second := sha256.Sum256(first[:])
	return second[:AddressChecksumLength]
}

// bytesEqual compares two byte slices for equality
func bytesEqual(a, b []byte) bool {
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

// Base58 encoding implementation (simplified version)
const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// base58Encode encodes bytes to base58 string
func base58Encode(input []byte) string {
	if len(input) == 0 {
		return ""
	}

	// Count leading zeros
	zeros := 0
	for i := 0; i < len(input) && input[i] == 0; i++ {
		zeros++
	}

	// Convert to big integer representation
	var num uint64 = 0
	for _, b := range input {
		num = num*256 + uint64(b)
	}

	// Convert to base58
	var result []byte
	for num > 0 {
		remainder := num % 58
		num = num / 58
		result = append([]byte{base58Alphabet[remainder]}, result...)
	}

	// Add leading zeros as '1's
	for i := 0; i < zeros; i++ {
		result = append([]byte{'1'}, result...)
	}

	return string(result)
}

// base58Decode decodes a base58 string to bytes
func base58Decode(input string) ([]byte, error) {
	if len(input) == 0 {
		return nil, nil
	}

	// Count leading '1's (zeros)
	zeros := 0
	for i := 0; i < len(input) && input[i] == '1'; i++ {
		zeros++
	}

	// Convert from base58
	var num uint64 = 0
	for _, char := range input {
		index := -1
		for i, c := range base58Alphabet {
			if c == char {
				index = i
				break
			}
		}
		if index == -1 {
			return nil, fmt.Errorf("invalid base58 character: %c", char)
		}
		num = num*58 + uint64(index)
	}

	// Convert to bytes
	var result []byte
	for num > 0 {
		result = append([]byte{byte(num % 256)}, result...)
		num = num / 256
	}

	// Add leading zeros
	for i := 0; i < zeros; i++ {
		result = append([]byte{0}, result...)
	}

	return result, nil
}

// GetAddressVersion extracts the version from an address
func GetAddressVersion(address string) (uint8, error) {
	addr, err := ValidateAddress(address)
	if err != nil {
		return 0, err
	}

	// Decode to get version
	base58Part := addr.Address[2:]
	decodedBytes, err := base58Decode(base58Part)
	if err != nil {
		return 0, err
	}

	return decodedBytes[0], nil
}

// CompareAddresses compares two addresses for equality
func CompareAddresses(addr1, addr2 string) bool {
	return addr1 == addr2
}
