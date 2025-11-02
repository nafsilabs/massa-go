module github.com/nafsilabs/massa-go/client

go 1.25.2

require (
	github.com/btcsuite/btcutil v1.0.2
	github.com/nafsilabs/massa-go/utils v0.0.0
	github.com/nafsilabs/massa-go/wallet v0.0.0
	github.com/shopspring/decimal v1.4.0
	github.com/stretchr/testify v1.11.1
	github.com/ybbus/jsonrpc/v3 v3.1.6
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/klauspost/cpuid/v2 v2.0.12 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/zeebo/blake3 v0.2.4 // indirect
	golang.org/x/crypto v0.42.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/nafsilabs/massa-go/wallet => ../wallet

replace github.com/nafsilabs/massa-go/utils => ../utils

// Map the upstream massa module imported by generated protos to our local generated tree
replace github.com/massalabs/massa => ./proto/massa
