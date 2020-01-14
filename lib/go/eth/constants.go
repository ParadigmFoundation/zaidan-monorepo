package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// MAX_UINT256 represents the maximum 256-bit unsigned integer (2^256 - 1)
var MAX_UINT256 = new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(1))

var NULL_ADDRESS = common.Address{}
