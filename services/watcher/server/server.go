package server

import (
	pb "../proto"
	"../watching"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"strings"
)

type WatcherServer struct {
	GethAddress string
	MakerEndpoint string
	GethClient *ethclient.Client
	TxWatching *watching.TxWatching
}

func (s *WatcherServer) WatchTransaction(ctx context.Context, in *pb.WatchTransactionRequest) (*pb.WatchTransactionResponse, error) {
	log.Printf("Received: %v", in.TxHash)
	txHash := common.HexToHash(strings.TrimSpace(in.TxHash))

	_, isPending, err:= s.GethClient.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Println("Transaction lookup failure", err)
		return nil, err
		//TODO do a thing
	}

	isWatched := false

	if isPending {
		for _, currentHash := range s.TxWatching.WatchedTransactions {
			if currentHash.TxHash == txHash {
				isWatched = true
			}
		}
		if !isWatched {
			s.TxWatching.WatchedTransactions = append(s.TxWatching.WatchedTransactions, watching.WatchedTransaction{TxHash: txHash, QuoteId: in.QuoteId})
		}
	} else {
		log.Println("Transaction mined")
		// TODO: if pending start listener if not make calls
	}

	return &pb.WatchTransactionResponse{ IsPending: isPending, TxHash: txHash.String(), QuoteId: in.QuoteId }, nil
}

func (s *WatcherServer) Init() error {
	client, err := ethclient.Dial(s.GethAddress)

	if err != nil {
		return fmt.Errorf("failed to connect to rpc" + err.Error())
	}

	txWatching := watching.TxWatching{
		EthClient: client,
		MakerEndpoint: s.MakerEndpoint,
	}

	err = txWatching.Init()

	if err != nil {
		return fmt.Errorf("failed to start tx watching" + err.Error())
	}

	s.GethClient = client
	s.TxWatching = &txWatching

	return nil
}