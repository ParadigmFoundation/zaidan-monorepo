package grpc

import (
	"context"
	"log"
	"net"
	"strings"

	pb "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/eth"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/watching"
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc"
)

type WatcherServer struct {
	EthToolkit *eth.EthereumToolkit
	TxWatching *watching.TxWatching
	grpc *grpc.Server
}

func NewServer(ethToolkit *eth.EthereumToolkit, txWatching *watching.TxWatching) *WatcherServer {
	watcherServer := &WatcherServer{
		EthToolkit: ethToolkit,
		TxWatching: txWatching,
		grpc:  grpc.NewServer(),
	}

	pb.RegisterWatcherServer(watcherServer.grpc, watcherServer)

	return watcherServer
}

func (s *WatcherServer) Listen(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	log.Println("Watcher listening on port", port)

	return s.grpc.Serve(lis)
}

func (s *WatcherServer) Stop() { s.grpc.GracefulStop() }

func (s *WatcherServer) WatchTransaction(ctx context.Context, in *pb.WatchTransactionRequest) (*pb.WatchTransactionResponse, error) {
	log.Printf("Received: %v", in.TxHash)
	txHash := common.HexToHash(strings.TrimSpace(in.TxHash))
	//TODO: validate transaction hash

	s.TxWatching.Lock()
	isPending, status, err := getTransactionInfo(ctx, s, txHash)
	if err != nil {
		log.Println("Transaction lookup failure", err)
		return nil, err
	}

	_/*watchedTx*/, isWatched := s.TxWatching.IsWatched(txHash) //TODO use watched tx?
	if isPending && !isWatched {
		log.Printf("Now watching: %v", in.TxHash)
		s.TxWatching.Watch(txHash, in.QuoteId)
		isWatched = true
	} else if isPending && isWatched {
		log.Printf("Already watching") //TODO
	} else {
		if isWatched {
			log.Println("This is not pending but is watched so needs to be handled") //TODO
		} else {
			log.Println("Transaction mined and does not need to be watched") // TODO
		}
	}
	s.TxWatching.Unlock()

	return &pb.WatchTransactionResponse{ IsWatched: isWatched, IsPending: isPending, TxStatus: status, TxHash: txHash.String(), QuoteId: in.QuoteId }, nil
}

func getTransactionInfo(c context.Context, s *WatcherServer, txHash common.Hash) (bool, uint32, error) {
	_, isPending, err:= s.EthToolkit.Client.TransactionByHash(c, txHash)
	if err != nil {
		return false, 0, err
	}


	if isPending {
		return isPending, 0, nil
	}

	receipt, err := s.EthToolkit.Client.TransactionReceipt(c, txHash)
	if err != nil {
		return isPending, 0, err
	}

	return isPending, uint32(receipt.Status), nil
}