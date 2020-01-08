module github.com/ParadigmFoundation/zaidan-monorepo/lib/go

go 1.13

require (
	github.com/0xProject/0x-mesh v0.0.0-20191204233214-2a293766deaa
	github.com/albrow/stringset v2.1.0+incompatible // indirect
	github.com/benbjohnson/clock v1.0.0 // indirect
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/caarlos0/env/v6 v6.1.0
	github.com/ethereum/go-ethereum v1.9.9
	github.com/gogo/protobuf v1.1.1
	github.com/golang/protobuf v1.3.2
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/miguelmota/go-ethereum-hdwallet v0.0.0-20191015012459-abf3d7f7f00c
	github.com/multiformats/go-multiaddr-dns v0.2.0 // indirect
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413
	google.golang.org/grpc v1.23.1
)

replace github.com/ParadigmFoundation/zaidan-monorepo/services/dealer => ../../services/dealer
