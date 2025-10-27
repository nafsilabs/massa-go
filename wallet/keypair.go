package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/pbkdf2"
)

// KeyPair represents a cryptographic key pair for Massa blockchain
type KeyPair struct {
	// Nonce used by the AES-GCM algorithm to protect the private key
	Nonce string `json:"nonce"`

	// Private key (encrypted with password)
	PrivateKey string `json:"privateKey"`

	// Public key (unencrypted)
	PublicKey string `json:"publicKey"`

	// Salt used by PBKDF2 for key derivation
	Salt string `json:"salt"`
}

// KeyPairRaw represents unencrypted key pair
type KeyPairRaw struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

const (
	// PBKDF2 parameters (Massa standard)
	pbkdf2Iterations = 600000 // Massa standard: 600,000 iterations
	pbkdf2KeyLength  = 32     // 32-byte derived key
	pbkdf2SaltLength = 16     // Massa standard: 16-byte salt

	// AES-GCM parameters
	aesGCMNonceLength = 12
)

// GenerateKeyPair generates a new Ed25519 key pair for Massa blockchain
func GenerateKeyPair() (*KeyPairRaw, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Ed25519 key pair: %w", err)
	}

	return &KeyPairRaw{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}

// EncryptKeyPair encrypts a key pair with a password using AES-GCM
func EncryptKeyPair(keyPair *KeyPairRaw, password string) (*KeyPair, error) {
	if len(keyPair.PrivateKey) != ed25519.PrivateKeySize {
		return nil, errors.New("invalid private key size")
	}

	// Generate random salt for PBKDF2
	salt := make([]byte, pbkdf2SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive encryption key using PBKDF2
	derivedKey := pbkdf2.Key([]byte(password), salt, pbkdf2Iterations, pbkdf2KeyLength, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, aesGCMNonceLength)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt private key
	encryptedPrivateKey := gcm.Seal(nil, nonce, keyPair.PrivateKey, nil)

	return &KeyPair{
		Nonce:      base64.StdEncoding.EncodeToString(nonce),
		PrivateKey: base64.StdEncoding.EncodeToString(encryptedPrivateKey),
		PublicKey:  base64.StdEncoding.EncodeToString(keyPair.PublicKey),
		Salt:       base64.StdEncoding.EncodeToString(salt),
	}, nil
}

// DecryptKeyPair decrypts an encrypted key pair using the password
func DecryptKeyPair(encryptedKeyPair *KeyPair, password string) (*KeyPairRaw, error) {
	// Decode base64 fields
	salt, err := base64.StdEncoding.DecodeString(encryptedKeyPair.Salt)
	if err != nil {
		return nil, fmt.Errorf("failed to decode salt: %w", err)
	}

	nonce, err := base64.StdEncoding.DecodeString(encryptedKeyPair.Nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to decode nonce: %w", err)
	}

	encryptedPrivateKey, err := base64.StdEncoding.DecodeString(encryptedKeyPair.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encrypted private key: %w", err)
	}

	publicKey, err := base64.StdEncoding.DecodeString(encryptedKeyPair.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	// Derive decryption key using PBKDF2
	derivedKey := pbkdf2.Key([]byte(password), salt, pbkdf2Iterations, pbkdf2KeyLength, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt private key
	privateKey, err := gcm.Open(nil, nonce, encryptedPrivateKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key (wrong password?): %w", err)
	}

	// Validate key pair
	if len(privateKey) != ed25519.PrivateKeySize {
		return nil, errors.New("invalid decrypted private key size")
	}

	if len(publicKey) != ed25519.PublicKeySize {
		return nil, errors.New("invalid public key size")
	}

	return &KeyPairRaw{
		PrivateKey: ed25519.PrivateKey(privateKey),
		PublicKey:  ed25519.PublicKey(publicKey),
	}, nil
}

// ToHex returns the hexadecimal representation of the key pair
func (kp *KeyPairRaw) ToHex() (privateKeyHex, publicKeyHex string) {
	return hex.EncodeToString(kp.PrivateKey), hex.EncodeToString(kp.PublicKey)
}

// FromHex creates a KeyPairRaw from hexadecimal strings
func FromHex(privateKeyHex, publicKeyHex string) (*KeyPairRaw, error) {
	privateKey, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key hex: %w", err)
	}

	// Handle both seed (32 bytes) and full private key (64 bytes)
	var actualPrivateKey ed25519.PrivateKey
	if len(privateKey) == ed25519.SeedSize {
		// This is a seed, generate the full private key
		actualPrivateKey = ed25519.NewKeyFromSeed(privateKey)
	} else if len(privateKey) == ed25519.PrivateKeySize {
		// This is already a full private key
		actualPrivateKey = ed25519.PrivateKey(privateKey)
	} else {
		return nil, fmt.Errorf("invalid private key size: expected %d (seed) or %d (full key), got %d", ed25519.SeedSize, ed25519.PrivateKeySize, len(privateKey))
	}

	// Handle public key
	var actualPublicKey ed25519.PublicKey
	if publicKeyHex == "" {
		// Derive public key from private key
		actualPublicKey = actualPrivateKey.Public().(ed25519.PublicKey)
	} else {
		publicKey, err := hex.DecodeString(publicKeyHex)
		if err != nil {
			return nil, fmt.Errorf("invalid public key hex: %w", err)
		}

		if len(publicKey) != ed25519.PublicKeySize {
			return nil, fmt.Errorf("invalid public key size: expected %d, got %d", ed25519.PublicKeySize, len(publicKey))
		}
		actualPublicKey = ed25519.PublicKey(publicKey)
	}

	return &KeyPairRaw{
		PrivateKey: actualPrivateKey,
		PublicKey:  actualPublicKey,
	}, nil
}

// FromBase58PrivateKey creates a KeyPairRaw from a Massa base58-encoded private key
func FromBase58PrivateKey(privateKeyBase58 string) (*KeyPairRaw, error) {
	// Check if it starts with 'S' prefix
	if len(privateKeyBase58) == 0 || privateKeyBase58[0] != 'S' {
		return nil, fmt.Errorf("invalid Massa private key format: must start with 'S'")
	}

	// Remove the 'S' prefix and decode base58
	base58Part := privateKeyBase58[1:]
	privateKeyBytes := base58.Decode(base58Part)

	if len(privateKeyBytes) == 0 {
		return nil, fmt.Errorf("invalid base58 encoding")
	}

	// Handle different private key formats
	var seed []byte

	if len(privateKeyBytes) == 33 {
		// New format: version (1 byte) + 32-byte seed
		if privateKeyBytes[0] != 0x00 {
			return nil, fmt.Errorf("invalid private key version: expected 0x00, got 0x%02x", privateKeyBytes[0])
		}
		seed = privateKeyBytes[1:]
	} else if len(privateKeyBytes) == 38 {
		// Legacy format: version + type + 32-byte seed + checksum
		seed = privateKeyBytes[2:34]
	} else if len(privateKeyBytes) == 37 {
		// For now, use position 1-33 as before
		seed = privateKeyBytes[1:33]
	} else if len(privateKeyBytes) == 32 {
		// Raw seed format
		seed = privateKeyBytes
	} else {
		return nil, fmt.Errorf("invalid private key length: expected 33, 38, 37, or 32 bytes, got %d", len(privateKeyBytes))
	}

	// Create ed25519 private key from seed
	privateKey := ed25519.NewKeyFromSeed(seed)

	// Derive public key
	publicKey := privateKey.Public().(ed25519.PublicKey)
	return &KeyPairRaw{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}

// ToBase58 returns the base58-encoded representation of the private key (Massa format)
func (kp *KeyPairRaw) ToBase58() string {
	// Extract the seed from ed25519 private key (first 32 bytes)
	seed := kp.PrivateKey.Seed()

	// Create Massa private key format: version (0x00) + 32-byte seed
	versionedKey := make([]byte, 33)
	versionedKey[0] = 0x00       // Version byte
	copy(versionedKey[1:], seed) // Copy 32-byte seed

	// Return in Massa format: 'S' + base58(versioned_key)
	return "S" + base58.Encode(versionedKey)
}

// ToBase58PublicKey returns the public key in Massa base58 format
func (kp *KeyPairRaw) ToBase58PublicKey() string {
	// Create Massa public key format: version (0x00) + 32-byte public key
	versionedKey := make([]byte, 33)
	versionedKey[0] = 0x00               // Version byte
	copy(versionedKey[1:], kp.PublicKey) // Copy 32-byte public key

	// Return in Massa format: 'P' + base58(versioned_key)
	return "P" + base58.Encode(versionedKey)
}

// Sign signs a message using the private key
func (kp *KeyPairRaw) Sign(message []byte) []byte {
	return ed25519.Sign(kp.PrivateKey, message)
}

// Verify verifies a signature using the public key
func (kp *KeyPairRaw) Verify(message, signature []byte) bool {
	return ed25519.Verify(kp.PublicKey, message, signature)
}

// PublicKeyFromPrivate derives the public key from a private key
func PublicKeyFromPrivate(privateKey ed25519.PrivateKey) ed25519.PublicKey {
	return privateKey.Public().(ed25519.PublicKey)
}

// VersionedSignatureBytes returns the Ed25519 signature prefixed with a version byte.
// It follows the convention used across Massa libraries where signatures are
// encoded with a single version byte (0x00) followed by the raw signature bytes.
// Example output shape: [0x00 || 64-byte ed25519 signature]
func (kp *KeyPairRaw) VersionedSignatureBytes(message []byte) []byte {
	sig := kp.Sign(message)
	out := make([]byte, 1+len(sig))
	out[0] = 0x00 // version
	copy(out[1:], sig)
	return out
}

// VersionedSignatureBase58 returns a base58-encoded, prefixed signature string
// suitable for human-friendly transport. The returned string is prefixed with
// 'SIG' to make the format self-describing (similar to other Massa prefixed
// encodings like 'S' for private keys and 'P' for public keys).
func (kp *KeyPairRaw) VersionedSignatureBase58(message []byte) string {
	versioned := kp.VersionedSignatureBytes(message)
	return "SIG" + base58.Encode(versioned)
}

// ParseVersionedSignatureBytes parses a versioned signature byte slice
// Expected format: [version (1 byte) || signature (64 bytes)].
// Returns version byte and signature bytes.
func ParseVersionedSignatureBytes(data []byte) (byte, []byte, error) {
	if len(data) < 1 {
		return 0x00, nil, fmt.Errorf("versioned signature too short: %d", len(data))
	}
	version := data[0]
	sig := data[1:]
	if len(sig) != ed25519.SignatureSize {
		return 0x00, nil, fmt.Errorf("invalid signature length: %d", len(sig))
	}
	return version, sig, nil
}

// ParseVersionedSignatureBase58 parses a base58-encoded versioned signature string.
// Accepts strings prefixed with "SIG" or raw base58 string. Returns version byte and signature bytes.
func ParseVersionedSignatureBase58(s string) (byte, []byte, error) {
	if s == "" {
		return 0x00, nil, fmt.Errorf("empty signature string")
	}
	// Allow optional "SIG" prefix
	if len(s) > 3 && s[:3] == "SIG" {
		s = s[3:]
	}
	decoded := base58.Decode(s)
	if len(decoded) == 0 {
		return 0x00, nil, fmt.Errorf("invalid base58 signature")
	}
	return ParseVersionedSignatureBytes(decoded)
}

// VerifyVersionedSignature verifies a versioned signature bytes against the message and public key.
// It ignores the version byte for verification purposes but ensures signature length is valid.
func VerifyVersionedSignature(message, versionedSig []byte, publicKey ed25519.PublicKey) bool {
	_, sig, err := ParseVersionedSignatureBytes(versionedSig)
	if err != nil {
		return false
	}
	return ed25519.Verify(publicKey, message, sig)
}

// VerifyVersionedSignatureBase58 verifies a base58-encoded versioned signature string.
func VerifyVersionedSignatureBase58(message []byte, s string, publicKey ed25519.PublicKey) bool {
	_, sig, err := ParseVersionedSignatureBase58(s)
	if err != nil {
		return false
	}
	return ed25519.Verify(publicKey, message, sig)
}
