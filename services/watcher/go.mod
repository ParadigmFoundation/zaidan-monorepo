module github.com/ParadigmFoundation/zaidan-monorepo/services/watcher

go 1.13

require (
	github.com/ParadigmFoundation/zaidan-monorepo/lib/go v0.0.0-00010101000000-000000000000
	github.com/ParadigmFoundation/zaidan-monorepo/services/dealer v0.0.0-00010101000000-000000000000 // indirect
	github.com/ethereum/go-ethereum v1.9.9
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.4.0
	google.golang.org/grpc v1.25.1
)

replace github.com/ParadigmFoundation/zaidan-monorepo/lib/go => ../../lib/go

replace github.com/ParadigmFoundation/zaidan-monorepo/services/dealer => ../dealer
