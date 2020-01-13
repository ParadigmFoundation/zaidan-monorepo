package core

import (
	"context"
	"math/big"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/zrx"

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
func (hw *HotWallet) SetAllowance(ctx context.Context, req *grpc.SetAllowanceRequest) (*grpc.SetAllowanceResponse, error) {
	token := common.HexToAddress(req.TokenAddress)
	spender := common.HexToAddress(req.SpenderAddress)

	if spender == zrx.NULL_ADDRESS {
		spender = hw.zrxHelper.ContractAddresses.ERC20Proxy
	}

	allowance := new(big.Int)
	if req.Allowance == "" {
		u256 := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
		_ = allowance.Sub(u256, big.NewInt(1))
	} else {
		_, ok := allowance.SetString(req.Allowance, 10)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "unable to parse value for 'allowance'")
		}
	}

	tx, err := hw.tokenManager.Approve(token, spender, allowance)
	if err != nil {
		return nil, err
	}

	return &grpc.SetAllowanceResponse{
		OwnerAddress:    grpc.NormalizeAddress(hw.makerAddress),
		SpenderAddress:  grpc.NormalizeAddress(spender),
		TokenAddress:    grpc.NormalizeAddress(token),
		Allowance:       allowance.String(),
		TransactionHash: tx.Hash().Hex(),
	}, nil
}
