package eth

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// known max uint256 as a string
// https://www.wolframalpha.com/input/?i=%282+%5E+256%29+-+1
const knownMaxUint256 = "115792089237316195423570985008687907853269984665640564039457584007913129639935"

func TestMaxUint256(t *testing.T) {
	knownMax, ok := new(big.Int).SetString(knownMaxUint256, 10)
	require.True(t, ok)

	assert.Equal(t, 0, knownMax.Cmp(MAX_UINT256))
}

func TestNullAddress(t *testing.T) {
	knownNull := common.HexToAddress("")
	assert.Equal(t, knownNull, NULL_ADDRESS)
}
