package core

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"

	"github.com/0xProject/0x-mesh/zeroex"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/logger"
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
	log *logger.Logger

	makerClient   types.MakerClient
	hwClient      types.HotWalletClient
	watcherClient types.WatcherClient

	orderDuration int64

	db store.Store
}

// NewDealer creates a new Dealer given ctx context and cfg configuration
func NewDealer(ctx context.Context, db store.Store, cfg DealerConfig) (*Dealer, error) {
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
		log:           logger.New("core"),
	}, nil
}

func (d *Dealer) FetchQuote(ctx context.Context, req *types.GetQuoteRequest) (*types.Quote, error) {
	now := time.Now()
	res, err := d.makerClient.GetQuote(ctx, req)
	if err != nil {
		d.log.WithError(err).Error("failed fetching quote from maker")
		return nil, err
	}

	orderReq := &types.CreateOrderRequest{
		TakerAddress:          req.TakerAddress,
		MakerAssetAddress:     res.MakerAsset,
		TakerAssetAddress:     res.TakerAsset,
		MakerAssetAmount:      res.MakerSize,
		TakerAssetAmount:      res.TakerSize,
		ExpirationTimeSeconds: now.Unix() + d.orderDuration,
	}

	orderRes, err := d.hwClient.CreateOrder(ctx, orderReq)
	if err != nil {
		d.log.WithError(err).Error("failed creating signed order")
		return nil, err
	}

	txInfo := &types.ZeroExTransactionInfo{
		Transaction: orderRes.ZeroExTransaction,
		Order:       orderRes.Order,
	}

	quote := &types.Quote{
		QuoteId:               res.QuoteId,
		MakerAssetAddress:     res.MakerAsset,
		TakerAssetAddress:     res.TakerAsset,
		MakerAssetSize:        res.MakerSize,
		TakerAssetSize:        res.TakerSize,
		Expiration:            res.Expiration,
		ServerTime:            now.UnixNano() / 1e6, // conversion from nanoseconds to milliseconds = ns / 1e6
		ZeroExTransactionHash: orderRes.ZeroExTransactionHash,
		ZeroExTransactionInfo: txInfo,
	}

	if err := d.db.CreateQuote(quote); err != nil {
		d.log.WithError(err).Error("failed creating quote in DB")
		return nil, err
	}

	d.log.WithFields(logrus.Fields{"id": res.QuoteId, "taker": req.TakerAddress}).Info("created new quote")
	return quote, nil
}

func (d *Dealer) ValidateOrder(ctx context.Context, req *types.ValidateOrderRequest) error {
	res, err := d.hwClient.ValidateOrder(ctx, req)
	if err != nil {
		return nil
	}

	if !res.Valid {
		d.log.WithField("reason", res.Info).Error("fill failed validation")
		return nil
	}
	return nil
}

func (d *Dealer) ExecuteZeroExTransaction(ctx context.Context, req *types.ExecuteZeroExTransactionRequest) (*types.ExecuteZeroExTransactionResponse, error) {
	d.log.WithField("taker", req.Transaction.SignerAddress).Info("executing 0x transaction")
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

	order, err := quote.ZeroExTransactionInfo.Order.ToZeroExSignedOrder()
	if err != nil {
		return nil, err
	}

	return order, nil
}
