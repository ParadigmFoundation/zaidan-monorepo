package core

import (
	"context"
	"testing"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/eth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

func testGetSetAllowance(hw *HotWallet, t *testing.T) {
	assetData, _ := hw.zrxHelper.DevUtils().EncodeERC20AssetData(nil, hw.zrxHelper.ContractAddresses.ZRXToken)
	allowanceBefore, err := hw.zrxHelper.DevUtils().GetAssetProxyAllowance(nil, hw.makerAddress, assetData)
	require.NoError(t, err)
	assert.Equal(t, "0", allowanceBefore.String())

	setReq := &grpc.SetAllowanceRequest{
		TokenAddress: hw.zrxHelper.ContractAddresses.ZRXToken.Hex(),
	}

	setRes, err := hw.SetAllowance(context.Background(), setReq)
	require.NoError(t, err)

	assert.Equal(t, hw.zrxHelper.ContractAddresses.ZRXToken, common.HexToAddress(setRes.TokenAddress))
	assert.Equal(t, hw.zrxHelper.ContractAddresses.ERC20Proxy, common.HexToAddress(setRes.ProxyAddress))
	assert.Equal(t, hw.makerAddress, common.HexToAddress(setRes.OwnerAddress))

	allowanceAfter, err := hw.zrxHelper.DevUtils().GetAssetProxyAllowance(nil, hw.makerAddress, assetData)
	require.NoError(t, err)
	assert.Equal(t, eth.MAX_UINT256.String(), allowanceAfter.String())

	getReq := &grpc.GetAllowanceRequest{
		OwnerAddress: hw.makerAddress.Hex(),
		TokenAddress: hw.zrxHelper.ContractAddresses.ZRXToken.Hex(),
	}

	getRes, err := hw.GetAllowance(context.Background(), getReq)
	require.NoError(t, err)

	assert.Equal(t, hw.zrxHelper.ContractAddresses.ZRXToken, common.HexToAddress(getRes.TokenAddress))
	assert.Equal(t, hw.zrxHelper.ContractAddresses.ERC20Proxy, common.HexToAddress(getRes.ProxyAddress))
	assert.Equal(t, hw.makerAddress, common.HexToAddress(getRes.OwnerAddress))
	assert.Equal(t, eth.MAX_UINT256.String(), getRes.Allowance)
}
