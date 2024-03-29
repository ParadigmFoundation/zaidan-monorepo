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

	txData, err := hw.zrxHelper.GetFillOrderCallData(signedOrder.Order, signedOrder.TakerAssetAmount, signedOrder.Signature)
	if err != nil {
		return nil, err
	}

	salt, err := zrx.GeneratePseudoRandomSalt()
	if err != nil {
		return nil, err
	}

	gasPrice, err := hw.provider.Client().SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	zrxTx := &zrx.Transaction{
		Salt:                  salt,
		ExpirationTimeSeconds: new(big.Int).Set(signedOrder.ExpirationTimeSeconds),
		GasPrice:              gasPrice,
		SignerAddress:         signedOrder.TakerAddress,
		Data:                  txData,
	}

	zrxTxHash, err := zrxTx.ComputeHashForChainID(hw.chainId)
	if err != nil {
		return nil, err
	}

	return &grpc.CreateOrderResponse{
		Order:                 grpc.SignedOrderToProto(signedOrder),
		OrderHash:             orderHash.Hex(),
		ZeroExTransaction:     grpc.ZeroExTransactionToProto(zrxTx),
		ZeroExTransactionHash: zrxTxHash.Hex(),
	}, nil
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
		big.NewInt(cfg.ExpirationTimeSeconds),
	)
	if err != nil {
		return nil, err
	}

	return zeroex.SignOrder(hw.provider, order)
}
