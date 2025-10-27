module github.com/nafsilabs/massa-go/examples/sc_update_age

go 1.25.2

require github.com/nafsilabs/massa-go/client v0.0.0

require (
	github.com/nafsilabs/massa-go/wallet v0.0.0 // indirect
	github.com/nafsilabs/massa-go/utils v0.0.0 // indirect
)

replace github.com/nafsilabs/massa-go/client => ../../client

replace github.com/nafsilabs/massa-go/wallet => ../../wallet

replace github.com/nafsilabs/massa-go/utils => ../../utils
