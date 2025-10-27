package client

import (
	"fmt"
	"regexp"

	"github.com/nafsilabs/massa-go/client/sendoperation/executesc"
	"github.com/nafsilabs/massa-go/wallet"
)

// DeploySC deploys a smart contract on the blockchain.
// The smart contract is deployed with the given account nickname.

func DeploySC(
	c *Client,
	nickname string,
	maxGas uint64,
	maxCoins uint64,
	coins uint64,
	fees uint64,
	expiry uint64,
	parameters []byte,
	smartContractByteCode []byte,
	deployerByteCode []byte,
	account *wallet.Account,
	password string,
	description string,
) (*OperationResponse, []Event, error) {

	contract := ContractDatastore{
		ByteCode: smartContractByteCode,
		Args:     parameters,
		Coins:    coins,
	}

	dataStore, err := populateDatastore(contract)
	if err != nil {
		return nil, nil, fmt.Errorf("populating datastore: %w", err)
	}

	exeSCOperation := executesc.New(
		deployerByteCode,
		maxGas,
		maxCoins,
		dataStore,
	)

	operationResponse, err := Call(
		c,
		expiry,
		fees,
		exeSCOperation,
		account,
		password,
		description,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("calling executeSC: %w", err)
	}

	events, err := ListenEvents(c, nil, nil, nil, &operationResponse.OperationID, nil, true)
	if err != nil {
		return nil, nil, fmt.Errorf("listening events for opId at %s : %w", operationResponse.OperationID, err)
	}

	return operationResponse, events, nil
}

func FindDeployedAddress(events []Event) (string, bool) {
	pattern := "Contract deployed at address: (.+)"
	re := regexp.MustCompile(pattern)

	for _, event := range events {
		match := re.FindStringSubmatch(event.Data)
		if len(match) > 1 {
			return match[1], true
		}
	}

	return "", false
}
