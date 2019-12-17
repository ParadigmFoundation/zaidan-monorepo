package grpc

import (
	"context"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"google.golang.org/grpc"
)

// Client is a gRPC client for the hot wallet service
type Client struct {
	types.HotWalletClient
}

// NewClient returns a new gRPC client for the hot wallet
func NewClient(ctx context.Context, addr string) (*Client, error) {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{
		types.NewHotWalletClient(conn),
	}, nil
}
