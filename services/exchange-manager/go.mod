module github.com/ParadigmFoundation/zaidan-monorepo/services/em

go 1.13

require (
	github.com/ParadigmFoundation/zaidan-monorepo/lib/go v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.3.2
	github.com/google/uuid v1.0.0
	github.com/peterbourgon/ff v1.7.0
	github.com/preichenberger/go-coinbasepro/v2 v2.0.4
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.4.0
	google.golang.org/grpc v1.26.0
)

replace github.com/ParadigmFoundation/zaidan-monorepo/lib/go => ../../lib/go
