package core

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/0xProject/0x-mesh/ethereum"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/zrx"
)

func testCreateOrder(hw *HotWallet, t *testing.T) {
	zrxAddresses, err := ethereum.GetContractAddressesForChainID(zrx.ZeroExTestChainID)
	require.NoError(t, err)

	testReq := &grpc.CreateOrderRequest{
		MakerAssetAddress:     zrxAddresses.WETH9.Hex(),
		TakerAssetAddress:     zrxAddresses.ZRXToken.Hex(),
		MakerAssetAmount:      "1234",
		TakerAssetAmount:      "5678",
		ExpirationTimeSeconds: 4242,
	}

	orderRes, err := hw.CreateOrder(context.Background(), testReq)
	require.NoError(t, err)

	order, err := orderRes.Order.ToZeroExSignedOrder()
	assert.NoError(t, err)

	expectedMakerAssetAmount := big.NewInt(1234)
	expectedTakerAssetAmount := big.NewInt(5678)
	expectedExpirationTime := big.NewInt(4242)

	assert.Equal(t, expectedMakerAssetAmount, order.MakerAssetAmount)
	assert.Equal(t, expectedTakerAssetAmount, order.TakerAssetAmount)
	assert.Equal(t, expectedExpirationTime, order.ExpirationTimeSeconds)
}
