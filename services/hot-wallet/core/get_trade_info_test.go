package core

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/zrx"

	"github.com/stretchr/testify/require"
)

func testGetTradeInfo(hw *HotWallet, t *testing.T) {
	expectedGasPrice, err := hw.provider.Client().SuggestGasPrice(context.Background())
	require.NoError(t, err)

	expectedGasLimit := strconv.FormatUint(zrx.EXECUTE_FILL_TX_GAS_LIMIT, 10)
	expectedChainId, err := hw.provider.Client().ChainID(context.Background())
	require.NoError(t, err)

	actualTradeInfo, err := hw.GetTradeInfo(context.Background(), &empty.Empty{})
	assert.NoError(t, err)

	assert.Equal(t, expectedGasPrice.String(), actualTradeInfo.GasPrice)
	assert.Equal(t, expectedGasLimit, actualTradeInfo.GasLimit)
	assert.Equal(t, uint32(expectedChainId.Uint64()), actualTradeInfo.ChainId)
}
