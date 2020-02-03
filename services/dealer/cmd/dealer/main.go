package main

import (
	"log"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/rpc"
)

func main() {
	log.Fatal(rpc.StartServer())
}
