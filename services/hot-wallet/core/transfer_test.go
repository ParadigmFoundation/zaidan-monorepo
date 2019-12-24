package core

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testTransferEther(hw *HotWallet, t *testing.T) {
	from := hw.makerAddress
	to := hw.senderAddress

	fromBeforeBalance, err := hw.provider.Client().BalanceAt(context.Background(), from, nil)
	require.NoError(t, err)

	toBeforeBalance, err := hw.provider.Client().BalanceAt(context.Background(), to, nil)
	require.NoError(t, err)

	testReq := &grpc.TransferRequest{
		ToAddress: to.Hex(),
		Amount:    "42",
	}

	res, err := hw.Transfer(context.Background(), testReq)
	assert.NoError(t, err)

	bts := common.HexToHash(res.TransactionHash)
	assert.Equal(t, 32, len(bts))

	fromAfterBalance, err := hw.provider.Client().BalanceAt(context.Background(), from, nil)
	require.NoError(t, err)

	toAfterBalance, err := hw.provider.Client().BalanceAt(context.Background(), to, nil)
	require.NoError(t, err)

	// to account should have balance equal to before + transfer amount
	assert.Equal(t, 0, new(big.Int).Sub(toAfterBalance, big.NewInt(42)).Cmp(toBeforeBalance))

	// from account should have balance less than before - transfer amount (because of gas)
	assert.Equal(t, -1, new(big.Int).Add(fromAfterBalance, big.NewInt(42)).Cmp(fromBeforeBalance))
	t.Log(fromBeforeBalance)
	t.Log(fromAfterBalance)
}
