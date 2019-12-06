package main

import (
	pb "../proto"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
)

func main() {

	// Set up a connection to the Server.
	const address = "localhost:5001"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewWatcherClient(conn)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Enter txHash:")
		txHash, _ := reader.ReadString('\n')

		resp, err := c.WatchTransaction(context.Background(), &pb.WatchTransactionRequest{ TxHash: strings.Replace(txHash, "\n", "", -1)})

		if (err != nil) {
			log.Print("Error: ", err)
		} else {
			log.Print("Success: ", fmt.Sprint(resp.IsPending))
		}
	}

}