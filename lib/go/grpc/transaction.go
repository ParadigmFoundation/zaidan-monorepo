package grpc

import (
	fmt "fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/0xProject/0x-mesh/ethereum/wrappers"
)

// ToStruct3 converts an ExecuteZeroExTransactionRequest to a low-level ZeroExTransaction binding
func (etx ExecuteZeroExTransactionRequest) ToStruct3() (wrappers.Struct3, error) {
	salt, ok := new(big.Int).SetString(etx.Salt, 10)
	if !ok {
		return wrappers.Struct3{}, fmt.Errorf("unable to parse 'salt' (%v)", etx.Salt)
	}

	expiration, ok := new(big.Int).SetString(etx.ExpirationTimeSeconds, 10)
	if !ok {
		return wrappers.Struct3{}, fmt.Errorf("unable to parse 'expirationTimeSeconds' (%v)", etx.ExpirationTimeSeconds)
	}

	gasPrice, ok := new(big.Int).SetString(etx.GasPrice, 10)
	if !ok {
		return wrappers.Struct3{}, fmt.Errorf("unable to parse 'gasPrice' (%v)", etx.GasPrice)
	}

	return wrappers.Struct3{
		Salt:                  salt,
		ExpirationTimeSeconds: expiration,
		GasPrice:              gasPrice,
		SignerAddress:         common.HexToAddress(etx.SignerAddress),
		Data:                  etx.Data,
	}, nil
}
