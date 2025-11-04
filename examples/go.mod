module github.com/nafsilabs/massa-go/examples

go 1.25.2

replace github.com/nafsilabs/massa-go/client => ../client

replace github.com/nafsilabs/massa-go/wallet => ../wallet

replace github.com/nafsilabs/massa-go/utils => ../utils

require github.com/nafsilabs/massa-go/client v0.0.0

require (
	github.com/btcsuite/btcutil v1.0.2 // indirect
	github.com/klauspost/cpuid/v2 v2.0.12 // indirect
	github.com/nafsilabs/massa-go/utils v0.0.0 // indirect
	github.com/nafsilabs/massa-go/wallet v0.0.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/ybbus/jsonrpc/v3 v3.1.6 // indirect
	github.com/zeebo/blake3 v0.2.4 // indirect
	golang.org/x/crypto v0.42.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251014184007-4626949a642f // indirect
	google.golang.org/grpc v1.76.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
