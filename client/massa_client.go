package client

import (
	"context"
	"encoding/base64"
	"fmt"

	pb "github.com/nafsilabs/massa-go/client/proto/massa/api/v1"
	model "github.com/nafsilabs/massa-go/client/proto/massa/model/v1"
	"github.com/nafsilabs/massa-go/wallet"
	"google.golang.org/protobuf/proto"
)

// Helper function to create NativeAmount from uint64 (assuming mantissa with scale 0)
func NewNativeAmount(amount uint64) *model.NativeAmount {
	return &model.NativeAmount{
		Mantissa: amount,
		Scale:    0,
	}
}

// MassaOperation represents a Massa blockchain operation that can be signed and sent
type MassaOperation interface {
	// ToOperationType converts the operation to its protobuf OperationType representation
	ToOperationType() *model.OperationType
}

// TransactionOp represents a simple transaction operation
type TransactionOp struct {
	RecipientAddress string
	Amount           uint64
}

func (op *TransactionOp) ToOperationType() *model.OperationType {
	return &model.OperationType{
		Type: &model.OperationType_Transaction{
			Transaction: &model.Transaction{
				RecipientAddress: op.RecipientAddress,
				Amount:           NewNativeAmount(op.Amount),
			},
		},
	}
}

// BuyRollsOp represents a roll purchase operation
type BuyRollsOp struct {
	RollCount uint64
}

func (op *BuyRollsOp) ToOperationType() *model.OperationType {
	return &model.OperationType{
		Type: &model.OperationType_RollBuy{
			RollBuy: &model.RollBuy{
				RollCount: op.RollCount,
			},
		},
	}
}

// SellRollsOp represents a roll sale operation
type SellRollsOp struct {
	RollCount uint64
}

func (op *SellRollsOp) ToOperationType() *model.OperationType {
	return &model.OperationType{
		Type: &model.OperationType_RollSell{
			RollSell: &model.RollSell{
				RollCount: op.RollCount,
			},
		},
	}
}

// ExecuteSCOp represents a smart contract execution operation (deploy bytecode)
type ExecuteSCOp struct {
	Bytecode  []byte
	MaxGas    uint64
	MaxCoins  uint64
	Datastore []*model.BytesMapFieldEntry
}

func (op *ExecuteSCOp) ToOperationType() *model.OperationType {
	return &model.OperationType{
		Type: &model.OperationType_ExecuteSc{
			ExecuteSc: &model.ExecuteSC{
				Data:      op.Bytecode,
				MaxGas:    op.MaxGas,
				MaxCoins:  op.MaxCoins,
				Datastore: op.Datastore,
			},
		},
	}
}

// CallSCOp represents a smart contract call operation
type CallSCOp struct {
	TargetAddress  string
	TargetFunction string
	Parameter      []byte
	MaxGas         uint64
	Coins          uint64
}

func (op *CallSCOp) ToOperationType() *model.OperationType {
	return &model.OperationType{
		Type: &model.OperationType_CallSc{
			CallSc: &model.CallSC{
				TargetAddress:  op.TargetAddress,
				TargetFunction: op.TargetFunction,
				Parameter:      op.Parameter,
				MaxGas:         op.MaxGas,
				Coins:          NewNativeAmount(op.Coins),
			},
		},
	}
}

// OperationSender provides high-level methods for signing and sending operations
type OperationSender struct {
	client *MassaGrpcClient
}

// NewOperationSender creates a new OperationSender
func NewOperationSender(client *MassaGrpcClient) *OperationSender {
	return &OperationSender{
		client: client,
	}
}

