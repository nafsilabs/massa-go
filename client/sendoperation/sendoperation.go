package sendoperation

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/nafsilabs/massa-go/utils"
	"github.com/nafsilabs/massa-go/wallet"
)

const (
	MaxGasAllowedExecuteSC     = 3_980_167_295
	MaxGasAllowedCallSC        = 4_294_167_295
	DefaultExpiryInSlot        = 3
	DefaultFee                 = 0
	accountCreationStorageCost = 1_000_000
	StorageCostPerByte         = 100_000
	StorageEntryBaseBytes      = 4
	OneMassa                   = 1_000_000_000
)

type Operation interface {
	Content() (interface{}, error)
	Message() []byte
}

type OperationResponse struct {
	OperationID string
}

type OperationContent struct {
	Description string `json:"description"`
	Operation   string `json:"operation"`
	//nolint:tagliatelle
	ChainID uint64 `json:"chainId"`
}

// // MakeOperation concatinets fee, expiry and operation into a single message byte slice,
// then encodes it as base64 for RPC transmission.
// It returns the message byte slice, its base64 encoding, and an error if any.
func MakeOperation(
	expiry uint64,
	fee uint64,
	operation Operation,
	account *wallet.Account,
	password string,
	chainID utils.NetworkType,
) ([]byte, error) {
	operationData := message(expiry, fee, operation)
	keyPairRaw, err := account.Unlock(password)
	if err != nil {
		return nil, err
	}

	versionedPublicKey := keyPairRaw.VersionedPublicKeyBytes()

	content := createOperationContent(chainID, versionedPublicKey, operationData)

	signature := keyPairRaw.VersionedSignatureBytes(content)

	// //data to be transmitted
	msg := make([]byte, 0)
	msg = append(msg, signature...)
	msg = append(msg, versionedPublicKey...)
	msg = append(msg, operationData...)

	return msg, nil
}

// message constructs a message byte slice with the provided expiry, fee, and operation.
// It returns the composed message.
func message(expiry uint64, fee uint64, operation Operation) []byte {
	msg := make([]byte, 0)
	buf := make([]byte, binary.MaxVarintLen64)
	// fee
	nbBytes := binary.PutUvarint(buf, fee)
	msg = append(msg, buf[:nbBytes]...)

	// expiration
	nbBytes = binary.PutUvarint(buf, expiry)
	msg = append(msg, buf[:nbBytes]...)

	// operation
	msg = append(msg, operation.Message()...)

	return msg
}

func createOperationContent(chainID utils.NetworkType, versionedPublicKey, operation []byte) []byte {
	msg := make([]byte, 0)
	// network type
	networkTypeData := chainID.Serialize()
	msg = append(msg, networkTypeData...)

	//public key
	msg = append(msg, versionedPublicKey...)

	// operation
	msg = append(msg, operation...)

	return msg
}

// DecodeMessage64 decodes a base64-encoded message and extracts the fee, expiry,
// and the actual message content(operation : CallSc, ExecuteSC, BuyRoll, SellRoll).
// It returns the decoded message, fee, expiry, and an error if any.
func DecodeMessage64(msgB64 string) ([]byte, uint64, uint64, error) {
	// Decode the base64-encoded message
	decodedMsg, err := b64.StdEncoding.DecodeString(msgB64)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("base64 decoding error: %w", err)
	}

	// Read the encoded fee from the decoded message and move the buffer index
	fee, bytesRead := binary.Uvarint(decodedMsg)
	if bytesRead <= 0 {
		return nil, 0, 0, errors.New("failed to read fee")
	}

	decodedMsg = decodedMsg[bytesRead:]

	// Read the encoded expiry from the decoded message and move the buffer index
	expiry, bytesRead := binary.Uvarint(decodedMsg)
	if bytesRead <= 0 {
		return nil, 0, 0, errors.New("failed to read expiry")
	}

	decodedMsg = decodedMsg[bytesRead:]

	return decodedMsg, fee, expiry, nil
}

// DecodeOperationType decodes a byte slice to retrieve the operation ID.
// Can only be used on operation specific data i.e. after fee and expiry have been removed
// i.e. The parameter should be the result of DecodeMessage64.
func DecodeOperationType(data []byte) (uint64, error) {
	buf := bytes.NewReader(data)

	// Read operation type
	opType, err := binary.ReadUvarint(buf)
	if err != nil {
		return 0, fmt.Errorf("failed to read operation type: %w", err)
	}

	return opType, nil
}

func StorageCostForEntry(keyByteLength, valueByteLength int) (int, error) {
	return (valueByteLength + keyByteLength + StorageEntryBaseBytes) * StorageCostPerByte, nil
}

func AccountCreationStorageCost() (int, error) {
	return accountCreationStorageCost, nil
}
