// +build maker

package main

import (
	"context"
	"net"
	"strconv"

	pb "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/logger"
	"google.golang.org/grpc"
)

var log = logger.New("app")

type TransactionStatusServer struct {
}

func (oss *TransactionStatusServer) TransactionStatusUpdate(c context.Context, r *pb.TransactionStatusUpdateRequest) (*pb.TransactionStatusUpdateResponse, error) {
	log.Info("Received update for QuoteId: ", r.QuoteId, ", TransactionHash: ", r.TxHash, " with Status: ", r.Status)
	return &pb.TransactionStatusUpdateResponse{ Status: 1 }, nil
}

func main() {
	log.Info("Starting Test Maker Endpoint on port 5002")
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(5002))
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	transactionStatusServer := TransactionStatusServer{}

	pb.RegisterTransactionStatusServer(s, &transactionStatusServer)
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve: %v", err)
	}
}