// SendOperation signs and sends a single operation to the Massa network
// Parameters:
//   - ctx: Context for the request
//   - op: The operation to send
//   - fee: Fee amount in NanoMassa
//   - expirePeriod: Period after which the operation expires
//   - signature: Base64 or hex encoded signature of the operation
//   - publicKey: Base64 or hex encoded public key of the signer
//
// Returns:
//   - operationID: The ID of the sent operation
//   - error: Any error that occurred
func (s *OperationSender) SendOperation(
	ctx context.Context,
	signedOpBytes []byte,
	// op MassaOperation,
	// fee uint64,
	// expirePeriod uint64,
	// signature string,
	// publicKey string,
) (string, error) {
	// // Create the operation
	// operation := &model.Operation{
	// 	Fee:          NewNativeAmount(fee),
	// 	ExpirePeriod: expirePeriod,
	// 	Op:           op.ToOperationType(),
	// }

	// // Serialize the operation content
	// _, err := proto.Marshal(operation)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to serialize operation: %w", err)
	// }

	// // Create signed operation
	// signedOp := &model.SignedOperation{
	// 	ContentCreatorPubKey: publicKey,
	// 	Signature:            signature,
	// 	//Content:              operation,
	// }

	// // Serialize the signed operation
	// signedOpBytes, err := proto.Marshal(signedOp)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to serialize signed operation: %w", err)
	// }

	// Send via bidirectional gRPC stream
	stream, err := s.client.PublicClient.SendOperations(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to open stream: %w", err)
	}

	// Send the operation
	req := &pb.SendOperationsRequest{
		Operations: [][]byte{signedOpBytes},
	}
	if err := stream.Send(req); err != nil {
		return "", fmt.Errorf("failed to send operation: %w", err)
	}

	// Close send and receive response
	if err := stream.CloseSend(); err != nil {
		return "", fmt.Errorf("failed to close send: %w", err)
	}

	resp, err := stream.Recv()
	if err != nil {
		return "", fmt.Errorf("failed to receive response: %w", err)
	}

	// Check for error in response
	if respError := resp.GetError(); respError != nil {
		return "", fmt.Errorf("operation error: %s", respError.Message)
	}

	// Get operation IDs
	opIds := resp.GetOperationIds()
	if opIds == nil || len(opIds.OperationIds) == 0 {
		return "", fmt.Errorf("no operation ID returned")
	}

	return opIds.OperationIds[0], nil
}

// SendOperations signs and sends multiple operations to the Massa network
// Parameters:
//   - ctx: Context for the request
//   - operations: List of operations with their signatures and public keys
//   - fee: Fee amount in NanoMassa (applied to all operations)
//   - expirePeriod: Period after which operations expire (applied to all)
//
// Returns:
//   - operationIDs: The IDs of the sent operations
//   - error: Any error that occurred
func (s *OperationSender) SendOperations(
	ctx context.Context,
	operations []struct {
		Op        MassaOperation
		Signature string
		PublicKey string
	},
	fee uint64,
	expirePeriod uint64,
) ([]string, error) {
	serializedOps := make([][]byte, 0, len(operations))

	for _, opData := range operations {
		// Create the operation
		operation := &model.Operation{
			Fee:          NewNativeAmount(fee),
			ExpirePeriod: expirePeriod,
			Op:           opData.Op.ToOperationType(),
		}

		// Serialize the operation content
		_, err := proto.Marshal(operation)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize operation: %w", err)
		}

		// Create signed operation
		signedOp := &model.SignedOperation{
			ContentCreatorPubKey: opData.PublicKey,
			Signature:            opData.Signature,
		}

		// Serialize the signed operation
		signedOpBytes, err := proto.Marshal(signedOp)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize signed operation: %w", err)
		}

		serializedOps = append(serializedOps, signedOpBytes)
	}

	// Send via bidirectional gRPC stream
	stream, err := s.client.PublicClient.SendOperations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open stream: %w", err)
	}

	// Send the operations
	req := &pb.SendOperationsRequest{
		Operations: serializedOps,
	}
	if err := stream.Send(req); err != nil {
		return nil, fmt.Errorf("failed to send operations: %w", err)
	}

	// Close send and receive response
	if err := stream.CloseSend(); err != nil {
		return nil, fmt.Errorf("failed to close send: %w", err)
	}

	resp, err := stream.Recv()
	if err != nil {
		return nil, fmt.Errorf("failed to receive response: %w", err)
	}

	// Check for error in response
	if respError := resp.GetError(); respError != nil {
		return nil, fmt.Errorf("operation error: %s", respError.Message)
	}

	// Get operation IDs
	opIds := resp.GetOperationIds()
	if opIds == nil {
		return nil, fmt.Errorf("no operation IDs returned")
	}

	return opIds.OperationIds, nil
}

// Helper method to encode bytes to base64
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Helper method to decode base64 to bytes
func DecodeBase64(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}

// // SendTransaction is a convenience method for sending a simple transaction
// func (s *OperationSender) SendTransaction(
// 	ctx context.Context,
// 	recipientAddress string,
// 	amount uint64,
// 	fee uint64,
// 	expirePeriod uint64,
// 	signature string,
// 	publicKey string,
// ) (string, error) {
// 	op := &TransactionOp{
// 		RecipientAddress: recipientAddress,
// 		Amount:           amount,
// 	}
// 	return s.SendOperation(ctx, op, fee, expirePeriod, signature, publicKey)
// }

