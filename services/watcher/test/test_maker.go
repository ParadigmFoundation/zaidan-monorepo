// +build maker

package main

import (
	"context"
	"net"
	"strconv"

	pb "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/sirupsen/logrus"
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
	logrus.Infof("Received update for %s Status: %d", r.QuoteId, r.Status)
	return &pb.OrderStatusUpdateResponse{Status: 1}, nil
}

func main() {
	logrus.Info("Starting Test Maker Endpoint on port 5002")
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(5002))
	if err != nil {
		logrus.Fatal("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	makerServer := MakerServer{}

	pb.RegisterMakerServer(s, &makerServer)
	if err := s.Serve(lis); err != nil {
		logrus.Fatal("failed to serve: %v", err)
	}
}
