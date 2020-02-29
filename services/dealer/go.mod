module github.com/ParadigmFoundation/zaidan-monorepo/services/dealer

go 1.13

require (
	github.com/0xProject/0x-mesh v0.0.0-20191204233214-2a293766deaa
	github.com/ParadigmFoundation/go-logrus-bugsnag v0.0.0-20200227164141-fdbc509dbe35
	github.com/ParadigmFoundation/zaidan-monorepo/lib/go v0.0.0-00010101000000-000000000000
	github.com/ethereum/go-ethereum v1.9.11
	github.com/go-sql-driver/mysql v1.4.1 // indirect
	github.com/gogo/protobuf v1.1.1
	github.com/golang/protobuf v1.3.2
	github.com/google/uuid v1.1.1
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.1.1
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/mattn/go-sqlite3 v1.11.0
	github.com/peterbourgon/ff v1.6.0
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0
	google.golang.org/appengine v1.6.1 // indirect
	google.golang.org/grpc v1.27.1
)

replace github.com/ParadigmFoundation/zaidan-monorepo/lib/go => ../../lib/go
