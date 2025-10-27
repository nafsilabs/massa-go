module github.com/nafsilabs/massa-go/examples/clientexample

go 1.25.2

require github.com/nafsilabs/massa-go/client v0.0.0

require (
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/ybbus/jsonrpc/v3 v3.1.6 // indirect
)

replace github.com/nafsilabs/massa-go/client => ../../client
