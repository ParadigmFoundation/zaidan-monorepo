package grpc

import (
	"context"

	pb "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"google.golang.org/grpc"
)

type MockOrderStatusServer struct {

}

func (MockOrderStatusServer) OrderStatusUpdate(ctx context.Context, in *pb.OrderStatusUpdateRequest, opts ...grpc.CallOption) (*pb.OrderStatusUpdateResponse, error) {
	return &pb.OrderStatusUpdateResponse{Status: 1}, nil
}

