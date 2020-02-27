package grpc

import (
	"context"
	"github.com/ParadigmFoundation/go-logrus-bugsnag/logger"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/store/sql"
	"net"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	ggrpc "google.golang.org/grpc"
)

type TransactionStatusServer struct {
	log   *logger.Logger
	store *sql.Store
}

func (oss *TransactionStatusServer) TransactionStatusUpdate(c context.Context, r *grpc.TransactionStatusUpdateRequest) (*grpc.TransactionStatusUpdateResponse, error) {
	oss.log.Info("Received update for QuoteId: ", r.QuoteId, ", TransactionHash: ", r.TxHash, " with Status: ", r.Status)
	if err := oss.store.UpdateTradeStatus(r.QuoteId, grpc.Trade_Status(r.Status + 1)); err != nil {
		return &grpc.TransactionStatusUpdateResponse{ Status: 0 }, err
	}
	return &grpc.TransactionStatusUpdateResponse{ Status: 1 }, nil
}

func CreateAndListen(store *sql.Store, binding string) {
	log := logger.New("transaction_status_grpc")

	ln, err := net.Listen("tcp", binding)
	if err != nil {
		log.Fatal("failed to listen: ", err)
	}
	ggrpcServer := ggrpc.NewServer()

	grpc.RegisterTransactionStatusServer(ggrpcServer, &TransactionStatusServer{
		log: log,
		store: store,
	})
	if err := ggrpcServer.Serve(ln); err != nil {
		log.WithError(err).Fatal("failed to serve TransactionStatus:")
	}
}
