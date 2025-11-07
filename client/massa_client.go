package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	pb "github.com/nafsilabs/massa-go/client/proto/massa/api/v1"
	model "github.com/nafsilabs/massa-go/client/proto/massa/model/v1"
	"github.com/nafsilabs/massa-go/client/sendoperation"
	"github.com/nafsilabs/massa-go/client/sendoperation/buyrolls"
	"github.com/nafsilabs/massa-go/client/sendoperation/callsc"
	"github.com/nafsilabs/massa-go/client/sendoperation/executesc"
	"github.com/nafsilabs/massa-go/client/sendoperation/sellrolls"
	"github.com/nafsilabs/massa-go/client/sendoperation/transaction"
	"github.com/nafsilabs/massa-go/utils"
	"github.com/nafsilabs/massa-go/wallet"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// ClientConfig holds configuration for the gRPC client
type ClientConfig struct {
	Address        string
	UseTLS         bool
	TLSConfig      *tls.Config
	DefaultTimeout time.Duration
	DialOptions    []grpc.DialOption
	ChainID        utils.NetworkType
	Account        *wallet.Account
}

// MassaClient provides high-level methods for signing and sending
// operations using wallet accounts (automatic signing)
type MassaClient struct {
	conn           *grpc.ClientConn
	PublicClient   pb.PublicServiceClient
	PrivateClient  pb.PrivateServiceClient
	defaultTimeout time.Duration
	ChainID        utils.NetworkType
	Account        *wallet.Account
}

// NewMassaClient creates a new wallet-integrated operation sender
func NewMassaClient(cfg *ClientConfig) (*MassaClient, error) {

	if cfg.DefaultTimeout == 0 {
		cfg.DefaultTimeout = 30 * time.Second
	}

	var opts []grpc.DialOption

	// Add transport credentials
	if cfg.UseTLS {
		if cfg.TLSConfig != nil {
			opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(cfg.TLSConfig)))
		} else {
			opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
		}
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Add any additional dial options
	opts = append(opts, cfg.DialOptions...)

	// Establish connection
	conn, err := grpc.NewClient(cfg.Address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	return &MassaClient{
		conn:           conn,
		PublicClient:   pb.NewPublicServiceClient(conn),
		PrivateClient:  pb.NewPrivateServiceClient(conn),
		defaultTimeout: cfg.DefaultTimeout,
		ChainID:        cfg.ChainID,
		Account:        cfg.Account,
	}, nil

}

