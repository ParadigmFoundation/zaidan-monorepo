package grpc

import (
	"context"

	"google.golang.org/grpc"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

type Client struct {
	types.OrderBookManagerClient
}

func NewClient(ctx context.Context, addr string) (*Client, error) {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{
		types.NewOrderBookManagerClient(conn),
	}, nil
}
