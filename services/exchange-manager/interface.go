package em

import (
	"context"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/golang/protobuf/ptypes/empty"
)

type Exchange interface {
	CreateOrder(context.Context, *types.ExchangeOrder) error
	GetOrder(context.Context, string) (*types.ExchangeOrder, error)
	CancelOrder(context.Context, string) (*empty.Empty, error)
}
