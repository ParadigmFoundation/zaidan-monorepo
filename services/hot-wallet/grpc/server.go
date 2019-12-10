package grpc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/0xProject/0x-mesh/ethereum"

	"github.com/0xProject/0x-mesh/zeroex"

	hw "github.com/ParadigmFoundation/zaidan-monorepo/services/hot-wallet"
	"github.com/ParadigmFoundation/zaidan-monorepo/common/eth"
	"google.golang.org/grpc"
)

type Server struct {
	pvr  *eth.Provider
	grpc *grpc.Server

	makerAccountIndex  int
	senderAccountIndex int
}

func NewServer(provider *eth.Provider, makerAccountIndex int, senderAccountIndex int) *Server {
	srv := &Server{
		pvr:                provider,
		grpc:               grpc.NewServer(),
		makerAccountIndex:  makerAccountIndex,
		senderAccountIndex: senderAccountIndex,
	}
	hw.RegisterHotWalletServer(srv.grpc, srv)
	return srv
}

// HashOrder implements hw.HotWalletServer
func (srv *Server) HashOrder(ctx context.Context, req *hw.HashOrderRequest) (*hw.HashOrderResponse, error) {
	order, err := req.GetOrder().ToZeroExOrder()
	if err != nil {
		return nil, err
	}

	hash, err := order.ComputeOrderHash()
	if err != nil {
		return nil, err
	}

	return &hw.HashOrderResponse{Hash: hash.Bytes()}, nil
}

// SignOrder implements hw.HotWalletServer
func (srv *Server) SignOrder(ctx context.Context, req *hw.SignOrderRequest) (*hw.SignOrderResponse, error) {
	order, err := req.GetOrder().ToZeroExOrder()
	if err != nil {
		return nil, err
	}

	signedOrder, err := zeroex.SignOrder(srv.pvr, order)
	if err != nil {
		return nil, err
	}

	protoOrder := hw.SignedOrderToProto(signedOrder)
	return &hw.SignOrderResponse{SignedOrder: protoOrder}, nil
}

// CreateOrder implements hw.HotWalletServer
func (srv *Server) CreateOrder(ctx context.Context, req *hw.CreateOrderRequest) (*hw.CreateOrderResponse, error) {
	chainId, err := srv.pvr.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	salt, err := eth.GeneratePseudoRandomSalt()
	if err != nil {
		return nil, err
	}

	addrs, err := ethereum.GetContractAddressesForChainID(int(chainId.Int64()))
	if err != nil {
		return nil, err
	}

	makerAssetAmount, ok := new(big.Int).SetString(req.MakerAssetAmount, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse 'makerAssetAmount'")
	}

	takerAssetAmount, ok := new(big.Int).SetString(req.TakerAssetAmount, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse 'takerAssetAmount'")
	}

	expirationTimeSeconds, ok := new(big.Int).SetString(req.ExpirationTimeSeconds, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse 'expirationTimeSeconds'")
	}

	order := &zeroex.Order{
		ChainID:               chainId,
		ExchangeAddress:       addrs.Exchange,
		MakerAddress:          srv.pvr.Accounts()[srv.makerAccountIndex].Address,
		MakerAssetData:        eth.EncodeERC20AssetData(common.HexToAddress(req.MakerAssetAddress)),
		MakerFeeAssetData:     common.LeftPadBytes([]byte{0x0}, eth.AssetDataLength),
		MakerAssetAmount:      makerAssetAmount,
		MakerFee:              big.NewInt(0),
		TakerAddress:          common.HexToAddress(req.TakerAddress),
		TakerAssetData:        eth.EncodeERC20AssetData(common.HexToAddress(req.TakerAssetAddress)),
		TakerFeeAssetData:     common.LeftPadBytes([]byte{0x0}, eth.AssetDataLength),
		TakerAssetAmount:      takerAssetAmount,
		TakerFee:              big.NewInt(0),
		SenderAddress:         srv.pvr.Accounts()[srv.senderAccountIndex].Address,
		FeeRecipientAddress:   common.Address{},
		ExpirationTimeSeconds: expirationTimeSeconds,
		Salt:                  salt,
	}

	signedOrder, err := zeroex.SignOrder(srv.pvr, order)
	if err != nil {
		return nil, err
	}

	return &hw.CreateOrderResponse{Order: hw.SignedOrderToProto(signedOrder)}, nil
}
