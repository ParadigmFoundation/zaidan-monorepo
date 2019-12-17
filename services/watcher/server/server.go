package server

import (
	"context"
	"fmt"
	"log"
	"strings"

	pb "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/watching"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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
	//TODO: validate transaction hash

	isPending, status, err := getTransactionInfo(ctx, s, txHash)
	if err != nil {
		log.Println("Transaction lookup failure", err)
		return nil, err
		//TODO do a thing
	}

	_/*watchedTx*/, isWatched := s.TxWatching.WatchedTransactions[txHash] //TODO use watched tx?

	if isPending && !isWatched {
		s.TxWatching.WatchedTransactions[txHash] = watching.WatchedTransaction{TxHash: txHash, QuoteId: in.QuoteId}

func getTransactionInfo(c context.Context, s *WatcherServer, txHash common.Hash) (bool, uint32, error) {
	_, isPending, err:= s.GethClient.TransactionByHash(c, txHash)
	if err != nil {
		return false, 0, err
	}


	if isPending {
		return isPending, 0, nil
	}

	receipt, err := s.GethClient.TransactionReceipt(c, txHash)
	if err != nil {
		return isPending, 0, err
	}

	return isPending, uint32(receipt.Status), nil
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