package core

import (
	"context"
	"errors"
	"time"

	"github.com/golang/protobuf/ptypes/empty"

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
	DealerGrpcBindAddress string

	OrderDuration     int64 // the number of seconds order should be valid for after initial quote is provided
	DealerBindAddress string
}

// Dealer is the core dealer service that interacts with other services
type Dealer struct {
	log *logger.Logger

	makerClient   types.MakerClient
	hwClient      types.HotWalletClient
	watcherClient types.WatcherClient

	orderDuration int64

	db  store.Store
	cfg DealerConfig
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
		cfg:           cfg,
		db:            db,
		log:           logger.New("core"),
	}, nil
}

func (d *Dealer) WithMakerClient(c types.MakerClient) *Dealer {
	d.makerClient = c
	return d
}

func (d *Dealer) WithHWClient(c types.HotWalletClient) *Dealer {
	d.hwClient = c
	return d
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

func (d *Dealer) ValidateOrder(ctx context.Context, req *types.ValidateOrderRequest, quoteId string) error {
	res, err := d.hwClient.ValidateOrder(ctx, req)
	if err != nil {
		d.log.WithError(err).Error("failed to call validate order")
		return err
	}

	if !res.Valid {
		d.log.WithField("reason", res.Info).Error("fill failed validation")
		return nil
	}

	makerRes, err := d.makerClient.CheckQuote(ctx, &types.CheckQuoteRequest{QuoteId: quoteId})
	if err != nil {
		d.log.WithError(err).Error("failed to validate quote with maker")
		return err
	}

	if !makerRes.IsValid {
		d.log.WithField("statusCode", makerRes.Status).Info("validation failed for quote")
		return errors.New("maker rejected fill")
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
		StatusUrls: []string{ d.cfg.DealerGrpcBindAddress, d.cfg.MakerBindAddress },
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

type paginatedMakets []*types.Market

func (p paginatedMakets) Paginate(page, perPage int) []*types.Market {
	if perPage == 0 {
		return p
	}

	offset := page * perPage
	end := offset + perPage
	if end > len(p) {
		end = len(p)
	}

	if offset > len(p) {
		return nil
	}

	return p[offset:end]
}

func (d *Dealer) GetMarkets(mAddr, tAddr string, page, perPage int) ([]*types.Market, error) {
	ctx := context.Background()
	req := &types.GetMarketsRequest{
		MakerAssetAddress: mAddr,
		TakerAssetAddress: tAddr,
	}
	resp, err := d.makerClient.GetMarkets(ctx, req)
	if err != nil {
		return nil, err
	}

	tradeInfo, err := d.hwClient.GetTradeInfo(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	for _, market := range resp.Markets {
		market.TradeInfo = tradeInfo
	}

	mkts := paginatedMakets(resp.Markets).Paginate(page, perPage)
	return mkts, nil
}

func (d *Dealer) CreateTrade(t *types.Trade) error {
	return d.db.CreateTrade(t)
}
