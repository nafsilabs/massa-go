module github.com/nafsilabs/massa-go/examples/sc_readonly_call

go 1.25.2

require github.com/nafsilabs/massa-go/client v0.0.0

require (
	github.com/btcsuite/btcutil v1.0.2 // indirect
	github.com/klauspost/cpuid/v2 v2.0.12 // indirect
	github.com/nafsilabs/massa-go/utils v0.0.0 // indirect
	github.com/nafsilabs/massa-go/wallet v0.0.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/ybbus/jsonrpc/v3 v3.1.6 // indirect
	github.com/zeebo/blake3 v0.2.4 // indirect
	golang.org/x/crypto v0.40.0 // indirect
)

replace github.com/nafsilabs/massa-go/client => ../../client

replace github.com/nafsilabs/massa-go/wallet => ../../wallet

replace github.com/nafsilabs/massa-go/utils => ../../utils
