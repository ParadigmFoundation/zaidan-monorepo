module github.com/ParadigmFoundation/zaidan-monorepo/services/obm

go 1.13

require (
	github.com/ParadigmFoundation/go-logrus-bugsnag v0.0.0-20200227164141-fdbc509dbe35
	github.com/ParadigmFoundation/zaidan-monorepo/lib/go v0.0.0-00010101000000-000000000000
	github.com/ParadigmFoundation/zaidan-monorepo/services/dealer v0.0.0-00010101000000-000000000000 // indirect
	github.com/adshao/go-binance v0.0.0-20191107145944-a468f0b0c2f0
	github.com/bitbandi/go-hitbtc v0.0.0-20190201230334-2adae5a2f724
	github.com/emirpasic/gods v1.12.0
	github.com/golang/protobuf v1.3.2
	github.com/gorilla/websocket v1.4.1
	github.com/guptarohit/asciigraph v0.4.1
	github.com/juju/errors v0.0.0-20190930114154-d42613fe1ab9 // indirect
	github.com/preichenberger/go-coinbasepro/v2 v2.0.4
	github.com/sourcegraph/jsonrpc2 v0.0.0-20191113080033-cee7209801bf // indirect
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413
	google.golang.org/grpc v1.27.1
	gotest.tools v2.2.0+incompatible
)

replace github.com/ParadigmFoundation/zaidan-monorepo/lib/go => ../../lib/go

replace github.com/ParadigmFoundation/zaidan-monorepo/services/dealer => ../dealer
