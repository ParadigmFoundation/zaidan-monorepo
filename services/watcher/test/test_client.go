// +build client

package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	pb "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/logger"
	"google.golang.org/grpc"
)

var log = logger.New("app")

func main() {

	// Set up a connection to the Server.
	const address = "localhost:5001"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(fmt.Errorf("did not connect: %v", err))
	}
	defer conn.Close()
	c := pb.NewWatcherClient(conn)

	reader := bufio.NewReader(os.Stdin)

	for {
		log.Info("Enter txHash:")
		txHash, _ := reader.ReadString('\n')

		resp, err := c.WatchTransaction(
			context.Background(),
			&pb.WatchTransactionRequest{
				QuoteId: "Random from test_client",
				TxHash: strings.Replace(txHash, "\n", "", -1),
				StatusUrls: []string{"http://localhost:5002"},
			},
		)

		if err != nil {
			log.Error(fmt.Errorf("Error: %s", err))
		} else {
			log.Info("Call succeeded: { txHash: ", resp.TxHash, ", quoteId: ", resp.QuoteId, ", isPending: ", fmt.Sprint(resp.IsPending), ", txStatus: ", resp.TxStatus, " }")
		}
	}

}
