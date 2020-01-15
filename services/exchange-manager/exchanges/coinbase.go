package exchanges

import (
	"context"
	"strings"

	"github.com/golang/protobuf/ptypes/empty"
	coinbasepro "github.com/preichenberger/go-coinbasepro/v2"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/em"
)

var _ em.Exchange = &Coinbase{}

type Coinbase struct {
	client *coinbasepro.Client
}

func NewCoinbase(cfg *coinbasepro.ClientConfig) *Coinbase {
	client := coinbasepro.NewClient()
	client.UpdateConfig(cfg)
	return &Coinbase{client: client}
}

func (cb *Coinbase) Client() *coinbasepro.Client { return cb.client }

func (cb *Coinbase) CreateOrder(ctx context.Context, req *types.ExchangeOrder) error {
	order, err := cb.client.CreateOrder(&coinbasepro.Order{
		Side:      strings.ToLower(req.Side.String()),
		ProductID: cb.fromSym(req.Symbol),
		Price:     req.Price,
		Size:      req.Amount,
	})
	if err != nil {
		return err
	}

	req.Id = order.ID
	return nil
}

func (cb *Coinbase) GetOrder(ctx context.Context, id string) (*types.ExchangeOrder, error) {
	order, err := cb.client.GetOrder(id)
	if err != nil {
		return nil, err
	}

	resp := &types.ExchangeOrder{
		Price:  order.Price,
		Symbol: cb.toSym(order.ProductID),
		Amount: order.Size,
		//	Side:   types.ExchangeOrder_BUY,
	}
	return resp, nil
}

func (cb *Coinbase) CancelOrder(ctx context.Context, id string) (*empty.Empty, error) {
	return nil, cb.client.CancelOrder(id)
}

// fromSym convert a symbol to use coinbase notation, this is:
// from `BTC/USD` to `BTC-USD`
func (*Coinbase) fromSym(s string) string { return strings.Replace(s, "/", "-", 1) }

// toSym convert a coinbase symbol into our own symbol notation, this is:
// from `BTC-USD` to `BTC/USD`
func (*Coinbase) toSym(s string) string { return strings.Replace(s, "-", "/", 1) }
