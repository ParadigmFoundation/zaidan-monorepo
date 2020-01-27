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

type OrderStatusServer struct {
}

func (oss *OrderStatusServer) OrderStatusUpdate(c context.Context, r *pb.OrderStatusUpdateRequest) (*pb.OrderStatusUpdateResponse, error) {
	logrus.Info("Received update for ", r.QuoteId, " Status: ", r.Status)
	return &pb.OrderStatusUpdateResponse{ Status: 1 }, nil
}

func main() {
	logrus.Info("Starting Test Maker Endpoint on port 5002")
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(5002))
	if err != nil {
		logrus.Fatal("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	orderStatusServer := OrderStatusServer{}

	pb.RegisterOrderStatusServer(s, &orderStatusServer)
	if err := s.Serve(lis); err != nil {
		logrus.Fatal("failed to serve: %v", err)
	}
}
