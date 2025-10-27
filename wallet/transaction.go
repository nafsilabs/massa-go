package wallet

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/ed25519"
)

// TransactionSigner provides transaction signing capabilities
type TransactionSigner struct {
	accountManager *AccountManager
}

// NewTransactionSigner creates a new transaction signer
func NewTransactionSigner(accountManager *AccountManager) *TransactionSigner {
	return &TransactionSigner{
		accountManager: accountManager,
	}
}

// TransactionType represents the type of transaction
type TransactionType uint8

const (
	TransactionTypeTransfer TransactionType = iota
	TransactionTypeCallSC
	TransactionTypeExecuteSC
	TransactionTypeBuyRolls
	TransactionTypeSellRolls
)

// Operation represents a Massa blockchain operation
type Operation struct {
	// Operation type
	Type TransactionType `json:"type"`

	// Fee in nanoMAS
	Fee uint64 `json:"fee"`

	// Expiry period
	ExpiryPeriod uint64 `json:"expiryPeriod"`

	// Operation content (varies by type)
	Content interface{} `json:"content"`
}

// TransferOperation represents a coin transfer operation
type TransferOperation struct {
	// Recipient address
	RecipientAddress string `json:"recipientAddress"`

	// Amount in nanoMAS
	Amount uint64 `json:"amount"`
}

// CallSCOperation represents a smart contract call operation
type CallSCOperation struct {
	// Target smart contract address
	TargetAddress string `json:"targetAddress"`

	// Function name to call
	FunctionName string `json:"functionName"`

	// Function parameters
	Parameters []byte `json:"parameters"`

	// Maximum gas for execution
	MaxGas uint64 `json:"maxGas"`

	// Coins to send with the call (in nanoMAS)
	Coins uint64 `json:"coins"`
}

// ExecuteSCOperation represents smart contract execution (deployment)
type ExecuteSCOperation struct {
	// Bytecode to execute
	Bytecode []byte `json:"bytecode"`

	// Maximum gas for execution
	MaxGas uint64 `json:"maxGas"`

	// Coins to send to the contract (in nanoMAS)
	Coins uint64 `json:"coins"`

	// Datastore for the contract
	Datastore map[string][]byte `json:"datastore"`
}

// RollOperation represents roll buying/selling operations
type RollOperation struct {
	// Number of rolls
	RollCount uint64 `json:"rollCount"`
}

// SignedOperation represents a signed operation ready for broadcast
type SignedOperation struct {
	// Serialized operation content
	SerializedContent []byte `json:"serializedContent"`

	// Creator's public key
	CreatorPublicKey string `json:"creatorPublicKey"`

	// Operation signature
	Signature string `json:"signature"`
}

// SignatureRequest represents a request to sign an operation
type SignatureRequest struct {
	// Account nickname to use for signing
	Nickname string `json:"nickname"`

	// Account password
	Password string `json:"password"`

	// Operation to sign
	Operation *Operation `json:"operation"`

	// Allow fee edition (for wallets)
	AllowFeeEdition bool `json:"allowFeeEdition"`
}

// SignOperation signs an operation using the specified account
func (ts *TransactionSigner) SignOperation(request *SignatureRequest) (*SignedOperation, error) {
	if request == nil {
		return nil, errors.New("signature request cannot be nil")
	}

	if request.Operation == nil {
		return nil, errors.New("operation cannot be nil")
	}

	// Unlock the account
	keyPair, err := ts.accountManager.UnlockAccount(request.Nickname, request.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to unlock account: %w", err)
	}

	// Serialize the operation
	serializedContent, err := ts.serializeOperation(request.Operation)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize operation: %w", err)
	}

	// Create the message to sign (operation hash)
	messageToSign := sha256.Sum256(serializedContent)

	// Sign the message
	signature := keyPair.Sign(messageToSign[:])

	// Get public key
	publicKeyBytes := keyPair.PublicKey

	return &SignedOperation{
		SerializedContent: serializedContent,
		CreatorPublicKey:  base64.StdEncoding.EncodeToString(publicKeyBytes),
		Signature:         base64.StdEncoding.EncodeToString(signature),
	}, nil
}

