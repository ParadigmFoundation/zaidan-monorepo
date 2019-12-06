package server

import (
	pb "../proto"
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

type WatcherServer struct {
	GethAddress string
	GethClient *ethclient.Client
}

func (s *WatcherServer) WatchTransaction(ctx context.Context, in *pb.WatchTransactionRequest) (*pb.WatchTransactionResponse, error) {
	log.Printf("Received: %v", in.TxHash)
	_, isPending, err:= s.GethClient.TransactionByHash(context.Background(), common.HexToHash(in.TxHash))
	if err != nil {
		return nil, err
		//TODO do a thing
	}

	//TODO if pending start listener if not make calls

	return &pb.WatchTransactionResponse{ IsPending: isPending }, nil
}

func (s *WatcherServer) Init() error {
	client, err := ethclient.Dial(s.GethAddress)

	if err != nil {
		return err
	}

	s.GethClient = client

	return nil
}