package grpc

import (
	"context"

	"google.golang.org/grpc"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

type MakerClient struct {
	types.MakerClient
}

func NewMakerClient(ctx context.Context, addr string) (*MakerClient, error) {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &MakerClient{types.NewMakerClient(conn)}, nil
}