// // SendBuyRolls is a convenience method for buying rolls
// func (s *OperationSender) SendBuyRolls(
// 	ctx context.Context,
// 	rollCount uint64,
// 	fee uint64,
// 	expirePeriod uint64,
// 	signature string,
// 	publicKey string,
// ) (string, error) {
// 	op := &BuyRollsOp{
// 		RollCount: rollCount,
// 	}
// 	return s.SendOperation(ctx, op, fee, expirePeriod, signature, publicKey)
// }

// // SendSellRolls is a convenience method for selling rolls
// func (s *OperationSender) SendSellRolls(
// 	ctx context.Context,
// 	rollCount uint64,
// 	fee uint64,
// 	expirePeriod uint64,
// 	signature string,
// 	publicKey string,
// ) (string, error) {
// 	op := &SellRollsOp{
// 		RollCount: rollCount,
// 	}
// 	return s.SendOperation(ctx, op, fee, expirePeriod, signature, publicKey)
// }

// // SendExecuteSC is a convenience method for executing/deploying smart contract bytecode
// func (s *OperationSender) SendExecuteSC(
// 	ctx context.Context,
// 	bytecode []byte,
// 	maxGas uint64,
// 	maxCoins uint64,
// 	datastore []*model.BytesMapFieldEntry,
// 	fee uint64,
// 	expirePeriod uint64,
// 	signature string,
// 	publicKey string,
// ) (string, error) {
// 	op := &ExecuteSCOp{
// 		Bytecode:  bytecode,
// 		MaxGas:    maxGas,
// 		MaxCoins:  maxCoins,
// 		Datastore: datastore,
// 	}
// 	return s.SendOperation(ctx, op, fee, expirePeriod, signature, publicKey)
// }

// // SendCallSC is a convenience method for calling a smart contract function
// func (s *OperationSender) SendCallSC(
// 	ctx context.Context,
// 	targetAddress string,
// 	targetFunction string,
// 	parameter []byte,
// 	maxGas uint64,
// 	coins uint64,
// 	fee uint64,
// 	expirePeriod uint64,
// 	signature string,
// 	publicKey string,
// ) (string, error) {
// 	op := &CallSCOp{
// 		TargetAddress:  targetAddress,
// 		TargetFunction: targetFunction,
// 		Parameter:      parameter,
// 		MaxGas:         maxGas,
// 		Coins:          coins,
// 	}
// 	return s.SendOperation(ctx, op, fee, expirePeriod, signature, publicKey)
// }

// ============================================================================
// Wallet-Integrated Operation Sender
// ============================================================================

// WalletOperationSender provides high-level methods for signing and sending
// operations using wallet accounts (automatic signing)
type WalletOperationSender struct {
	client         *MassaGrpcClient
	accountManager *wallet.AccountManager
}

// NewWalletOperationSender creates a new wallet-integrated operation sender
func NewWalletOperationSender(client *MassaGrpcClient, accountManager *wallet.AccountManager) *WalletOperationSender {
	return &WalletOperationSender{
		client:         client,
		accountManager: accountManager,
	}
}

// signOperation signs an operation using the wallet account
func (w *WalletOperationSender) signOperation(
	nickname string,
	password string,
	operation *model.Operation,
) (signature string, publicKey string, err error) {
	// Unlock the account to get the key pair
	keyPair, err := w.accountManager.UnlockAccount(nickname, password)
	if err != nil {
		return "", "", fmt.Errorf("failed to unlock account: %w", err)
	}

	// Serialize the operation for signing
	operationBytes, err := proto.Marshal(operation)
	if err != nil {
		return "", "", fmt.Errorf("failed to serialize operation: %w", err)
	}

	// Hash the operation
	//operationHash := sha256.Sum256(operationBytes)
	//operationHash := utils.brak

	// Sign the hash
	signatureBytes := keyPair.VersionedSignatureBytes(operationBytes)

	// Encode to base64
	signature = base64.StdEncoding.EncodeToString(signatureBytes)
	publicKey = base64.StdEncoding.EncodeToString(keyPair.VersionedPublicKeyBytes())

	return signature, publicKey, nil
}