func (c *MassaClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
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
func (c *MassaClient) SendOperation(
	ctx context.Context,
	nickname string,
	password string,
	op sendoperation.Operation,
	fee uint64,
) (string, error) {

	// Get blockchain status
	status, err := c.PublicClient.GetStatus(ctx, &pb.GetStatusRequest{})
	if err != nil {
		log.Fatal(err)
	}
	expirely := status.Status.LastExecutedFinalSlot.Period + status.Status.Config.OperationValidityPeriods

	signedOpBytes, err := sendoperation.MakeOperation(expirely, fee, op, c.Account, password, c.ChainID)
	if err != nil {
		return "", fmt.Errorf("failed to create signed operation: %w", err)
	}

	// Send via bidirectional gRPC stream
	stream, err := c.PublicClient.SendOperations(ctx)
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

// SendTransaction is a convenience method for sending a transaction using wallet
func (c *MassaClient) SendTransaction(
	ctx context.Context,
	nickname string,
	password string,
	recipientAddress string,
	amount float64,
	fee float64,
) (string, error) {
	amountInt := utils.FromMAS(amount)
	feeInt := utils.FromMAS(fee)
	op, err := transaction.New(recipientAddress, amountInt.Uint64())
	if err != nil {
		return "", fmt.Errorf("failed to create transaction operation: %w", err)
	}

	return c.SendOperation(ctx, nickname, password, op, feeInt.Uint64())
}

// BuyRolls is a convenience method for buying rolls using wallet
func (c *MassaClient) BuyRolls(
	ctx context.Context,
	nickname string,
	password string,
	rollCount uint64,
	fee float64,
) (string, error) {
	feeInt := utils.FromMAS(fee)
	op := buyrolls.New(rollCount)
	return c.SendOperation(ctx, nickname, password, op, feeInt.Uint64())
}

// SellRolls is a conven`ience method for selling rolls using wallet
func (c *MassaClient) SellRolls(
	ctx context.Context,
	nickname string,
	password string,
	rollCount uint64,
	fee float64) (string, error) {
	feeInt := utils.FromMAS(fee)
	op := sellrolls.New(rollCount)
	return c.SendOperation(ctx, nickname, password, op, feeInt.Uint64())
}

// SendExecuteSCWithWallet is a convenience method for deploying smart contracts using wallet
func (c *MassaClient) ExecuteSC(
	ctx context.Context,
	nickname string,
	password string,
	bytecode []byte,
	maxGas float64,
	maxCoins float64,
	datastore []byte,
	fee float64,
) (string, error) {
	feeInt := utils.FromMAS(fee)
	coinsInt := utils.FromMAS(maxCoins)
	maxGasInt := utils.FromMAS(maxGas)
	op := executesc.New(bytecode, maxGasInt.Uint64(), coinsInt.Uint64(), datastore)
	return c.SendOperation(ctx, nickname, password, op, feeInt.Uint64())
}

// CallSC is a convenience method for calling smart contracts using wallet
func (c *MassaClient) CallSC(
	ctx context.Context,
	nickname string,
	password string,
	targetAddress string,
	targetFunction string,
	parameters []byte,
	maxGas float64,
	coins float64,
	fee float64,
) (string, error) {
	feeInt := utils.FromMAS(fee)
	coinsInt := utils.FromMAS(coins)
	maxGasInt := utils.FromMAS(maxGas)

	op, err := callsc.New(targetAddress, targetFunction, parameters, maxGasInt.Uint64(), coinsInt.Uint64())
	if err != nil {
		return "", fmt.Errorf("failed to create call smart contract operation: %w", err)
	}
	return c.SendOperation(ctx, nickname, password, op, feeInt.Uint64())
}

// ReadOnlyCallSC performs a read-only smart contract call (no signature required)
// It uses the node JSON-RPC "execute_read_only_call" under the hood and returns the result.
// If callerAddress is empty and the MassaClient was created with a wallet account, the caller
// address will be taken from the wallet account.
func (c *MassaClient) ReadOnlyCallSC(
	ctx context.Context,
	targetAddress string,
	targetFunction string,
	parameters []byte,
	coins float64,
	fee float64,
	callerAddress string,
) (*model.ReadOnlyExecutionOutput, error) {
	// If caller address not provided, try using the wallet account address
	if callerAddress == "" && c.Account != nil && c.Account.Address != nil {
		callerAddress = c.Account.Address.String()
	}

	// Build NativeAmount for coins and fee using FromMAS (nanoMassa) and MANTISSA_SCALE
	coinsNano := utils.FromMAS(coins).Uint64()
	feeNano := utils.FromMAS(fee).Uint64()

	// Prepare ReadOnlyExecutionCall
	call := &model.ReadOnlyExecutionCall{
		MaxGas: sendoperation.MaxGasAllowedCallSC,
		Target: &model.ReadOnlyExecutionCall_FunctionCall{
			FunctionCall: &model.FunctionCall{
				TargetAddress:  targetAddress,
				TargetFunction: targetFunction,
				Parameter:      parameters,
				Coins: &model.NativeAmount{
					Mantissa: coinsNano,
					Scale:    MANTISSA_SCALE,
				},
			},
		},
		Fee: &model.NativeAmount{
			Mantissa: feeNano,
			Scale:    MANTISSA_SCALE,
		},
	}

	// Optionally set caller address
	if callerAddress != "" {
		call.CallerAddress = &wrapperspb.StringValue{Value: callerAddress}
	}

	// Execute gRPC ReadOnly call
	req := &pb.ExecuteReadOnlyCallRequest{Call: call}

	resp, err := c.PublicClient.ExecuteReadOnlyCall(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ExecuteReadOnlyCall gRPC error: %w", err)
	}

	if resp == nil || resp.Output == nil {
		return nil, fmt.Errorf("empty read-only call response")
	}

	return resp.Output, nil
}

// GetDatastoreEntries fetches datastore entries for a given contract address
// and a list of keys. It returns the protobuf model DatastoreEntry objects
// (which contain candidate and final values).
func (c *MassaClient) GetDatastoreEntries(ctx context.Context, address string, keys [][]byte) ([]*model.DatastoreEntry, error) {
	// Build filters for each address-key pair
	filters := make([]*pb.GetDatastoreEntryFilter, len(keys))
	for i, key := range keys {
		filters[i] = &pb.GetDatastoreEntryFilter{
			Filter: &pb.GetDatastoreEntryFilter_AddressKey{
				AddressKey: &model.AddressKeyEntry{
					Address: address,
					Key:     key,
				},
			},
		}
	}

	req := &pb.GetDatastoreEntriesRequest{Filters: filters}

	resp, err := c.PublicClient.GetDatastoreEntries(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetDatastoreEntries gRPC error: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("empty GetDatastoreEntries response")
	}

	return resp.GetDatastoreEntries(), nil
}
