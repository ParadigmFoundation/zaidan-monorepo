package grpc

import (
	"context"

	"google.golang.org/grpc"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

type HotWalletClient struct {
	types.HotWalletClient
}

func NewHotWalletClient(ctx context.Context, addr string) (*HotWalletClient, error) {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &HotWalletClient{types.NewHotWalletClient(conn)}, nil
}
