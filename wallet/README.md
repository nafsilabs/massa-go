# Massa Go Wallet

A comprehensive Go library for managing Massa blockchain accounts, including key generation, transaction signing, and wallet operations.

## 🚀 Features

- **Account Generation**: Create new Ed25519 key pairs for Massa blockchain
- **Key Management**: Secure password-based encryption for private keys
- **Transaction Signing**: Sign transfers, smart contract calls, and deployments
- **Address Derivation**: Generate and validate Massa blockchain addresses
- **Wallet Persistence**: Save and load wallet data from JSON files
- **Multi-Account Support**: Manage multiple accounts in a single wallet
- **Comprehensive Security**: AES-GCM encryption with PBKDF2 key derivation

## 📦 Installation

```bash
go get github.com/nafsilabs/massa-go/wallet
```

## 🔧 Dependencies

- `golang.org/x/crypto` - Cryptographic functions
- Go 1.19+ - Required for modern Go features

## 📖 Quick Start

### Basic Key Pair Operations

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/nafsilabs/massa-go/wallet"
)

func main() {
    // Generate a new key pair
    keyPair, err := wallet.GenerateKeyPair()
    if err != nil {
        log.Fatal(err)
    }
    
    // Get hex representation
    privateKeyHex, publicKeyHex := keyPair.ToHex()
    fmt.Printf("Private Key: %s\n", privateKeyHex)
    fmt.Printf("Public Key:  %s\n", publicKeyHex)
    
    // Generate address from public key
    address, err := wallet.AddressFromPublicKey(keyPair.PublicKey)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Address: %s\n", address.String())
    fmt.Printf("Is EOA: %t\n", address.IsEOA)
}
```

### Wallet Management

```go
// Create a new wallet
config := &wallet.WalletConfig{
    WalletPath: "my_wallet.json",
    Network:    "mainnet",
}

w := wallet.NewWallet(config)

// Create an account
account, err := w.CreateAccount("my-account", "secure-password")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created account: %s\n", account.Address.String())

// Save wallet to file
err = w.Save()
if err != nil {
    log.Fatal(err)
}
```

### Transaction Signing

```go
// Sign a transfer transaction
recipientAddress := "AU1234567890123456789012345678901234567890"
amount := uint64(1000000000) // 1 MAS in nanoMAS

signedTx, err := w.TransferCoins("my-account", "secure-password", recipientAddress, amount)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Signed transaction: %s\n", signedTx.Signature)

// Sign a smart contract call
contractAddress := "AS1234567890123456789012345678901234567890"
functionName := "transfer"
parameters := []byte(`{"to": "AU...", "amount": 1000}`)
maxGas := uint64(1000000)
coins := uint64(0)

