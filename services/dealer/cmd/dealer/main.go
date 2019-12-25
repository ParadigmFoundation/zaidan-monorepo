package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/store/sql"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/peterbourgon/ff"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/rpc"
)

func main() {
	fs := flag.NewFlagSet("dealer", flag.ExitOnError)
	var (
		db  = fs.String("db", "sqlite3", "Database driver [sqlite3|postgres]")
		dsn = fs.String("dsn", ":memory:", "Database's Data Source Name (see driver's doc for details)")
	)
	ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("DEALER"),
	)

	store, err := sql.New(*db, *dsn)
	if err != nil {
		log.Fatal(err)
	}
	_ = store

	server := gethrpc.NewServer()
	service, err := rpc.NewService()
	if err != nil {
		panic(err)
	}

	if err := server.RegisterName("dealer", service); err != nil {
		panic(err)
	}

	if err := http.ListenAndServe("0.0.0.0:8000", server.WebsocketHandler([]string{"*"})); err != nil {
		log.Fatal(err)
	}
}
