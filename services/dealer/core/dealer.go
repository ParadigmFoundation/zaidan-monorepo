package core

import (
	"context"
	"fmt"
	"time"

	"github.com/0xProject/0x-mesh/zeroex"

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
	WatcherBindAddress   string

	OrderDuration int64 // the number of seconds order should be valid for after initial quote is provided
}

// Dealer is the core dealer service that interacts with other services
type Dealer struct {
	Logger log.Logger

	makerClient   types.MakerClient
	hwClient      types.HotWalletClient
	watcherClient types.WatcherClient

	orderDuration int64

	db store.Store
}

// NewDealer creates a new Dealer given ctx context and cfg configuration
func NewDealer(ctx context.Context, db store.Store, cfg DealerConfig) (*Dealer, error) {
	logger := log.New(ctx)

	makerConn, err := grpc.DialContext(ctx, cfg.MakerBindAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	hwConn, err := grpc.DialContext(ctx, cfg.HotWalletBindAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	watcherConn, err := grpc.DialContext(ctx, cfg.WatcherBindAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Dealer{
		makerClient:   types.NewMakerClient(makerConn),
		hwClient:      types.NewHotWalletClient(hwConn),
		watcherClient: types.NewWatcherClient(watcherConn),
		orderDuration: cfg.OrderDuration,
		db:            db,
		Logger:        logger,
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

	if err := d.db.CreateQuote(quote); err != nil {
		return nil, err
	}

	fmt.Println(fmt.Sprintf("created order for quote with ID (%s)", res.QuoteId))
	return quote, nil
}

func (d *Dealer) ValidateOrder(ctx context.Context, req *types.ValidateOrderRequest) error {
	res, err := d.hwClient.ValidateOrder(ctx, req)
	if err != nil {
		return nil
	}

	if !res.Valid {
		return fmt.Errorf("fill failed validation: %s", res.Info)
	}
	return nil
}

func (d *Dealer) ExecuteZeroExTransaction(ctx context.Context, req *types.ExecuteZeroExTransactionRequest) (*types.ExecuteZeroExTransactionResponse, error) {
	return d.hwClient.ExecuteZeroExTransaction(ctx, req)
}

func (d *Dealer) WatchTransaction(ctx context.Context, quoteId string, txHash string) (*types.WatchTransactionResponse, error) {
	req := &types.WatchTransactionRequest{
		QuoteId: quoteId,
		TxHash:  txHash,
	}

	return d.watcherClient.WatchTransaction(ctx, req)
}

func (d *Dealer) GetQuote(quoteId string) (*types.Quote, error) {
	return d.db.GetQuote(quoteId)
}

func (d *Dealer) GetOrder(quoteId string) (*zeroex.SignedOrder, error) {
	quote, err := d.GetQuote(quoteId)
	if err != nil {
		return nil, err
	}

	order, err := quote.Order.ToZeroExSignedOrder()
	if err != nil {
		return nil, err
	}

	return order, nil
}
