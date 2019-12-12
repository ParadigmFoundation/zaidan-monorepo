module github.com/ParadigmFoundation/zaidan-monorepo/services/watcher

go 1.13

replace github.com/ParadigmFoundation/zaidan-monorepo/lib/go => ../../lib/go

require (
	github.com/ParadigmFoundation/zaidan-monorepo/lib/go v0.0.0-00010101000000-000000000000
	github.com/ethereum/go-ethereum v1.9.9
	github.com/spf13/cobra v0.0.5 // indirect
	google.golang.org/grpc v1.25.1
)
