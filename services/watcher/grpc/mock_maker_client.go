package grpc

import (
	"context"

	pb "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"google.golang.org/grpc"
)

type MockMakerClient struct {

}


func (MockMakerClient) GetQuote(ctx context.Context, in *pb.GetQuoteRequest, opts ...grpc.CallOption) (*pb.GetQuoteResponse, error) {
	return &pb.GetQuoteResponse{}, nil
}

func (MockMakerClient) CheckQuote(ctx context.Context, in *pb.CheckQuoteRequest, opts ...grpc.CallOption) (*pb.CheckQuoteResponse, error) {
	return &pb.CheckQuoteResponse{}, nil
}

func (MockMakerClient) OrderStatusUpdate(ctx context.Context, in *pb.OrderStatusUpdateRequest, opts ...grpc.CallOption) (*pb.OrderStatusUpdateResponse, error) {
	return &pb.OrderStatusUpdateResponse{Status: 1}, nil
}

