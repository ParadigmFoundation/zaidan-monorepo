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

type OrderStatusServer struct {
}

func (oss *OrderStatusServer) OrderStatusUpdate(c context.Context, r *pb.OrderStatusUpdateRequest) (*pb.OrderStatusUpdateResponse, error) {
	logging.Info("Received update for ", r.QuoteId, " Status: ", r.Status)
	return &pb.OrderStatusUpdateResponse{ Status: 1 }, nil
}

func main() {
	logging.Info("Starting Test Maker Endpoint on port 5002")
	lis, err := net.Listen("tcp", ":" + strconv.Itoa(5002))
	if err != nil {
		logging.Fatal(fmt.Errorf("failed to listen: %v", err))
	}
	s := grpc.NewServer()
	orderStatusServer := OrderStatusServer{}

	pb.RegisterOrderStatusServer(s, &orderStatusServer)
	if err := s.Serve(lis); err != nil {
		logging.Fatal(fmt.Errorf("failed to serve: %v", err))
	}
}
