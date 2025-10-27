package wallet

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Wallet represents a complete Massa wallet with account management and transaction signing
type Wallet struct {
	// Account manager for handling accounts
	AccountManager *AccountManager

	// Transaction signer for signing operations
	TransactionSigner *TransactionSigner

	// Wallet file path for persistence
	walletPath string
}

// WalletConfig represents wallet configuration
type WalletConfig struct {
	// Path to store wallet data
	WalletPath string

	// Default network (mainnet, testnet, etc.)
	Network string
}

// WalletData represents serializable wallet data
type WalletData struct {
	// All accounts in the wallet
	Accounts []*Account `json:"accounts"`

	// Wallet metadata
	Version   string `json:"version"`
	Network   string `json:"network"`
	CreatedAt string `json:"createdAt"`
}

const (
	// Current wallet version
	WalletVersion = "1.0.0"

	// Default wallet filename
	DefaultWalletFilename = "massa_wallet.json"
)

// NewWallet creates a new wallet instance
func NewWallet(config *WalletConfig) *Wallet {
	accountManager := NewAccountManager()
	transactionSigner := NewTransactionSigner(accountManager)

	walletPath := config.WalletPath
	if walletPath == "" {
		walletPath = DefaultWalletFilename
	}

	return &Wallet{
		AccountManager:    accountManager,
		TransactionSigner: transactionSigner,
		walletPath:        walletPath,
	}
}

// LoadWallet loads a wallet from file
func LoadWallet(walletPath string) (*Wallet, error) {
	if walletPath == "" {
		walletPath = DefaultWalletFilename
	}

	// Check if wallet file exists
	if _, err := os.Stat(walletPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("wallet file does not exist: %s", walletPath)
	}

	// Read wallet file
	data, err := os.ReadFile(walletPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read wallet file: %w", err)
	}

	// Parse wallet data
	var walletData WalletData
	if err := json.Unmarshal(data, &walletData); err != nil {
		return nil, fmt.Errorf("failed to parse wallet file: %w", err)
	}

	// Create wallet instance
	wallet := &Wallet{
		AccountManager: NewAccountManager(),
		walletPath:     walletPath,
	}
	wallet.TransactionSigner = NewTransactionSigner(wallet.AccountManager)

	// Load accounts
	for _, account := range walletData.Accounts {
		wallet.AccountManager.accounts[account.Nickname] = account
	}

	return wallet, nil
}