// SendOperationWithWallet signs and sends a single operation using a wallet account
// Parameters:
//   - ctx: Context for the request
//   - nickname: Account nickname in the wallet
//   - password: Account password
//   - op: The operation to send
//   - fee: Fee amount in NanoMassa
//   - expirePeriod: Period after which the operation expires
//
// Returns:
//   - operationID: The ID of the sent operation
//   - error: Any error that occurred
func (w *WalletOperationSender) SendOperationWithWallet(
	ctx context.Context,
	nickname string,
	password string,
	op MassaOperation,
	fee uint64,
	expirePeriod uint64,
) (string, error) {
	// Create the operation
	operation := &model.Operation{
		Fee:          NewNativeAmount(fee),
		ExpirePeriod: expirePeriod,
		Op:           op.ToOperationType(),
	}

	// Sign the operation
	signature, publicKey, err := w.signOperation(nickname, password, operation)
	if err != nil {
		return "", fmt.Errorf("failed to sign operation: %w", err)
	}

	// Create signed operation
	signedOp := &model.SignedOperation{
		ContentCreatorPubKey: publicKey,
		Signature:            signature,
	}

	// Serialize the signed operation
	signedOpBytes, err := proto.Marshal(signedOp)
	if err != nil {
		return "", fmt.Errorf("failed to serialize signed operation: %w", err)
	}

	// Send via bidirectional gRPC stream
	stream, err := w.client.PublicClient.SendOperations(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to open stream: %w", err)
	}

	// Send the operation
	req := &pb.SendOperationsRequest{
		Operations: [][]byte{signedOpBytes},
	}
	if err := stream.Send(req); err != nil {
		return "", fmt.Errorf("failed to send operation: %w", err)
	}

	// Close send and receive response
	if err := stream.CloseSend(); err != nil {
		return "", fmt.Errorf("failed to close send: %w", err)
	}

	resp, err := stream.Recv()
	if err != nil {
		return "", fmt.Errorf("failed to receive response: %w", err)
	}

	// Check for error in response
	if respError := resp.GetError(); respError != nil {
		return "", fmt.Errorf("operation error: %s", respError.Message)
	}

	// Get operation IDs
	opIds := resp.GetOperationIds()
	if opIds == nil || len(opIds.OperationIds) == 0 {
		return "", fmt.Errorf("no operation ID returned")
	}

	return opIds.OperationIds[0], nil
}

// SendTransactionWithWallet is a convenience method for sending a transaction using wallet
func (w *WalletOperationSender) SendTransactionWithWallet(
	ctx context.Context,
	nickname string,
	password string,
	recipientAddress string,
	amount uint64,
	fee uint64,
	expirePeriod uint64,
) (string, error) {
	op := &TransactionOp{
		RecipientAddress: recipientAddress,
		Amount:           amount,
	}
	return w.SendOperationWithWallet(ctx, nickname, password, op, fee, expirePeriod)
}

// SendBuyRollsWithWallet is a convenience method for buying rolls using wallet
func (w *WalletOperationSender) SendBuyRollsWithWallet(
	ctx context.Context,
	nickname string,
	password string,
	rollCount uint64,
	fee uint64,
	expirePeriod uint64,
) (string, error) {
	op := &BuyRollsOp{
		RollCount: rollCount,
	}
	return w.SendOperationWithWallet(ctx, nickname, password, op, fee, expirePeriod)
}

// SendSellRollsWithWallet is a convenience method for selling rolls using wallet
func (w *WalletOperationSender) SendSellRollsWithWallet(
	ctx context.Context,
	nickname string,
	password string,
	rollCount uint64,
	fee uint64,
	expirePeriod uint64,
) (string, error) {
	op := &SellRollsOp{
		RollCount: rollCount,
	}
	return w.SendOperationWithWallet(ctx, nickname, password, op, fee, expirePeriod)
}

// SendExecuteSCWithWallet is a convenience method for deploying smart contracts using wallet
func (w *WalletOperationSender) SendExecuteSCWithWallet(
	ctx context.Context,
	nickname string,
	password string,
	bytecode []byte,
	maxGas uint64,
	maxCoins uint64,
	datastore []*model.BytesMapFieldEntry,
	fee uint64,
	expirePeriod uint64,
) (string, error) {
	op := &ExecuteSCOp{
		Bytecode:  bytecode,
		MaxGas:    maxGas,
		MaxCoins:  maxCoins,
		Datastore: datastore,
	}
	return w.SendOperationWithWallet(ctx, nickname, password, op, fee, expirePeriod)
}

// SendCallSCWithWallet is a convenience method for calling smart contracts using wallet
func (w *WalletOperationSender) SendCallSCWithWallet(
	ctx context.Context,
	nickname string,
	password string,
	targetAddress string,
	targetFunction string,
	parameter []byte,
	maxGas uint64,
	coins uint64,
	fee uint64,
	expirePeriod uint64,
) (string, error) {
	op := &CallSCOp{
		TargetAddress:  targetAddress,
		TargetFunction: targetFunction,
		Parameter:      parameter,
		MaxGas:         maxGas,
		Coins:          coins,
	}
	return w.SendOperationWithWallet(ctx, nickname, password, op, fee, expirePeriod)
}
