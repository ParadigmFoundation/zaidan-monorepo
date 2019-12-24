module github.com/ParadigmFoundation/zaidan-monorepo/services/dealer

go 1.13

require (
	github.com/ParadigmFoundation/zaidan-monorepo/lib/go v0.0.0-00010101000000-000000000000
	github.com/ethereum/go-ethereum v1.9.9
	github.com/go-sql-driver/mysql v1.4.1 // indirect
	github.com/gogo/protobuf v1.2.0
	github.com/google/uuid v1.0.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/levenlabs/golib v0.0.0-20180911183212-0f8974794783
	github.com/lib/pq v1.1.1
	github.com/mattn/go-sqlite3 v1.11.0
	github.com/stretchr/testify v1.4.0
	google.golang.org/genproto v0.0.0-20190404172233-64821d5d2107 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
)

replace github.com/ParadigmFoundation/zaidan-monorepo/lib/go => ../../lib/go
