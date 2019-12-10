package hw

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"
)

const TestOrderJSON = `{"chainId":1337,"exchangeAddress":"0x48bacb9266a570d521063ef5dd96e61686dbe788","makerAssetData":"0xf47261b00000000000000000000000006b175474e89094c44da98b954eedeac495271d0f","takerAssetData":"0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","takerAddress":"0x0000000000000000000000000000000000000000","makerAssetAmount":"42000000000000000000","takerAssetAmount":"43022100000000000000","expirationTimeSeconds":"4200000000300","makerFee":"0","makerFeeAssetData":"0x0000000000000000000000000000000000000000","takerFee":"0","takerFeeAssetData":"0x0000000000000000000000000000000000000000","salt":"72637260637360712271732819822715422542517835146337292562747232053943662131164","makerAddress":"0x6ecbe1db9ef729cbe972c83fb886247691fb6beb","senderAddress":"0x6ecbe1db9ef729cbe972c83fb886247691fb6beb","feeRecipientAddress":"0x0000000000000000000000000000000000000000","signature":"0x1b154911f2b115a267265d8a34eeb39ad520b1e5d3399803021cda3a6192d9eb9e4562cdff119319021ee1c5976d2a6f93bde3338fa73a3651f509400d209cb6cf03"}`

const TestAddressStringCheckSummed = "0x48BaCB9266a570d521063EF5dD96e61686DbE788"
const TestAddressStringNormalized = "0x48bacb9266a570d521063ef5dd96e61686dbe788"

var TestAddress = common.Address{0x48, 0xba, 0xcb, 0x92, 0x66, 0xa5, 0x70, 0xd5, 0x21, 0x06, 0x3e, 0xf5, 0xdd, 0x96, 0xe6, 0x16, 0x86, 0xdb, 0xe7, 0x88}

func TestNormalizeAddress(t *testing.T) {
	require.Equal(t, TestAddress, common.HexToAddress(TestAddressStringCheckSummed))
	require.Equal(t, TestAddress, common.HexToAddress(TestAddressStringNormalized))

	assert.Equal(t, TestAddressStringCheckSummed, TestAddress.Hex())
	assert.Equal(t, TestAddressStringNormalized, NormalizeAddress(TestAddress))
}

func TestSignedOrder(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var order SignedOrder
	err := json.NewDecoder(strings.NewReader(TestOrderJSON)).Decode(&order)
	require.NoError(err)
	require.Equal("0x48bacb9266a570d521063ef5dd96e61686dbe788", order.ExchangeAddress)

	zrxOrder, err := order.ToZeroExSignedOrder()
	require.NoError(err)

	assert.Equal(order.ExchangeAddress, NormalizeAddress(zrxOrder.ExchangeAddress))
	assert.Equal(order.MakerAssetData, hexutil.Encode(zrxOrder.MakerAssetData))
	assert.Equal(order.TakerAssetAmount, zrxOrder.TakerAssetAmount.String())
	assert.Equal(order.Signature, hexutil.Encode(zrxOrder.Signature))
}

func TestOrder(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var order Order
	err := json.NewDecoder(strings.NewReader(TestOrderJSON)).Decode(&order)
	require.NoError(err)
	require.Equal("0x48bacb9266a570d521063ef5dd96e61686dbe788", order.ExchangeAddress)

	zrxOrder, err := order.ToZeroExOrder()
	require.NoError(err)

	assert.Equal(order.ExchangeAddress, NormalizeAddress(zrxOrder.ExchangeAddress))
	assert.Equal(order.MakerAssetData, hexutil.Encode(zrxOrder.MakerAssetData))
	assert.Equal(order.TakerAssetAmount, zrxOrder.TakerAssetAmount.String())
}
