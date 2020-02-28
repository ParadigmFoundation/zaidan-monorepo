module github.com/ParadigmFoundation/zaidan-monorepo/services/hot-wallet

go 1.13

require (
	github.com/0xProject/0x-mesh v0.0.0-20191212004844-e881b7dcd31a
	github.com/ParadigmFoundation/go-logrus-bugsnag v0.0.0-20200227164141-fdbc509dbe35
	github.com/ParadigmFoundation/zaidan-monorepo/lib/go v0.0.0-00010101000000-000000000000
	github.com/ParadigmFoundation/zaidan-monorepo/services/dealer v0.0.0-00010101000000-000000000000 // indirect
	github.com/VividCortex/gohistogram v1.0.0 // indirect
	github.com/caarlos0/env/v6 v6.1.0
	github.com/ethereum/go-ethereum v1.9.11
	github.com/gogo/protobuf v1.2.0
	github.com/golang/protobuf v1.3.2
	github.com/prometheus/tsdb v0.7.1 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	golang.org/x/sys v0.0.0-20191210023423-ac6580df4449 // indirect
	google.golang.org/grpc v1.27.1
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

replace github.com/ParadigmFoundation/zaidan-monorepo/lib/go => ../../lib/go

replace github.com/ParadigmFoundation/zaidan-monorepo/services/dealer => ../dealer
