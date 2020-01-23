package core

import (
	"context"
	"math/big"
	"testing"

	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/0xProject/0x-mesh/zeroex"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/zrx"
)

func testExecuteZrxTransaction(hw *HotWallet, t *testing.T) {
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

	zrxAssetData := zrx.EncodeERC20AssetData(hw.zrxHelper.ContractAddresses.ZRXToken)
	wethAssetData := zrx.EncodeERC20AssetData(hw.zrxHelper.ContractAddresses.WETH9)

	// fetch balance and allowance for maker/taker
	makerZrxBalancePreTrade, err := hw.zrxHelper.DevUtils().GetBalance(nil, hw.makerAddress, zrxAssetData)
	require.NoError(t, err)
	takerWethBalancePreTrade, err := hw.zrxHelper.DevUtils().GetBalance(nil, takerAcct.Address, wethAssetData)
	require.NoError(t, err)
	makerWethBalancePreTrade, err := hw.zrxHelper.DevUtils().GetBalance(nil, hw.makerAddress, wethAssetData)
	require.NoError(t, err)
	takerZrxBalancePreTrade, err := hw.zrxHelper.DevUtils().GetBalance(nil, takerAcct.Address, zrxAssetData)
	makerZrxAllowance, err := hw.zrxHelper.DevUtils().GetAssetProxyAllowance(nil, hw.makerAddress, zrxAssetData)
	require.NoError(t, err)
	takerWethAllowance, err := hw.zrxHelper.DevUtils().GetAssetProxyAllowance(nil, takerAcct.Address, wethAssetData)
	require.NoError(t, err)

	// ensure maker has allowance and balance > trade amount
	assert.Equal(t, 1, makerZrxBalancePreTrade.Cmp(makerAmount))
	assert.Equal(t, 1, makerZrxAllowance.Cmp(makerAmount))

	// ensure taker has allowance and balance > trade amount
	assert.Equal(t, 1, takerWethBalancePreTrade.Cmp(takerAmount))
	assert.Equal(t, 1, takerWethAllowance.Cmp(takerAmount))

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
	require.NoError(t, hw.zrxHelper.ValidateFill(context.Background(), signedOrder, order.TakerAssetAmount))

	// prep and sign ZEIP-18 fill transaction
	txData, err := hw.zrxHelper.GetFillOrderCallData(signedOrder.Order, order.TakerAssetAmount, signedOrder.Signature)
	require.NoError(t, err)
	fillTx := &zrx.Transaction{
		Salt:                  signedOrder.Salt,
		ExpirationTimeSeconds: signedOrder.ExpirationTimeSeconds,
		GasPrice:              big.NewInt(1),
		SignerAddress:         takerAcct.Address,
		Data:                  txData,
	}
	signedTransaction, err := zrx.SignTransaction(hw.provider, fillTx, int(order.ChainID.Int64()))
	require.NoError(t, err)

	// create and send gRPC request to execute
	executeReq := &grpc.ExecuteZeroExTransactionRequest{
		Transaction: grpc.ZeroExTransactionToProto(&signedTransaction.Transaction),
		Signature:   signedTransaction.Signature,
	}
	_, err = hw.ExecuteZeroExTransaction(context.Background(), executeReq)
	require.NoError(t, err)

	// check that balances have updated accordingly
	makerZrxBalancePostTrade, err := hw.zrxHelper.DevUtils().GetBalance(nil, hw.makerAddress, zrxAssetData)
	require.NoError(t, err)
	takerWethBalancePostTrade, err := hw.zrxHelper.DevUtils().GetBalance(nil, takerAcct.Address, wethAssetData)
	require.NoError(t, err)
	makerWethBalancePostTrade, err := hw.zrxHelper.DevUtils().GetBalance(nil, hw.makerAddress, wethAssetData)
	require.NoError(t, err)
	takerZrxBalancePostTrade, err := hw.zrxHelper.DevUtils().GetBalance(nil, takerAcct.Address, zrxAssetData)
	require.NoError(t, err)
	assert.Equal(t, 0, new(big.Int).Sub(makerZrxBalancePreTrade, makerAmount).Cmp(makerZrxBalancePostTrade))
	assert.Equal(t, 0, new(big.Int).Add(makerWethBalancePreTrade, takerAmount).Cmp(makerWethBalancePostTrade))
	assert.Equal(t, 0, new(big.Int).Sub(takerWethBalancePreTrade, takerAmount).Cmp(takerWethBalancePostTrade))
	assert.Equal(t, 0, new(big.Int).Add(takerZrxBalancePreTrade, makerAmount).Cmp(takerZrxBalancePostTrade))
}
