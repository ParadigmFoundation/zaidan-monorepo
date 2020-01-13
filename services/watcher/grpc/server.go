package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	pb "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/logging"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/eth"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/watcher/watching"
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

type WatcherServer struct {
	TxWatching *watching.TxWatching
	grpc *grpc.Server
}

func NewServer(txWatching *watching.TxWatching) *WatcherServer {
	watcherServer := &WatcherServer{
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
	logging.Info("Watcher listening on port", port)

	return s.grpc.Serve(lis)
}

func (s *WatcherServer) Stop() { s.grpc.GracefulStop() }

func (s *WatcherServer) WatchTransaction(ctx context.Context, in *pb.WatchTransactionRequest) (*pb.WatchTransactionResponse, error) {
	trimmed := strings.TrimSpace(in.TxHash)
	txHash := common.HexToHash(trimmed)
	updateUrls := in.StatusUrls
	if !strings.EqualFold(txHash.String(), trimmed) || len(txHash.String()) != 66 {
		return nil, errors.New("invalid txHash")
	}

	logging.Info(fmt.Sprintf("Received: %v", in.TxHash))
	s.TxWatching.Lock()
	isPending, status, err := getTransactionInfo(ctx, s, txHash)
	if err != nil {
		logging.SafeError(fmt.Errorf("transaction lookup failure: %s", err.Error()))
		return nil, grpcStatus.Error(grpcCodes.Internal, fmt.Sprintf("transaction lookup failure: %s", err.Error()))
	}

	_, isWatched := s.TxWatching.IsWatched(txHash)
	if isPending && !isWatched {
		logging.Info(fmt.Sprintf("Now watching: %v", in.TxHash))
		s.TxWatching.Watch(txHash, in.QuoteId, updateUrls)
		isWatched = true
	} else if !isPending {
		if isWatched {
			logging.Info("No longer watching previously mined transaction", txHash.String())
			s.TxWatching.RequestRemoval(txHash)
		}

		logging.Info("Transaction", txHash.String(), "mined and does not need to be watched notifying maker")

		for _, url := range updateUrls {
			conn, err := grpc.Dial(url, grpc.WithInsecure())
			if err != nil {
				logging.Fatal(fmt.Errorf("failed to connect maker client: %v", err))
			}

			_, _ = pb.NewOrderStatusClient(conn).OrderStatusUpdate(ctx, &pb.OrderStatusUpdateRequest{
				QuoteId: in.QuoteId,
				Status:  status,
			})
		}
	}
	s.TxWatching.Unlock()

	return &pb.WatchTransactionResponse{ IsWatched: isWatched, IsPending: isPending, TxStatus: status, TxHash: txHash.String(), QuoteId: in.QuoteId }, nil
}

func getTransactionInfo(c context.Context, s *WatcherServer, txHash common.Hash) (bool, uint32, error) {
	_, isPending, err:= eth.Client.TransactionByHash(c, txHash)
	if err != nil {
		return false, 0, err
	}


	if isPending {
		return isPending, 0, nil
	}

	receipt, err := eth.Client.TransactionReceipt(c, txHash)
	if err != nil {
		return isPending, 0, err
	}

	return isPending, uint32(receipt.Status), nil
}