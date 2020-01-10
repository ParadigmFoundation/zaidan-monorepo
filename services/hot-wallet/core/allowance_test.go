package core

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

const MAX_UINT_256 = "115792089237316195423570985008687907853269984665640564039457584007913129639935"

func testGetSetAllowance(hw *HotWallet, t *testing.T) {
	assetData, _ := hw.zrxHelper.DevUtils().EncodeERC20AssetData(nil, hw.zrxHelper.ContractAddresses.ZRXToken)
	allowanceBefore, err := hw.zrxHelper.DevUtils().GetAssetProxyAllowance(nil, hw.makerAddress, assetData)
	require.NoError(t, err)
	assert.Equal(t, "0", allowanceBefore.String())

	setReq := &grpc.SetAllowanceRequest{
		SpenderAddress: hw.zrxHelper.ContractAddresses.ERC20Proxy.Hex(),
		TokenAddress:   hw.zrxHelper.ContractAddresses.ZRXToken.Hex(),
	}

	setRes, err := hw.SetAllowance(context.Background(), setReq)
	require.NoError(t, err)

	assert.Equal(t, hw.zrxHelper.ContractAddresses.ZRXToken, common.HexToAddress(setRes.TokenAddress))
	assert.Equal(t, hw.zrxHelper.ContractAddresses.ERC20Proxy, common.HexToAddress(setRes.SpenderAddress))
	assert.Equal(t, hw.makerAddress, common.HexToAddress(setRes.OwnerAddress))

	allowanceAfter, err := hw.zrxHelper.DevUtils().GetAssetProxyAllowance(nil, hw.makerAddress, assetData)
	require.NoError(t, err)
	assert.Equal(t, MAX_UINT_256, allowanceAfter.String())

	getReq := &grpc.GetAllowanceRequest{
		SpenderAddress: hw.zrxHelper.ContractAddresses.ERC20Proxy.Hex(),
		OwnerAddress:   hw.makerAddress.Hex(),
		TokenAddress:   hw.zrxHelper.ContractAddresses.ZRXToken.Hex(),
	}

	getRes, err := hw.GetAllowance(context.Background(), getReq)
	require.NoError(t, err)

	assert.Equal(t, hw.zrxHelper.ContractAddresses.ZRXToken, common.HexToAddress(getRes.TokenAddress))
	assert.Equal(t, hw.zrxHelper.ContractAddresses.ERC20Proxy, common.HexToAddress(getRes.ProxyAddress))
	assert.Equal(t, hw.makerAddress, common.HexToAddress(getRes.OwnerAddress))
	assert.Equal(t, MAX_UINT_256, getRes.Allowance)
}
