package massa

import (
	"github.com/nafsilabs/massa-go/client"
	"github.com/nafsilabs/massa-go/wallet"
)

type ClientConfig = client.ClientConfig
type MassaClient = client.MassaClient
type Wallet = wallet.Wallet
type WalletConfig = wallet.WalletConfig

func NewMassaClient(cfg *ClientConfig) (*MassaClient, error) {
	return client.NewMassaClient(cfg)
}

func NewWallet(config *WalletConfig) *Wallet {
	return wallet.NewWallet(config)
}

func LoadWallet(path string) (*Wallet, error) {
	return wallet.LoadWallet(path)
}
