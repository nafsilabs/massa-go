package wallet

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Account represents a Massa wallet account
type Account struct {
	// Account's nickname/name
	Nickname string `json:"nickname"`

	// Account's Massa address
	Address *Address `json:"address"`

	// Encrypted key pair
	KeyPair *KeyPair `json:"keyPair"`

	// Account balance (in nanoMAS)
	Balance uint64 `json:"balance"`

	// Candidate balance (in nanoMAS)
	CandidateBalance uint64 `json:"candidateBalance"`

	// Account status
	Status AccountStatus `json:"status"`

	// Creation timestamp
	CreatedAt time.Time `json:"createdAt"`

	// Last update timestamp
	UpdatedAt time.Time `json:"updatedAt"`
}

// AccountStatus represents the status of an account
type AccountStatus string

const (
	AccountStatusOK        AccountStatus = "ok"
	AccountStatusCorrupted AccountStatus = "corrupted"
	AccountStatusLocked    AccountStatus = "locked"
)

// AccountManager manages wallet accounts
type AccountManager struct {
	accounts map[string]*Account
}

// NewAccountManager creates a new account manager
func NewAccountManager() *AccountManager {
	return &AccountManager{
		accounts: make(map[string]*Account),
	}
}

// CreateAccount creates a new account with the given nickname and password
func (am *AccountManager) CreateAccount(nickname, password string) (*Account, error) {
	if nickname == "" {
		return nil, errors.New("nickname cannot be empty")
	}

	if password == "" {
		return nil, errors.New("password cannot be empty")
	}

	// Check if account already exists
	if _, exists := am.accounts[nickname]; exists {
		return nil, fmt.Errorf("account with nickname '%s' already exists", nickname)
	}

	// Generate new key pair
	keyPairRaw, err := GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %w", err)
	}

	// Encrypt the key pair
	encryptedKeyPair, err := EncryptKeyPair(keyPairRaw, password)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt key pair: %w", err)
	}

	// Derive address from public key
	address, err := AddressFromPublicKey(keyPairRaw.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to derive address: %w", err)
	}

	// Create account
	account := &Account{
		Nickname:         nickname,
		Address:          address,
		KeyPair:          encryptedKeyPair,
		Balance:          0,
		CandidateBalance: 0,
		Status:           AccountStatusOK,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Store account
	am.accounts[nickname] = account

	return account, nil
}

// ImportAccount imports an account from a private key
func (am *AccountManager) ImportAccount(nickname, privateKeyHex, password string) (*Account, error) {
	if nickname == "" {
		return nil, errors.New("nickname cannot be empty")
	}

	if password == "" {
		return nil, errors.New("password cannot be empty")
	}

	if privateKeyHex == "" {
		return nil, errors.New("private key cannot be empty")
	}

	// Check if account already exists
	if _, exists := am.accounts[nickname]; exists {
		return nil, fmt.Errorf("account with nickname '%s' already exists", nickname)
	}

	// Create key pair from hex
	keyPairRaw, err := FromHex(privateKeyHex, "")
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Derive public key if not provided
	if len(keyPairRaw.PublicKey) == 0 {
		keyPairRaw.PublicKey = PublicKeyFromPrivate(keyPairRaw.PrivateKey)
	}

	// Encrypt the key pair
	encryptedKeyPair, err := EncryptKeyPair(keyPairRaw, password)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt key pair: %w", err)
	}

	// Derive address from public key
	address, err := AddressFromPublicKey(keyPairRaw.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to derive address: %w", err)
	}

	// Create account
	account := &Account{
		Nickname:         nickname,
		Address:          address,
		KeyPair:          encryptedKeyPair,
		Balance:          0,
		CandidateBalance: 0,
		Status:           AccountStatusOK,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Store account
	am.accounts[nickname] = account

	return account, nil
}

// ImportAccountFromPrivateKey imports an account from a Massa base58-encoded private key
func (am *AccountManager) ImportAccountFromPrivateKey(nickname, privateKeyBase58, password string) (*Account, error) {
	if nickname == "" {
		return nil, errors.New("nickname cannot be empty")
	}

	if password == "" {
		return nil, errors.New("password cannot be empty")
	}

	if privateKeyBase58 == "" {
		return nil, errors.New("private key cannot be empty")
	}

	// Check if account already exists
	if _, exists := am.accounts[nickname]; exists {
		return nil, fmt.Errorf("account with nickname '%s' already exists", nickname)
	}

	// Create key pair from base58-encoded private key
	keyPairRaw, err := FromBase58PrivateKey(privateKeyBase58)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Encrypt the key pair
	encryptedKeyPair, err := EncryptKeyPair(keyPairRaw, password)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt key pair: %w", err)
	}

	// Derive address from public key
	address, err := AddressFromPublicKey(keyPairRaw.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to derive address: %w", err)
	}

	// Create account
	account := &Account{
		Nickname:         nickname,
		Address:          address,
		KeyPair:          encryptedKeyPair,
		Balance:          0,
		CandidateBalance: 0,
		Status:           AccountStatusOK,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Store account
	am.accounts[nickname] = account

	return account, nil
}