// Save saves the wallet to file
func (w *Wallet) Save() error {
	// Prepare wallet data
	walletData := WalletData{
		Accounts:  w.AccountManager.ListAccounts(),
		Version:   WalletVersion,
		Network:   "mainnet", // Default network
		CreatedAt: "unknown", // Will be set if this is a new wallet
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(walletData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal wallet data: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(w.walletPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create wallet directory: %w", err)
	}

	// Write to file
	if err := os.WriteFile(w.walletPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write wallet file: %w", err)
	}

	return nil
}

// CreateAccount creates a new account in the wallet
func (w *Wallet) CreateAccount(nickname, password string) (*Account, error) {
	account, err := w.AccountManager.CreateAccount(nickname, password)
	if err != nil {
		return nil, err
	}

	// Save wallet after creating account
	if err := w.Save(); err != nil {
		// Remove the account if save fails
		w.AccountManager.DeleteAccount(nickname)
		return nil, fmt.Errorf("failed to save wallet after creating account: %w", err)
	}

	return account, nil
}

// ImportAccount imports an account into the wallet
func (w *Wallet) ImportAccount(nickname, privateKeyHex, password string) (*Account, error) {
	account, err := w.AccountManager.ImportAccount(nickname, privateKeyHex, password)
	if err != nil {
		return nil, err
	}

	// Save wallet after importing account
	if err := w.Save(); err != nil {
		// Remove the account if save fails
		w.AccountManager.DeleteAccount(nickname)
		return nil, fmt.Errorf("failed to save wallet after importing account: %w", err)
	}

	return account, nil
}

// ImportAccountFromPrivateKey imports an account from a Massa base58-encoded private key
func (w *Wallet) ImportAccountFromPrivateKey(nickname, privateKeyBase58, password string) (*Account, error) {
	account, err := w.AccountManager.ImportAccountFromPrivateKey(nickname, privateKeyBase58, password)
	if err != nil {
		return nil, err
	}

	// Save wallet after importing account
	if err := w.Save(); err != nil {
		// Remove the account if save fails
		w.AccountManager.DeleteAccount(nickname)
		return nil, fmt.Errorf("failed to save wallet after importing account: %w", err)
	}

	return account, nil
}

// DeleteAccount deletes an account from the wallet
func (w *Wallet) DeleteAccount(nickname string) error {
	if err := w.AccountManager.DeleteAccount(nickname); err != nil {
		return err
	}

	// Save wallet after deleting account
	return w.Save()
}

// TransferCoins creates and signs a coin transfer transaction
func (w *Wallet) TransferCoins(fromNickname, password, toAddress string, amount uint64) (*SignedOperation, error) {
	// Get current period for expiry
	expiryPeriod := GetCurrentPeriod()

	// Create transfer operation (with default fee)
	operation := CreateTransferOperation(toAddress, amount, 0, expiryPeriod)

	// Create signature request
	request := &SignatureRequest{
		Nickname:        fromNickname,
		Password:        password,
		Operation:       operation,
		AllowFeeEdition: true,
	}

	// Sign the operation
	return w.TransactionSigner.SignOperation(request)
}

// CallSmartContract creates and signs a smart contract call transaction
func (w *Wallet) CallSmartContract(fromNickname, password, contractAddress, functionName string, parameters []byte, maxGas, coins uint64) (*SignedOperation, error) {
	// Get current period for expiry
	expiryPeriod := GetCurrentPeriod()

	// Create call SC operation (with default fee)
	operation := CreateCallSCOperation(contractAddress, functionName, parameters, maxGas, coins, 0, expiryPeriod)

	// Create signature request
	request := &SignatureRequest{
		Nickname:        fromNickname,
		Password:        password,
		Operation:       operation,
		AllowFeeEdition: true,
	}

	// Sign the operation
	return w.TransactionSigner.SignOperation(request)
}

// DeploySmartContract creates and signs a smart contract deployment transaction
func (w *Wallet) DeploySmartContract(fromNickname, password string, bytecode []byte, maxGas, coins uint64, datastore map[string][]byte) (*SignedOperation, error) {
	// Get current period for expiry
	expiryPeriod := GetCurrentPeriod()

	// Create execute SC operation (with default fee)
	operation := CreateExecuteSCOperation(bytecode, maxGas, coins, 0, expiryPeriod, datastore)

	// Create signature request
	request := &SignatureRequest{
		Nickname:        fromNickname,
		Password:        password,
		Operation:       operation,
		AllowFeeEdition: true,
	}

	// Sign the operation
	return w.TransactionSigner.SignOperation(request)
}

// GetAccountInfo returns comprehensive account information
func (w *Wallet) GetAccountInfo(nickname string) (*Account, error) {
	return w.AccountManager.GetAccount(nickname)
}

// ListAccounts returns all accounts in the wallet
func (w *Wallet) ListAccounts() []*Account {
	return w.AccountManager.ListAccounts()
}

// ExportPrivateKey exports the private key of an account
func (w *Wallet) ExportPrivateKey(nickname, password string) (string, error) {
	return w.AccountManager.ExportPrivateKey(nickname, password)
}

// ExportPrivateKeyBase58 exports the private key of an account in Massa base58 format
func (w *Wallet) ExportPrivateKeyBase58(nickname, password string) (string, error) {
	return w.AccountManager.ExportPrivateKeyBase58(nickname, password)
}

// ChangeAccountPassword changes the password of an account
func (w *Wallet) ChangeAccountPassword(nickname, oldPassword, newPassword string) error {
	if err := w.AccountManager.ChangePassword(nickname, oldPassword, newPassword); err != nil {
		return err
	}

	// Save wallet after changing password
	return w.Save()
}

// ValidateAddress validates a Massa address
func (w *Wallet) ValidateAddress(address string) (*Address, error) {
	return ValidateAddress(address)
}

// SignArbitraryMessage signs an arbitrary message with an account
func (w *Wallet) SignArbitraryMessage(nickname, password string, message []byte) ([]byte, error) {
	return w.TransactionSigner.SignMessage(nickname, password, message)
}

// VerifySignature verifies a signature against a message and public key
func (w *Wallet) VerifySignature(message, signature, publicKey []byte) bool {
	return w.TransactionSigner.VerifySignature(message, signature, publicKey)
}

// GetWalletPath returns the path where the wallet is stored
func (w *Wallet) GetWalletPath() string {
	return w.walletPath
}

// SetWalletPath sets the path where the wallet should be stored
func (w *Wallet) SetWalletPath(path string) {
	w.walletPath = path
}

// BackupWallet creates a backup of the wallet file
func (w *Wallet) BackupWallet(backupPath string) error {
	// Read current wallet file
	data, err := os.ReadFile(w.walletPath)
	if err != nil {
		return fmt.Errorf("failed to read wallet file for backup: %w", err)
	}

	// Create backup directory if it doesn't exist
	dir := filepath.Dir(backupPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Write backup file
	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}

// WalletExists checks if a wallet file exists at the given path
func WalletExists(walletPath string) bool {
	if walletPath == "" {
		walletPath = DefaultWalletFilename
	}

	_, err := os.Stat(walletPath)
	return !os.IsNotExist(err)
}

// WalletStats returns statistics about the wallet
type WalletStats struct {
	AccountCount     int    `json:"accountCount"`
	TotalBalance     uint64 `json:"totalBalance"`
	WalletVersion    string `json:"walletVersion"`
	WalletPath       string `json:"walletPath"`
	HasEncryptedKeys bool   `json:"hasEncryptedKeys"`
}

// GetWalletStats returns statistics about the wallet
func (w *Wallet) GetWalletStats() *WalletStats {
	accounts := w.AccountManager.ListAccounts()
	totalBalance := uint64(0)
	hasEncryptedKeys := false

	for _, account := range accounts {
		totalBalance += account.Balance
		if account.KeyPair != nil {
			hasEncryptedKeys = true
		}
	}

	return &WalletStats{
		AccountCount:     len(accounts),
		TotalBalance:     totalBalance,
		WalletVersion:    WalletVersion,
		WalletPath:       w.walletPath,
		HasEncryptedKeys: hasEncryptedKeys,
	}
}
