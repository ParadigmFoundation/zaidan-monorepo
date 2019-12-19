module github.com/ParadigmFoundation/zaidan-monorepo/services/hot-wallet

go 1.13

require (
	github.com/0xProject/0x-mesh v0.0.0-20191212004844-e881b7dcd31a
	github.com/ParadigmFoundation/zaidan-monorepo/lib/go v0.0.0-00010101000000-000000000000
	github.com/albrow/stringset v2.1.0+incompatible // indirect
	github.com/benbjohnson/clock v1.0.0 // indirect
	github.com/ethereum/go-ethereum v1.9.9
	github.com/gogo/protobuf v1.1.1
	github.com/prometheus/tsdb v0.7.1 // indirect
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/spf13/pflag v1.0.5
	golang.org/x/sys v0.0.0-20191210023423-ac6580df4449 // indirect
	google.golang.org/grpc v1.25.1
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

replace github.com/ParadigmFoundation/zaidan-monorepo/lib/go => ../../lib/go
