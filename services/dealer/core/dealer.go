package core

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/grpc"

	"github.com/ethereum/go-ethereum/log"

	types "github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/services/dealer/store"
)

// DealerConfig defines configuration for a Dealer
type DealerConfig struct {
	MakerBindAddress     string
	HotWalletBindAddress string
}

// Dealer is the core dealer service that interacts with other services
type Dealer struct {
	makerClient types.MakerClient
	hwClient    types.HotWalletClient

	db         store.Store
	cancelFunc context.CancelFunc
	logger     log.Logger
}

// NewDealer creates a new Dealer given ctx context and cfg configuration
func NewDealer(ctx context.Context, cfg DealerConfig) (*Dealer, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	logger := log.New(ctx)

	makerConn, err := grpc.DialContext(ctx, cfg.MakerBindAddress)
	if err != nil {
		return nil, err
	}

	hwConn, err := grpc.DialContext(ctx, cfg.HotWalletBindAddress)
	if err != nil {
		return nil, err
	}

	return &Dealer{
		makerClient: types.NewMakerClient(makerConn),
		hwClient:    types.NewHotWalletClient(hwConn),
		cancelFunc:  cancelFunc,
		logger:      logger,
	}, nil
}

// @todo: hrharder - should eventually return types.Quote
func (d *Dealer) FetchQuote(ctx context.Context, req *types.GetQuoteRequest) error {
	res, err := d.makerClient.GetQuote(ctx, req)
	if err != nil {
		return err
	}

	// todo: remove
	str, err := json.Marshal(res)
	if err != nil {
		return err
	}

	fmt.Println(string(str))

	makerAssetAddress, err := d.getAssetAddress(req.MakerAsset)
	if err != nil {
		return err
	}

	takerAssetAddress, err := d.getAssetAddress(req.TakerAsset)
	if err != nil {
		return err
	}

	orderReq := &types.CreateOrderRequest{
		TakerAddress:          req.TakerAddress,
		MakerAssetAddress:     makerAssetAddress,
		TakerAssetAddress:     takerAssetAddress,
		MakerAssetAmount:      res.MakerSize,
		TakerAssetAmount:      res.TakerSize,
		ExpirationTimeSeconds: res.Expiration,
	}

	order, err := d.hwClient.CreateOrder(ctx, orderReq)
	if err != nil {
		return err
	}

	obts, err := json.Marshal(order)
	if err != nil {
		return err
	}

	fmt.Println(string(obts))
	return nil
}

func (d *Dealer) getAssetAddress(ticker string) (string, error) {
	return "", nil
}
