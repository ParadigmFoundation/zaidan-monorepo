module github.com/ParadigmFoundation/zaidan-monorepo/services/dealer

go 1.13

require (
	github.com/0xProject/0x-mesh v0.0.0-20191204233214-2a293766deaa
	github.com/ParadigmFoundation/zaidan-monorepo/lib/go v0.0.0-00010101000000-000000000000
	github.com/ethereum/go-ethereum v1.9.9
	github.com/golang/protobuf v1.3.2
	github.com/google/uuid v1.1.1
	github.com/levenlabs/golib v0.0.0-20180911183212-0f8974794783
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/stretchr/testify v1.4.0
	google.golang.org/grpc v1.25.1
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
)

replace github.com/ParadigmFoundation/zaidan-monorepo/lib/go => ../../lib/go
