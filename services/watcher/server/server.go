package server

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"

	pb "../proto"
)

type WatcherServer struct {
	GethAddress string
	GethClient *ethclient.Client
}

func (s *WatcherServer) WatchTransaction(ctx context.Context, in *pb.WatchTransactionRequest) (*pb.WatchTransactionResponse, error) {
	log.Printf("Received: %v", in.TxHash)
	tx, isPending, err:= s.GethClient.TransactionByHash(context.Background(), common.HexToHash(in.TxHash))
	if err != nil {
		//TODO do a thing
	}
	return &pb.WatchTransactionResponse{ TxStatus: tx.Value().String(), IsPending: isPending }, nil
}

func (s *WatcherServer) Init() error {
	client, err := ethclient.Dial(s.GethAddress)

	if err != nil {
		return err
	}

	s.GethClient = client

	return nil
}