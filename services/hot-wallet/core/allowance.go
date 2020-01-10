package core

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

// GetAllowance implements grpc.HotWalletServer
func (hw *HotWallet) GetAllowance(ctx context.Context, req *grpc.GetAllowanceRequest) (*grpc.GetAllowanceResponse, error) {
	owner := common.HexToAddress(req.OwnerAddress)
	token := common.HexToAddress(req.TokenAddress)

	devUtils := hw.zrxHelper.DevUtils()
	assetData, err := devUtils.EncodeERC20AssetData(nil, token)
	if err != nil {
		return nil, err
	}

	allowance, err := devUtils.GetAssetProxyAllowance(nil, owner, assetData)
	if err != nil {
		return nil, err
	}

	return &grpc.GetAllowanceResponse{
		OwnerAddress: grpc.NormalizeAddress(owner),
		TokenAddress: grpc.NormalizeAddress(token),
		ProxyAddress: grpc.NormalizeAddress(hw.zrxHelper.ContractAddresses.ERC20Proxy),
		Allowance:    allowance.String(),
	}, nil
}

// SetAllowance implements grpc.HotWalletServer
// func (hw *HotWallet) SetAllowance(ctx context.Context, req *grpc.SetAllowanceRequest)
