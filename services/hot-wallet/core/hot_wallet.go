package core

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	"github.com/0xProject/0x-mesh/zeroex"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/eth"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/zrx"
)

//
type HotWalletConfig struct {
	// Maximum request length used for the order validator
	OrderValidatorMaxReqLength int

	// Address (0x-prefixed) to use as the maker address for orders
	MakerAddress string

	// Address (0x-prefixed) to use as the sender for order
	SenderAddress string
}

// HotWallet represents a live hot wallet that can interact with the 0x contract system
type HotWallet struct {
	provider  *eth.Provider
	zrxHelper *zrx.ZeroExHelper

	makerAddress  common.Address
	senderAddress common.Address

	logger log.Logger
}

// NewHotWallet creates a new hot wallet with the supplied provider and configuration
func NewHotWallet(provider *eth.Provider, cfg HotWalletConfig) (*HotWallet, error) {
	zrxHelper, err := zrx.NewZeroExHelper(provider.Client(), cfg.OrderValidatorMaxReqLength)
	if err != nil {
		return nil, err
	}

	makerAddress := common.HexToAddress(cfg.MakerAddress)
	senderAddress := common.HexToAddress(cfg.SenderAddress)
	if !provider.CanSignWithAddress(makerAddress) || !provider.CanSignWithAddress(senderAddress) {
		return nil, fmt.Errorf("unable to sign with maker or sender address")
	}

	return &HotWallet{
		provider:      provider,
		zrxHelper:     zrxHelper,
		makerAddress:  makerAddress,
		senderAddress: senderAddress,
		logger:        log.New(context.Background()),
	}, nil
}

// CreateOrder implements grpc.HotWalletServer
func (hw *HotWallet) CreateOrder(ctx context.Context, req *grpc.CreateOrderRequest) (*grpc.CreateOrderResponse, error) {
	signedOrder, err := hw.createAndSignOrder(*req)
	if err != nil {
		return nil, err
	}

	orderHash, err := signedOrder.ComputeOrderHash()
	if err != nil {
		return nil, err
	}

	return &grpc.CreateOrderResponse{Order: grpc.SignedOrderToProto(signedOrder), Hash: orderHash.Bytes()}, nil
}

func (hw *HotWallet) createAndSignOrder(cfg grpc.CreateOrderRequest) (*zeroex.SignedOrder, error) {
	takerAddress := common.HexToAddress(cfg.TakerAddress)
	makerAssetAddress := common.HexToAddress(cfg.MakerAssetAddress)
	takerAssetAddress := common.HexToAddress(cfg.TakerAssetAddress)

	makerAssetAmount, ok := new(big.Int).SetString(cfg.MakerAssetAmount, 10)
	if !ok {
		return nil, fmt.Errorf(`unable to parse "makerAssetAmount"`)
	}
	takerAssetAmount, ok := new(big.Int).SetString(cfg.TakerAssetAmount, 10)
	if !ok {
		return nil, fmt.Errorf(`unable to parse "takerAssetAmount"`)
	}
	expirationTimeSeconds, ok := new(big.Int).SetString(cfg.ExpirationTimeSeconds, 10)
	if !ok {
		return nil, fmt.Errorf(`unable to parse "expirationTimeSeconds"`)
	}

	order, err := hw.zrxHelper.CreateOrder(hw.makerAddress, takerAddress, hw.senderAddress, zrx.NULL_ADDRESS, makerAssetAddress, takerAssetAddress, makerAssetAmount, takerAssetAmount, big.NewInt(0), big.NewInt(0), zrx.NULL_ADDRESS, zrx.NULL_ADDRESS, expirationTimeSeconds)
	if err != nil {
		return nil, err
	}

	return zeroex.SignOrder(hw.provider, order)
}
