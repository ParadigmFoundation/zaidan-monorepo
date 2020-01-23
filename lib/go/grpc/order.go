package grpc

// This file adds extra functionality to the generated types.pb.go

import (
	fmt "fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/0xProject/0x-mesh/zeroex"
)

// NormalizeAddress converts an Ethereum address to an all-lowercase string representation with a leading "0x"
func NormalizeAddress(address common.Address) string {
	return strings.ToLower(address.Hex())
}

// OrderToProto converts a 0x unsigned order to a wire-safe protobuf representation
func OrderToProto(order *zeroex.Order) *Order {
	return &Order{
		ChainId:               order.ChainID.Uint64(),
		ExchangeAddress:       NormalizeAddress(order.ExchangeAddress),
		MakerAddress:          NormalizeAddress(order.MakerAddress),
		MakerAssetData:        hexutil.Encode(order.MakerAssetData),
		MakerAssetAmount:      order.MakerAssetAmount.String(),
		MakerFee:              order.MakerFee.String(),
		TakerAddress:          NormalizeAddress(order.TakerAddress),
		TakerAssetData:        hexutil.Encode(order.TakerAssetData),
		TakerAssetAmount:      order.TakerAssetAmount.String(),
		TakerFee:              order.TakerFee.String(),
		SenderAddress:         NormalizeAddress(order.SenderAddress),
		FeeRecipientAddress:   NormalizeAddress(order.FeeRecipientAddress),
		ExpirationTimeSeconds: order.ExpirationTimeSeconds.String(),
		Salt:                  order.Salt.String(),
	}
}

// SignedOrderToProto converts a 0x signed order to a wire-safe protobuf representation
func SignedOrderToProto(order *zeroex.SignedOrder) *SignedOrder {
	return &SignedOrder{
		ChainId:               order.ChainID.Uint64(),
		ExchangeAddress:       NormalizeAddress(order.ExchangeAddress),
		MakerAddress:          NormalizeAddress(order.MakerAddress),
		MakerAssetData:        hexutil.Encode(order.MakerAssetData),
		MakerFeeAssetData:     hexutil.Encode(order.MakerFeeAssetData),
		MakerAssetAmount:      order.MakerAssetAmount.String(),
		MakerFee:              order.MakerFee.String(),
		TakerAddress:          NormalizeAddress(order.TakerAddress),
		TakerAssetData:        hexutil.Encode(order.TakerAssetData),
		TakerFeeAssetData:     hexutil.Encode(order.TakerFeeAssetData),
		TakerAssetAmount:      order.TakerAssetAmount.String(),
		TakerFee:              order.TakerFee.String(),
		SenderAddress:         NormalizeAddress(order.SenderAddress),
		FeeRecipientAddress:   NormalizeAddress(order.FeeRecipientAddress),
		ExpirationTimeSeconds: order.ExpirationTimeSeconds.String(),
		Salt:                  order.Salt.String(),
		Signature:             hexutil.Encode(order.Signature),
	}
}

// ToZeroExOrder converts a protobuf unsigned order message to a 0x SignedOrder
func (o *Order) ToZeroExOrder() (*zeroex.Order, error) {
	chainId, ok := new(big.Int).SetString(strconv.FormatUint(o.ChainId, 10), 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse 'chainId'")
	}

	makerAssetAmount, ok := new(big.Int).SetString(o.MakerAssetAmount, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse 'makerAssetAmount'")
	}

	makerFee, ok := new(big.Int).SetString(o.MakerFee, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse 'makerFee'")
	}

	takerAssetAmount, ok := new(big.Int).SetString(o.TakerAssetAmount, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse 'takerAssetAmount'")
	}

	takerFee, ok := new(big.Int).SetString(o.TakerFee, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse 'takerFee'")
	}

	salt, ok := new(big.Int).SetString(o.Salt, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse 'salt'")
	}

	expirationTimeSeconds, ok := new(big.Int).SetString(o.ExpirationTimeSeconds, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse 'expirationTimeSeconds'")
	}

	return &zeroex.Order{
		ChainID:               chainId,
		ExchangeAddress:       common.HexToAddress(o.ExchangeAddress),
		MakerAddress:          common.HexToAddress(o.MakerAddress),
		MakerAssetData:        common.FromHex(o.MakerAssetData),
		MakerFeeAssetData:     common.FromHex(o.MakerFeeAssetData),
		MakerAssetAmount:      makerAssetAmount,
		MakerFee:              makerFee,
		TakerAddress:          common.HexToAddress(o.TakerAddress),
		TakerAssetData:        common.FromHex(o.TakerAssetData),
		TakerFeeAssetData:     common.FromHex(o.TakerFeeAssetData),
		TakerAssetAmount:      takerAssetAmount,
		TakerFee:              takerFee,
		SenderAddress:         common.HexToAddress(o.SenderAddress),
		FeeRecipientAddress:   common.HexToAddress(o.FeeRecipientAddress),
		ExpirationTimeSeconds: expirationTimeSeconds,
		Salt:                  salt,
	}, nil
}

// ToZeroExSignedOrder converts a protobuf unsigned order message to a 0x SignedOrder
func (o *SignedOrder) ToZeroExSignedOrder() (*zeroex.SignedOrder, error) {
	chainId, ok := new(big.Int).SetString(strconv.FormatUint(o.ChainId, 10), 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse 'chainId'")
	}

	makerAssetAmount, ok := new(big.Int).SetString(o.MakerAssetAmount, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse 'makerAssetAmount'")
	}

	makerFee, ok := new(big.Int).SetString(o.MakerFee, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse 'makerFee'")
	}

	takerAssetAmount, ok := new(big.Int).SetString(o.TakerAssetAmount, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse 'takerAssetAmount'")
	}

	takerFee, ok := new(big.Int).SetString(o.TakerFee, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse 'takerFee'")
	}

	salt, ok := new(big.Int).SetString(o.Salt, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse 'salt'")
	}

	expirationTimeSeconds, ok := new(big.Int).SetString(o.ExpirationTimeSeconds, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse 'expirationTimeSeconds'")
	}

	order := zeroex.Order{
		ChainID:               chainId,
		ExchangeAddress:       common.HexToAddress(o.ExchangeAddress),
		MakerAddress:          common.HexToAddress(o.MakerAddress),
		MakerAssetData:        common.FromHex(o.MakerAssetData),
		MakerFeeAssetData:     common.FromHex(o.MakerFeeAssetData),
		MakerAssetAmount:      makerAssetAmount,
		MakerFee:              makerFee,
		TakerAddress:          common.HexToAddress(o.TakerAddress),
		TakerAssetData:        common.FromHex(o.TakerAssetData),
		TakerFeeAssetData:     common.FromHex(o.TakerFeeAssetData),
		TakerAssetAmount:      takerAssetAmount,
		TakerFee:              takerFee,
		SenderAddress:         common.HexToAddress(o.SenderAddress),
		FeeRecipientAddress:   common.HexToAddress(o.FeeRecipientAddress),
		ExpirationTimeSeconds: expirationTimeSeconds,
		Salt:                  salt,
	}

	return &zeroex.SignedOrder{
		Order:     order,
		Signature: common.FromHex(o.Signature),
	}, nil
}
