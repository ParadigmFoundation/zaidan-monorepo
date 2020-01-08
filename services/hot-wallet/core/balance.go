package core

import (
	"context"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/zrx"
)

// GetTokenBalance implements grpc.HotWalletServer
func (hw *HotWallet) GetBalance(ctx context.Context, req *grpc.GetBalanceRequest) (*grpc.GetBalanceResponse, error) {
	owner := common.HexToAddress(req.OwnerAddress)
	token := common.HexToAddress(req.TokenAddress)

	// special case for checking Ether balance
	if token == zrx.NULL_ADDRESS {
		return hw.getEtherBalance(ctx, owner)
	} else {
		return hw.getTokenBalance(ctx, owner, token)
	}
}

func (hw *HotWallet) getEtherBalance(ctx context.Context, owner common.Address) (*grpc.GetBalanceResponse, error) {
	balance, err := hw.provider.Client().BalanceAt(ctx, owner, nil)
	if err != nil {
		return nil, err
	}

	return &grpc.GetBalanceResponse{
		OwnerAddress: strings.ToLower(owner.Hex()),
		TokenAddress: strings.ToLower((common.Address{}).Hex()),
		Balance:      balance.String(),
	}, nil
}

func (hw *HotWallet) getTokenBalance(ctx context.Context, owner common.Address, token common.Address) (*grpc.GetBalanceResponse, error) {
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
		OwnerAddress: strings.ToLower(owner.Hex()),
		TokenAddress: strings.ToLower(token.Hex()),
		Balance:      balance.String(),
	}, nil
}
