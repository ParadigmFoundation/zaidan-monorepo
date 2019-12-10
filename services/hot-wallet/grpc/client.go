package grpc

import (
	"context"

	hw "github.com/ParadigmFoundation/zaidan-monorepo/services/hot-wallet"
	"google.golang.org/grpc"
)

// Client is a gRPC client for the hot wallet service
type Client struct {
	hw.HotWalletClient
}

// NewClient returns a new gRPC client for the hot wallet
func NewClient(ctx context.Context, addr string) (*Client, error) {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{
		hw.NewHotWalletClient(conn),
	}, nil
}