// SignMessage signs an arbitrary message using the specified account
func (ts *TransactionSigner) SignMessage(nickname, password string, message []byte) ([]byte, error) {
	// Unlock the account
	keyPair, err := ts.accountManager.UnlockAccount(nickname, password)
	if err != nil {
		return nil, fmt.Errorf("failed to unlock account: %w", err)
	}

	// Sign the message
	signature := keyPair.Sign(message)

	return signature, nil
}

// VerifySignature verifies a signature against a message and public key
func (ts *TransactionSigner) VerifySignature(message, signature, publicKey []byte) bool {
	if len(publicKey) != ed25519.PublicKeySize {
		return false
	}

	return ed25519.Verify(ed25519.PublicKey(publicKey), message, signature)
}

// serializeOperation serializes an operation into bytes for signing
func (ts *TransactionSigner) serializeOperation(op *Operation) ([]byte, error) {
	var buffer []byte

	// Add operation type
	buffer = append(buffer, byte(op.Type))

	// Add fee (8 bytes, little endian)
	feeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(feeBytes, op.Fee)
	buffer = append(buffer, feeBytes...)

	// Add expiry period (8 bytes, little endian)
	expiryBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(expiryBytes, op.ExpiryPeriod)
	buffer = append(buffer, expiryBytes...)

	// Serialize content based on operation type
	contentBytes, err := ts.serializeOperationContent(op.Type, op.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize operation content: %w", err)
	}

	buffer = append(buffer, contentBytes...)

	return buffer, nil
}

// serializeOperationContent serializes operation content based on type
func (ts *TransactionSigner) serializeOperationContent(opType TransactionType, content interface{}) ([]byte, error) {
	switch opType {
	case TransactionTypeTransfer:
		return ts.serializeTransferOperation(content)
	case TransactionTypeCallSC:
		return ts.serializeCallSCOperation(content)
	case TransactionTypeExecuteSC:
		return ts.serializeExecuteSCOperation(content)
	case TransactionTypeBuyRolls, TransactionTypeSellRolls:
		return ts.serializeRollOperation(content)
	default:
		return nil, fmt.Errorf("unsupported operation type: %d", opType)
	}
}

// serializeTransferOperation serializes a transfer operation
func (ts *TransactionSigner) serializeTransferOperation(content interface{}) ([]byte, error) {
	transfer, ok := content.(*TransferOperation)
	if !ok {
		return nil, errors.New("invalid transfer operation content")
	}

	var buffer []byte

	// Add recipient address
	recipientBytes := []byte(transfer.RecipientAddress)
	lengthBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthBytes, uint32(len(recipientBytes)))
	buffer = append(buffer, lengthBytes...)
	buffer = append(buffer, recipientBytes...)

	// Add amount (8 bytes, little endian)
	amountBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(amountBytes, transfer.Amount)
	buffer = append(buffer, amountBytes...)

	return buffer, nil
}

// serializeCallSCOperation serializes a smart contract call operation
func (ts *TransactionSigner) serializeCallSCOperation(content interface{}) ([]byte, error) {
	call, ok := content.(*CallSCOperation)
	if !ok {
		return nil, errors.New("invalid call SC operation content")
	}

	var buffer []byte

	// Add target address
	targetBytes := []byte(call.TargetAddress)
	lengthBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthBytes, uint32(len(targetBytes)))
	buffer = append(buffer, lengthBytes...)
	buffer = append(buffer, targetBytes...)

	// Add function name
	functionBytes := []byte(call.FunctionName)
	binary.LittleEndian.PutUint32(lengthBytes, uint32(len(functionBytes)))
	buffer = append(buffer, lengthBytes...)
	buffer = append(buffer, functionBytes...)

	// Add parameters
	binary.LittleEndian.PutUint32(lengthBytes, uint32(len(call.Parameters)))
	buffer = append(buffer, lengthBytes...)
	buffer = append(buffer, call.Parameters...)

	// Add max gas (8 bytes, little endian)
	gasBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(gasBytes, call.MaxGas)
	buffer = append(buffer, gasBytes...)

	// Add coins (8 bytes, little endian)
	coinsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(coinsBytes, call.Coins)
	buffer = append(buffer, coinsBytes...)

	return buffer, nil
}

