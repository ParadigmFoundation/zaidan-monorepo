package grpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/obm"
)

type Client struct {
	obm.OrderBookManagerClient
}

func NewClient(ctx context.Context, addr string) (*Client, error) {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{
		obm.NewOrderBookManagerClient(conn),
	}, nil
}