// GetAccount retrieves an account by nickname
func (am *AccountManager) GetAccount(nickname string) (*Account, error) {
	account, exists := am.accounts[nickname]
	if !exists {
		return nil, fmt.Errorf("account with nickname '%s' not found", nickname)
	}

	return account, nil
}

// ListAccounts returns all accounts
func (am *AccountManager) ListAccounts() []*Account {
	accounts := make([]*Account, 0, len(am.accounts))
	for _, account := range am.accounts {
		accounts = append(accounts, account)
	}
	return accounts
}

// DeleteAccount removes an account by nickname
func (am *AccountManager) DeleteAccount(nickname string) error {
	if _, exists := am.accounts[nickname]; !exists {
		return fmt.Errorf("account with nickname '%s' not found", nickname)
	}

	delete(am.accounts, nickname)
	return nil
}

// UnlockAccount unlocks an account and returns the raw key pair
func (am *AccountManager) UnlockAccount(nickname, password string) (*KeyPairRaw, error) {
	account, exists := am.accounts[nickname]
	if !exists {
		return nil, fmt.Errorf("account with nickname '%s' not found", nickname)
	}

	if account.Status == AccountStatusCorrupted {
		return nil, fmt.Errorf("account '%s' is corrupted", nickname)
	}

	// Decrypt key pair
	keyPairRaw, err := DecryptKeyPair(account.KeyPair, password)
	if err != nil {
		return nil, fmt.Errorf("failed to unlock account: %w", err)
	}

	return keyPairRaw, nil
}

// UpdateBalance updates the balance of an account
func (am *AccountManager) UpdateBalance(nickname string, balance, candidateBalance uint64) error {
	account, exists := am.accounts[nickname]
	if !exists {
		return fmt.Errorf("account with nickname '%s' not found", nickname)
	}

	account.Balance = balance
	account.CandidateBalance = candidateBalance
	account.UpdatedAt = time.Now()

	return nil
}

// ChangePassword changes the password of an account
func (am *AccountManager) ChangePassword(nickname, oldPassword, newPassword string) error {
	account, exists := am.accounts[nickname]
	if !exists {
		return fmt.Errorf("account with nickname '%s' not found", nickname)
	}

	if newPassword == "" {
		return errors.New("new password cannot be empty")
	}

	// Decrypt with old password
	keyPairRaw, err := DecryptKeyPair(account.KeyPair, oldPassword)
	if err != nil {
		return fmt.Errorf("failed to decrypt with old password: %w", err)
	}

	// Encrypt with new password
	newEncryptedKeyPair, err := EncryptKeyPair(keyPairRaw, newPassword)
	if err != nil {
		return fmt.Errorf("failed to encrypt with new password: %w", err)
	}

	// Update account
	account.KeyPair = newEncryptedKeyPair
	account.UpdatedAt = time.Now()

	return nil
}

// ExportPrivateKey exports the private key of an account in hexadecimal format
func (am *AccountManager) ExportPrivateKey(nickname, password string) (string, error) {
	keyPairRaw, err := am.UnlockAccount(nickname, password)
	if err != nil {
		return "", err
	}

	privateKeyHex, _ := keyPairRaw.ToHex()
	return privateKeyHex, nil
}

// ExportPrivateKeyBase58 exports the private key of an account in Massa base58 format
func (am *AccountManager) ExportPrivateKeyBase58(nickname, password string) (string, error) {
	keyPairRaw, err := am.UnlockAccount(nickname, password)
	if err != nil {
		return "", err
	}

	return keyPairRaw.ToBase58(), nil
}

// ToJSON converts an account to JSON
func (a *Account) ToJSON() ([]byte, error) {
	return json.Marshal(a)
}

// FromJSON creates an account from JSON
func FromJSON(data []byte) (*Account, error) {
	var account Account
	err := json.Unmarshal(data, &account)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal account: %w", err)
	}

	return &account, nil
}

// GetPublicKey returns the public key of the account
func (a *Account) GetPublicKey() (string, error) {
	if a.KeyPair == nil {
		return "", errors.New("account has no key pair")
	}

	return a.KeyPair.PublicKey, nil
}

// IsLocked checks if the account is locked (requires password to use)
func (a *Account) IsLocked() bool {
	return a.Status == AccountStatusLocked
}

// IsCorrupted checks if the account is corrupted
func (a *Account) IsCorrupted() bool {
	return a.Status == AccountStatusCorrupted
}

// String returns a string representation of the account
func (a *Account) String() string {
	return fmt.Sprintf("Account{Nickname: %s, Address: %s, Status: %s}", a.Nickname, a.Address.String(), a.Status)
}

// Unlock decrypts the account's KeyPair with the provided password and returns the raw KeyPair
func (a *Account) Unlock(password string) (*KeyPairRaw, error) {
	if a == nil {
		return nil, fmt.Errorf("nil account")
	}

	if a.KeyPair == nil {
		return nil, fmt.Errorf("account has no key pair")
	}

	kp, err := DecryptKeyPair(a.KeyPair, password)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt key pair: %w", err)
	}

	return kp, nil
}
