package grpc

import (
	"context"

	pb "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"google.golang.org/grpc"
)

type MockTransactionStatusServer struct {

}

func (MockTransactionStatusServer) TransactionStatusUpdate(ctx context.Context, in *pb.TransactionStatusUpdateRequest, opts ...grpc.CallOption) (*pb.TransactionStatusUpdateResponse, error) {
	return &pb.TransactionStatusUpdateResponse{Status: 1}, nil
}

