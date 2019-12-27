// +build maker

package main

import (
	"context"
	"fmt"
	"net"
	"strconv"

	pb "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/logging"
	"google.golang.org/grpc"
)

type MakerServer struct {
}

func (ms *MakerServer) GetQuote(c context.Context, r *pb.GetQuoteRequest) (*pb.GetQuoteResponse, error) {
	return nil, nil
}

func (ms *MakerServer) CheckQuote(c context.Context, r *pb.CheckQuoteRequest) (*pb.CheckQuoteResponse, error) {
	return nil, nil
}

func (ms *MakerServer) OrderStatusUpdate(c context.Context, r *pb.OrderStatusUpdateRequest) (*pb.OrderStatusUpdateResponse, error) {
	logging.Info("Received update for ", r.QuoteId, " Status: ", r.Status)
	return &pb.OrderStatusUpdateResponse{ Status: 1 }, nil
}

func main() {
	logging.Info("Starting Test Maker Endpoint on port 5002")
	lis, err := net.Listen("tcp", ":" + strconv.Itoa(5002))
	if err != nil {
		logging.FatalString(fmt.Sprintf("failed to listen: %v", err))
	}
	s := grpc.NewServer()
	makerServer := MakerServer{}

	pb.RegisterMakerServer(s, &makerServer)
	if err := s.Serve(lis); err != nil {
		logging.FatalString(fmt.Sprintf("failed to serve: %v", err))
	}
}
