package core

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"

	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func testValidateOrder(hw *HotWallet, t *testing.T) {
	_testValidateValidOrder(hw, t)
	_testValidateInvalidOrder(hw, t)
}

func _testValidateValidOrder(hw *HotWallet, t *testing.T) {
	// weth/zrx token wrappers
	wethToken, err := wrappers.NewWETH9(hw.zrxHelper.ContractAddresses.WETH9, hw.provider.Client())
	require.NoError(t, err)
	zrxToken, err := wrappers.NewZRXToken(hw.zrxHelper.ContractAddresses.ZRXToken, hw.provider.Client())
	require.NoError(t, err)

	// trade amount in base units
	makerAmount := big.NewInt(424)
	takerAmount := big.NewInt(242)

	// maker transaction signer
	makerAcct, err := hw.provider.GetAccount(hw.makerAddress)
	require.NoError(t, err)
	makerKey, err := hw.provider.Wallet().PrivateKey(makerAcct)
	require.NoError(t, err)
	makerTransactor := bind.NewKeyedTransactor(makerKey)

	// taker transaction signer
	takerAcct := hw.provider.Accounts()[2]
	takerKey, err := hw.provider.Wallet().PrivateKey(takerAcct)
	require.NoError(t, err)
	takerTransactor := bind.NewKeyedTransactor(takerKey)

	// wrap some ETH from the taker
	opts := &bind.TransactOpts{
		From:   takerAcct.Address,
		Signer: takerTransactor.Signer,
		Value:  big.NewInt(100000),
	}
	_, err = wethToken.Deposit(opts)
	require.NoError(t, err)

	// set erc-20 approvals for maker/taker accounts for both ZRX and wETH
	_, err = wethToken.Approve(takerTransactor, hw.zrxHelper.ContractAddresses.ERC20Proxy, big.NewInt(1000000))
	require.NoError(t, err)
	_, err = zrxToken.Approve(takerTransactor, hw.zrxHelper.ContractAddresses.ERC20Proxy, big.NewInt(1000000))
	require.NoError(t, err)
	_, err = wethToken.Approve(makerTransactor, hw.zrxHelper.ContractAddresses.ERC20Proxy, big.NewInt(1000000))
	require.NoError(t, err)
	_, err = zrxToken.Approve(makerTransactor, hw.zrxHelper.ContractAddresses.ERC20Proxy, big.NewInt(1000000))
	require.NoError(t, err)

	// create, sign and validate underlying order
	order, err := hw.zrxHelper.CreateOrder(
		hw.makerAddress,
		takerAcct.Address,
		hw.senderAddress,
		common.Address{},
		hw.zrxHelper.ContractAddresses.ZRXToken,
		hw.zrxHelper.ContractAddresses.WETH9,
		makerAmount,
		takerAmount,
		big.NewInt(0),
		big.NewInt(0),
		common.Address{},
		common.Address{},
		big.NewInt(16000000000),
	)
	require.NoError(t, err)

	signedOrder, err := zeroex.SignOrder(hw.provider, order)
	require.NoError(t, err)

	protoOrder := grpc.SignedOrderToProto(signedOrder)
	testReq := &grpc.ValidateOrderRequest{
		Order:            protoOrder,
		TakerAssetAmount: order.TakerAssetAmount.String(),
	}

	testRes, err := hw.ValidateOrder(context.Background(), testReq)
	require.NoError(t, err)
	assert.True(t, testRes.Valid)
}

func _testValidateInvalidOrder(hw *HotWallet, t *testing.T) {
	takerAcct := hw.provider.Accounts()[2]

	// create, sign and validate underlying order
	order, err := hw.zrxHelper.CreateOrder(
		hw.makerAddress,
		takerAcct.Address,
		hw.senderAddress,
		common.Address{},
		hw.zrxHelper.ContractAddresses.ZRXToken,
		hw.zrxHelper.ContractAddresses.WETH9,
		big.NewInt(1000001), // this value is > allowance set in previous test
		big.NewInt(1000002), // this value is > allowance set in previous test
		big.NewInt(0),
		big.NewInt(0),
		common.Address{},
		common.Address{},
		big.NewInt(16000000000),
	)
	require.NoError(t, err)

	signedOrder, err := zeroex.SignOrder(hw.provider, order)
	require.NoError(t, err)

	protoOrder := grpc.SignedOrderToProto(signedOrder)
	testReq := &grpc.ValidateOrderRequest{
		Order:            protoOrder,
		TakerAssetAmount: order.TakerAssetAmount.String(),
	}

	testRes, err := hw.ValidateOrder(context.Background(), testReq)
	require.NoError(t, err)
	assert.False(t, testRes.Valid)

	knownFailureReason := "maker has insufficient balance or allowance for this order to be filled"
	assert.Equal(t, knownFailureReason, testRes.Info)
}
