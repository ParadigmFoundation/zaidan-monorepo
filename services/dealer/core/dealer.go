package core

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/store"
)

// DealerConfig defines configuration for a Dealer
type DealerConfig struct {
	MakerBindAddress     string
	HotWalletBindAddress string

	OrderDuration int64 // the number of seconds order should be valid for after initial quote is provided
}

// Dealer is the core dealer service that interacts with other services
type Dealer struct {
	makerClient types.MakerClient
	hwClient    types.HotWalletClient

	orderDuration int64

	db     store.Store
	logger log.Logger
}

// NewDealer creates a new Dealer given ctx context and cfg configuration
func NewDealer(ctx context.Context, cfg DealerConfig) (*Dealer, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	logger := log.New(ctx)

	makerConn, err := grpc.DialContext(ctx, cfg.MakerBindAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	hwConn, err := grpc.DialContext(ctx, cfg.HotWalletBindAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Dealer{
		makerClient:   types.NewMakerClient(makerConn),
		hwClient:      types.NewHotWalletClient(hwConn),
		orderDuration: cfg.OrderDuration,
		logger:        logger,
	}, nil
}

// @todo: hrharder - should eventually return types.Quote
func (d *Dealer) FetchQuote(ctx context.Context, req *types.GetQuoteRequest) (*types.Quote, error) {
	now := time.Now()
	res, err := d.makerClient.GetQuote(ctx, req)
	if err != nil {
		return nil, err
	}

	orderReq := &types.CreateOrderRequest{
		TakerAddress:          req.TakerAsset,
		MakerAssetAddress:     res.MakerAsset,
		TakerAssetAddress:     res.TakerAsset,
		MakerAssetAmount:      res.MakerSize,
		TakerAssetAmount:      res.TakerSize,
		ExpirationTimeSeconds: now.Unix() + d.orderDuration,
	}

	orderRes, err := d.hwClient.CreateOrder(ctx, orderReq)
	if err != nil {
		return nil, err
	}

	quote := &types.Quote{
		QuoteId:           res.QuoteId,
		MakerAssetAddress: res.MakerAsset,
		TakerAssetAddress: res.TakerAsset,
		MakerAssetSize:    res.MakerSize,
		TakerAssetSize:    res.TakerSize,
		Expiration:        res.Expiration,
		ServerTime:        now.UnixNano() / 1e6, // conversion from nanoseconds to milliseconds = ns / 1e6
		OrderHash:         hexutil.Encode(orderRes.Hash),
		Order:             orderRes.Order,
		FillTx:            hexutil.Encode(orderRes.FillTxData),
	}

	return quote, nil
}
