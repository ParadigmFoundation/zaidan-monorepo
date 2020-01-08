package core

import (
	"context"
	"fmt"
	"math/big"

	"github.com/0xProject/0x-mesh/zeroex"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/zrx"
)

// CreateOrder implements grpc.HotWalletServer
func (hw *HotWallet) CreateOrder(ctx context.Context, req *grpc.CreateOrderRequest) (*grpc.CreateOrderResponse, error) {
	signedOrder, err := hw.createAndSignOrder(req)
	if err != nil {
		return nil, err
	}

	orderHash, err := signedOrder.ComputeOrderHash()
	if err != nil {
		return nil, err
	}

	return &grpc.CreateOrderResponse{Order: grpc.SignedOrderToProto(signedOrder), Hash: orderHash.Bytes()}, nil
}

func (hw *HotWallet) createAndSignOrder(cfg *grpc.CreateOrderRequest) (*zeroex.SignedOrder, error) {
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

	order, err := hw.zrxHelper.CreateOrder(
		hw.makerAddress,
		takerAddress,
		hw.senderAddress,
		zrx.NULL_ADDRESS,
		makerAssetAddress,
		takerAssetAddress,
		makerAssetAmount,
		takerAssetAmount,
		big.NewInt(0),
		big.NewInt(0),
		zrx.NULL_ADDRESS,
		zrx.NULL_ADDRESS,
		expirationTimeSeconds,
	)
	if err != nil {
		return nil, err
	}

	return zeroex.SignOrder(hw.provider, order)
}
