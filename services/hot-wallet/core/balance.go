package core

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

// GetEtherBalance implements grpc.HotWalletServer
func (hw *HotWallet) GetEtherBalance(ctx context.Context, req *grpc.GetBalanceRequest) (*grpc.GetBalanceResponse, error) {
	owner := common.HexToAddress(req.OwnerAddress)
	balance, err := hw.provider.Client().BalanceAt(ctx, owner, nil)
	if err != nil {
		return nil, err
	}

	return &grpc.GetBalanceResponse{
		OwnerAddress: grpc.NormalizeAddress(owner),
		Balance:      balance.String(),
	}, nil
}

// GetTokenBalance implements grpc.HotWalletServer
func (hw *HotWallet) GetTokenBalance(ctx context.Context, req *grpc.GetBalanceRequest) (*grpc.GetBalanceResponse, error) {
	owner := common.HexToAddress(req.OwnerAddress)
	token := common.HexToAddress(req.TokenAddress)

	devUtils := hw.zrxHelper.DevUtils()
	assetData, err := devUtils.EncodeERC20AssetData(nil, token)
	if err != nil {
		return nil, err
	}

	balance, err := devUtils.GetBalance(nil, owner, assetData)
	if err != nil {
		return nil, err
	}

	return &grpc.GetBalanceResponse{
		OwnerAddress: grpc.NormalizeAddress(owner),
		TokenAddress: grpc.NormalizeAddress(token),
		Balance:      balance.String(),
	}, nil
}
