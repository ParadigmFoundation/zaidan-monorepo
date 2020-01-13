package eth

import "math/big"

// MAX_UINT256 represents the maximum 256-bit unsigned integer (2^256 - 1)
var MAX_UINT256 *big.Int

func init() {
	tmp := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
	MAX_UINT256 = new(big.Int).Sub(tmp, big.NewInt(1))
}
