module github.com/ParadigmFoundation/zaidan-monorepo/services/hot-wallet

go 1.13

require (
	github.com/0xProject/0x-mesh v0.0.0-20191212004844-e881b7dcd31a
	github.com/ParadigmFoundation/zaidan-monorepo/lib/go v0.0.0-00010101000000-000000000000
	github.com/VividCortex/gohistogram v1.0.0 // indirect
	github.com/albrow/stringset v2.1.0+incompatible // indirect
	github.com/benbjohnson/clock v1.0.0 // indirect
	github.com/caarlos0/env/v6 v6.1.0
	github.com/ethereum/go-ethereum v1.9.9
	github.com/golang/protobuf v1.3.2
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/stretchr/testify v1.4.0
	google.golang.org/grpc v1.25.1
)

replace github.com/ParadigmFoundation/zaidan-monorepo/lib/go => ../../lib/go