signedCall, err := w.CallSmartContract(
    "my-account", 
    "secure-password", 
    contractAddress, 
    functionName, 
    parameters, 
    maxGas, 
    coins,
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Smart contract call signed: %s\n", signedCall.Signature)
```

## 🏗️ Architecture

### Core Components

1. **KeyPair Management** (`keypair.go`)
   - Ed25519 key generation
   - AES-GCM encryption with PBKDF2
   - Secure key storage and retrieval

2. **Address Operations** (`address.go`)
   - Address generation from public keys
   - Address validation and formatting
   - Support for EOA and Smart Contract addresses

3. **Account Management** (`account.go`)
   - Multi-account wallet support
   - Account metadata and status tracking
   - Password-protected account operations

4. **Transaction Signing** (`transaction.go`)
   - Support for all Massa transaction types
   - Binary serialization for blockchain compatibility
   - Digital signature generation and verification

5. **Wallet Operations** (`wallet.go`)
   - High-level wallet interface
   - Persistent storage management
   - Comprehensive wallet statistics

### Security Features

- **Password-Based Encryption**: Private keys are encrypted using AES-GCM
- **Key Derivation**: PBKDF2 with 100,000 iterations for password-based key derivation
- **Secure Random Generation**: Cryptographically secure random number generation
- **Memory Safety**: Secure handling of sensitive cryptographic material

## 🔐 Cryptographic Details

### Key Generation
- **Algorithm**: Ed25519 (Curve25519 with EdDSA)
- **Key Size**: 256-bit private keys, 256-bit public keys
- **Random Source**: `crypto/rand` for secure entropy

### Encryption
- **Algorithm**: AES-256-GCM
- **Key Derivation**: PBKDF2-SHA256 with 100,000 iterations
- **Salt**: 32-byte random salt per key pair
- **Nonce**: 12-byte random nonce per encryption

### Address Generation
- **Hash Functions**: SHA-256 → RIPEMD-160
- **Checksum**: Double SHA-256 (first 4 bytes)
- **Encoding**: Base58 with Massa-specific prefixes
- **Prefixes**: `AU` for EOAs, `AS` for Smart Contracts

## 📋 API Reference

### KeyPair Operations

```go
// Generate new key pair
keyPair, err := wallet.GenerateKeyPair()

// Encrypt with password
encrypted, err := wallet.EncryptKeyPair(keyPair, password)

// Decrypt with password
decrypted, err := wallet.DecryptKeyPair(encrypted, password)

// Sign message
signature := keyPair.Sign(message)

// Verify signature
isValid := keyPair.Verify(message, signature)
```

### Address Operations

```go
// Generate address from public key
address, err := wallet.AddressFromPublicKey(publicKey)

// Validate address format
validatedAddr, err := wallet.ValidateAddress(addressString)

// Compare addresses
isEqual := wallet.CompareAddresses(addr1, addr2)
```

### Account Management

```go
// Create account manager
am := wallet.NewAccountManager()

// Create account
account, err := am.CreateAccount(nickname, password)

// Import account from private key
account, err := am.ImportAccount(nickname, privateKeyHex, password)

// Unlock account
keyPair, err := am.UnlockAccount(nickname, password)

// List accounts
accounts := am.ListAccounts()
```

### Transaction Signing

```go
// Create transaction signer
ts := wallet.NewTransactionSigner(accountManager)

// Sign operation
signedOp, err := ts.SignOperation(request)

// Sign arbitrary message
signature, err := ts.SignMessage(nickname, password, message)

// Verify signature
isValid := ts.VerifySignature(message, signature, publicKey)
```

### Wallet Operations

```go
// Create wallet
config := &wallet.WalletConfig{WalletPath: "wallet.json"}
w := wallet.NewWallet(config)

// Load wallet from file
w, err := wallet.LoadWallet("wallet.json")

// Create account
account, err := w.CreateAccount(nickname, password)

// Transfer coins
signedTx, err := w.TransferCoins(from, password, to, amount)

// Call smart contract
signedCall, err := w.CallSmartContract(from, password, contract, function, params, gas, coins)

// Deploy smart contract
signedDeploy, err := w.DeploySmartContract(from, password, bytecode, gas, coins, datastore)

// Save wallet
err = w.Save()
```

## 🧪 Testing

Run the comprehensive test suite:

```bash
go test -v
```

Run benchmarks:

```bash
go test -bench=.
```

## 🔧 Configuration

### Wallet Configuration

```go
type WalletConfig struct {
    // Path to store wallet data
    WalletPath string
    
    // Network (mainnet, testnet, etc.)
    Network string
}
```

### Security Parameters

- **PBKDF2 Iterations**: 100,000 (configurable via constants)
- **Salt Length**: 32 bytes
- **AES Key Length**: 32 bytes (AES-256)
- **Nonce Length**: 12 bytes (GCM standard)

## 🚨 Security Considerations

1. **Password Security**: Use strong, unique passwords for each account
2. **Private Key Storage**: Private keys are encrypted and never stored in plaintext
3. **Memory Management**: Sensitive data is handled securely in memory
4. **File Permissions**: Wallet files are created with restrictive permissions (0600)
5. **Backup Safety**: Always backup your wallet files securely

## 📝 Transaction Types

The wallet supports all Massa transaction types:

- **Transfer**: Coin transfers between addresses
- **CallSC**: Smart contract function calls
- **ExecuteSC**: Smart contract deployment
- **BuyRolls**: Purchase validator rolls
- **SellRolls**: Sell validator rolls

## 🔄 Serialization Format

Transactions are serialized in Massa-compatible binary format:

1. **Operation Type** (1 byte)
2. **Fee** (8 bytes, little-endian)
3. **Expiry Period** (8 bytes, little-endian)
4. **Content** (variable length, type-specific)

## 💡 Best Practices

1. **Password Management**: Use a secure password manager
2. **Regular Backups**: Backup wallet files regularly
3. **Version Control**: Never commit wallet files to version control
4. **Testing**: Test with small amounts first
5. **Network Awareness**: Use appropriate networks for testing vs production

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.

## 🔗 Related Projects

- [massa-go/sc](../sc) - Smart Contract SDK
- [massa-go/client](../client) - Blockchain Client
- [Massa](https://massa.net/) - Massa Blockchain

---

Built with ❤️ for the Massa ecosystem