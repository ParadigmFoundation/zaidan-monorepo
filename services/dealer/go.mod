module github.com/ParadigmFoundation/zaidan-monorepo/services/dealer

go 1.13

require (
	github.com/ParadigmFoundation/zaidan-monorepo/lib/go v0.0.0-00010101000000-000000000000
	github.com/ethereum/go-ethereum v1.9.9
	github.com/go-sql-driver/mysql v1.4.1 // indirect
	github.com/gogo/protobuf v1.1.1
	github.com/google/uuid v1.1.1
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.1.1
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/mattn/go-sqlite3 v1.11.0
	github.com/peterbourgon/ff v1.6.0
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/stretchr/testify v1.4.0
	google.golang.org/appengine v1.6.1 // indirect
	google.golang.org/genproto v0.0.0-20190620144150-6af8c5fc6601 // indirect
	google.golang.org/grpc v1.23.1
)

replace github.com/ParadigmFoundation/zaidan-monorepo/lib/go => ../../lib/go
