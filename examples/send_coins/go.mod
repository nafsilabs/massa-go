module github.com/nafsilabs/massa-go/examples/send_coins

go 1.25.2

require github.com/nafsilabs/massa-go/client v0.0.0

require (
	github.com/k0kubun/pp v3.0.1+incompatible // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/nafsilabs/massa-go/utils v0.0.0 // indirect
	github.com/nafsilabs/massa-go/wallet v0.0.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
)

replace github.com/nafsilabs/massa-go/client => ../../client

replace github.com/nafsilabs/massa-go/wallet => ../../wallet

replace github.com/nafsilabs/massa-go/utils => ../../utils
