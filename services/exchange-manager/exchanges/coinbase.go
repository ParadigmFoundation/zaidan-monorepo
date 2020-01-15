package exchanges

import (
	"context"
	"encoding/json"
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

func (cb *Coinbase) GetOrder(ctx context.Context, id string) (*types.ExchangeOrderResponse, error) {
	order, err := cb.client.GetOrder(id)
	if err != nil {
		return nil, err
	}

	return cb.NewOrderResponse(&order)
}

func (cb *Coinbase) GetOpenOrders(ctx context.Context) (*types.ExchangeOrderArrayResponse, error) {
	res := &types.ExchangeOrderArrayResponse{}

	cursor := cb.client.ListOrders()
	for cursor.HasMore {
		var orders []*coinbasepro.Order
		if err := cursor.NextPage(&orders); err != nil {
			return nil, err
		}

		for _, order := range orders {
			newOrder, err := cb.NewOrderResponse(order)
			if err != nil {
				return nil, err
			}
			res.Array = append(res.Array, newOrder)
		}
	}

	return res, nil
}

func (cb *Coinbase) CancelOrder(ctx context.Context, id string) (*empty.Empty, error) {
	return nil, cb.client.CancelOrder(id)
}

// NewOrderResponse converts a coinbase Order type into our type
func (cb *Coinbase) NewOrderResponse(order *coinbasepro.Order) (*types.ExchangeOrderResponse, error) {
	// Convert the side
	var side types.ExchangeOrder_Side
	switch order.Side {
	case "buy":
		side = types.ExchangeOrder_BUY
	case "sell":
		side = types.ExchangeOrder_SELL
	}

	// encode original response
	infoBytes, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	return &types.ExchangeOrderResponse{
		Order: &types.ExchangeOrder{
			Price:  order.Price,
			Symbol: cb.toSym(order.ProductID),
			Amount: order.Size,
			Side:   side,
		},
		Status: &types.ExchangeOrderStatus{
			Timestamp: order.CreatedAt.Time().Unix(),
			Filled:    order.FilledSize,
			Info:      infoBytes,
		},
	}, nil
}

// fromSym convert a symbol to use coinbase notation, this is:
// from `BTC/USD` to `BTC-USD`
func (*Coinbase) fromSym(s string) string { return strings.Replace(s, "/", "-", 1) }

// toSym convert a coinbase symbol into our own symbol notation, this is:
// from `BTC-USD` to `BTC/USD`
func (*Coinbase) toSym(s string) string { return strings.Replace(s, "-", "/", 1) }
