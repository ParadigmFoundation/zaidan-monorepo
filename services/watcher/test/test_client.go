// +build client

package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	pb "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/logging"
	"google.golang.org/grpc"
)

func main() {

	// Set up a connection to the Server.
	const address = "localhost:5001"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logging.FatalString(fmt.Sprintf("did not connect: %v", err))
	}
	defer conn.Close()
	c := pb.NewWatcherClient(conn)

	reader := bufio.NewReader(os.Stdin)

	for {
		logging.Info("Enter txHash:")
		txHash, _ := reader.ReadString('\n')

		resp, err := c.WatchTransaction(context.Background(), &pb.WatchTransactionRequest{QuoteId: "Random from test_client", TxHash: strings.Replace(txHash, "\n", "", -1)})

		if err != nil {
			logging.SafeErrorString(fmt.Sprintf("Error: %s", err))
		} else {
			logging.Info("Call succeeded: { txHash: ", resp.TxHash, ", quoteId: ", resp.QuoteId, ", isPending: ", fmt.Sprint(resp.IsPending), ", txStatus: ", resp.TxStatus, " }")
		}
	}

}
