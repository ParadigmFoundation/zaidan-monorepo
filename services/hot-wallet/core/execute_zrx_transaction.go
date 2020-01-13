package core

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

// ExecuteZeroExTransaction implements grpc.HotWalletServer
func (hw *HotWallet) ExecuteZeroExTransaction(ctx context.Context, req *grpc.ExecuteZeroExTransactionRequest) (*grpc.ExecuteZeroExTransactionResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method not yet implemented")
}
