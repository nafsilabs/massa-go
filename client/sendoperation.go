package client

import (
	"bytes"
	"context"
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

//nolint:tagliatelle
type sendOperationsReq struct {
	Operations [][]byte `protobuf:"bytes,1,rep,name=operations,proto3" json:"operations,omitempty"`
}

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

// Call uses the plugin wallet to sign an operation, then send the call to blockchain.
func Call(
	c *Client,
	expiry uint64,
	fee uint64,
	operation Operation,
	account *wallet.Account,
	password string,
	description string,
) (*OperationResponse, error) {

	//msg contains expiry, fee and operation
	operationData, _, err := MakeOperation(c, expiry, fee, operation)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("operation data: %v\n", operationData)

	//versioned public key
	keyPairRaw, err := account.Unlock(password)
	if err != nil {
		return nil, err
	}

	versionedPublicKey := keyPairRaw.VersionedPublicKeyBytes()

	content := createOperationContent(c.ChainID, versionedPublicKey, operationData)
	//fmt.Printf("Signature data: %v\n", content)

	signature := keyPairRaw.VersionedSignatureBytes(content)

	//fmt.Printf("Signature: %v\n", signature)

	// //data to be transmitted
	msg := make([]byte, 0)
	msg = append(msg, signature...)
	msg = append(msg, versionedPublicKey...)
	msg = append(msg, content...)

	// //fmt.Printf("data to be transmitted: %v\n", msg)
	// //return nil, nil
	// resp, err := MakeRPCCall(msg, c)
	// if err != nil {
	// 	return nil, err
	// }

	// apiClient := NewApiClient()
	// if err := apiClient.Init(BuildnetURL); err != nil {
	// 	fmt.Printf("Error getting info: %v\n", err)
	// }

	// op, err := apiClient.SendOperation(msg)
	// if err != nil {
	// 	fmt.Printf("error sending operation")
	// }

	return &OperationResponse{OperationID: "nil"}, nil
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

func MakeRPCCall(msg []byte, c *Client) ([]string, error) {

	sendOpParams := sendOperationsReq{
		Operations: [][]byte{msg},
	}

	rawResponse, err := c.RPCClient.Call(
		context.Background(),
		"send_operations",
		sendOpParams,
	)
	if err != nil {
		return nil, fmt.Errorf("calling send_operations jsonrpc with '%+v': %w", sendOpParams, err)
	}

	if rawResponse.Error != nil {
		return nil, fmt.Errorf("receiving send_operations response: %w", rawResponse.Error)
	}

	var resp []string

	err = rawResponse.GetObject(&resp)
	if err != nil {
		return nil, fmt.Errorf("parsing send_operations jsonrpc response '%+v': %w", rawResponse, err)
	}

	return resp, nil
}

// Call uses the plugin wallet to sign an operation, then send the call to blockchain.
func MakeOperation2(
	c *Client,
	expiry uint64,
	fee uint64,
	operation Operation,
	account *wallet.Account,
	password string,
	description string,
) ([]byte, error) {

	//msg contains expiry, fee and operation
	operationData, _, err := MakeOperation(c, expiry, fee, operation)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("operation data: %v\n", operationData)

	keyPairRaw, err := account.Unlock(password)
	if err != nil {
		return nil, err
	}

	versionedPublicKey := keyPairRaw.VersionedPublicKeyBytes()

	content := createOperationContent(c.ChainID, versionedPublicKey, operationData)
	//fmt.Printf("Signature data: %v\n", content)

	signature := keyPairRaw.VersionedSignatureBytes(content)

	//fmt.Printf("Signature: %v\n", signature)

	// //data to be transmitted
	msg := make([]byte, 0)
	msg = append(msg, signature...)
	msg = append(msg, versionedPublicKey...)
	msg = append(msg, operationData...)

	return msg, nil
}

// MakeOperation concatinets fee, expiry and operation into a single message byte slice,
// then encodes it as base64 for RPC transmission.
// It returns the message byte slice, its base64 encoding, and an error if any.
func MakeOperation(c *Client, expiry uint64, fee uint64, operation Operation) ([]byte, string, error) {
	exp, err := NextSlot(c)
	if err != nil {
		return nil, "", fmt.Errorf("calling NextSlot: %w", err)
	}

	expiry += exp

	msg := message(expiry, fee, operation)

	msgB64 := b64.StdEncoding.EncodeToString(msg)

	return msg, msgB64, nil
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