// serializeExecuteSCOperation serializes a smart contract execution operation
func (ts *TransactionSigner) serializeExecuteSCOperation(content interface{}) ([]byte, error) {
	execute, ok := content.(*ExecuteSCOperation)
	if !ok {
		return nil, errors.New("invalid execute SC operation content")
	}

	var buffer []byte

	// Add bytecode
	lengthBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthBytes, uint32(len(execute.Bytecode)))
	buffer = append(buffer, lengthBytes...)
	buffer = append(buffer, execute.Bytecode...)

	// Add max gas (8 bytes, little endian)
	gasBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(gasBytes, execute.MaxGas)
	buffer = append(buffer, gasBytes...)

	// Add coins (8 bytes, little endian)
	coinsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(coinsBytes, execute.Coins)
	buffer = append(buffer, coinsBytes...)

	// Add datastore
	datastoreBytes, err := ts.serializeDatastore(execute.Datastore)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize datastore: %w", err)
	}
	binary.LittleEndian.PutUint32(lengthBytes, uint32(len(datastoreBytes)))
	buffer = append(buffer, lengthBytes...)
	buffer = append(buffer, datastoreBytes...)

	return buffer, nil
}

// serializeRollOperation serializes a roll operation
func (ts *TransactionSigner) serializeRollOperation(content interface{}) ([]byte, error) {
	roll, ok := content.(*RollOperation)
	if !ok {
		return nil, errors.New("invalid roll operation content")
	}

	// Add roll count (8 bytes, little endian)
	rollBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(rollBytes, roll.RollCount)

	return rollBytes, nil
}

// serializeDatastore serializes a datastore map
func (ts *TransactionSigner) serializeDatastore(datastore map[string][]byte) ([]byte, error) {
	var buffer []byte

	// Add number of entries
	countBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(countBytes, uint32(len(datastore)))
	buffer = append(buffer, countBytes...)

	// Add each key-value pair
	lengthBytes := make([]byte, 4)
	for key, value := range datastore {
		// Add key
		keyBytes := []byte(key)
		binary.LittleEndian.PutUint32(lengthBytes, uint32(len(keyBytes)))
		buffer = append(buffer, lengthBytes...)
		buffer = append(buffer, keyBytes...)

		// Add value
		binary.LittleEndian.PutUint32(lengthBytes, uint32(len(value)))
		buffer = append(buffer, lengthBytes...)
		buffer = append(buffer, value...)
	}

	return buffer, nil
}

// CreateTransferOperation creates a new transfer operation
func CreateTransferOperation(recipientAddress string, amount, fee, expiryPeriod uint64) *Operation {
	return &Operation{
		Type:         TransactionTypeTransfer,
		Fee:          fee,
		ExpiryPeriod: expiryPeriod,
		Content: &TransferOperation{
			RecipientAddress: recipientAddress,
			Amount:           amount,
		},
	}
}

// CreateCallSCOperation creates a new smart contract call operation
func CreateCallSCOperation(targetAddress, functionName string, parameters []byte, maxGas, coins, fee, expiryPeriod uint64) *Operation {
	return &Operation{
		Type:         TransactionTypeCallSC,
		Fee:          fee,
		ExpiryPeriod: expiryPeriod,
		Content: &CallSCOperation{
			TargetAddress: targetAddress,
			FunctionName:  functionName,
			Parameters:    parameters,
			MaxGas:        maxGas,
			Coins:         coins,
		},
	}
}

// CreateExecuteSCOperation creates a new smart contract execution operation
func CreateExecuteSCOperation(bytecode []byte, maxGas, coins, fee, expiryPeriod uint64, datastore map[string][]byte) *Operation {
	if datastore == nil {
		datastore = make(map[string][]byte)
	}

	return &Operation{
		Type:         TransactionTypeExecuteSC,
		Fee:          fee,
		ExpiryPeriod: expiryPeriod,
		Content: &ExecuteSCOperation{
			Bytecode:  bytecode,
			MaxGas:    maxGas,
			Coins:     coins,
			Datastore: datastore,
		},
	}
}

// ToJSON converts a signed operation to JSON
func (so *SignedOperation) ToJSON() ([]byte, error) {
	return json.Marshal(so)
}

// GetCurrentPeriod returns a reasonable expiry period (current + 3 periods)
func GetCurrentPeriod() uint64 {
	// This is a placeholder - in a real implementation, this would query the blockchain
	// For now, we'll use a timestamp-based approximation
	// Massa has approximately 16-second periods
	return uint64(time.Now().Unix()/16) + 3
}
