module github.com/nafsilabs/massa-go

go 1.25.2

require (
	github.com/btcsuite/btcutil v1.0.2
	github.com/json-iterator/go v1.1.12
	github.com/shopspring/decimal v1.4.0
	github.com/stretchr/testify v1.11.1
	github.com/ybbus/jsonrpc/v3 v3.1.6
	github.com/zeebo/blake3 v0.2.4
	golang.org/x/crypto v0.42.0
	google.golang.org/genproto/googleapis/api v0.0.0-20251029180050-ab9386a59fda
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/klauspost/cpuid/v2 v2.0.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251014184007-4626949a642f // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/massalabs/massa => ./client/proto/massa
