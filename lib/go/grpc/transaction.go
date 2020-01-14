package grpc

import (
	fmt "fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/zrx"
)

func (etx ExecuteZeroExTransactionRequest) ToZeroExTransaction() (*zrx.Transaction, error) {
	salt, ok := new(big.Int).SetString(etx.Salt, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse 'salt' (%v)", etx.Salt)
	}

	expiration, ok := new(big.Int).SetString(etx.ExpirationTimeSeconds, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse 'expirationTimeSeconds' (%v)", etx.ExpirationTimeSeconds)
	}

	gasPrice, ok := new(big.Int).SetString(etx.GasPrice, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse 'gasPrice' (%v)", etx.GasPrice)
	}

	return &zrx.Transaction{
		Salt:                  salt,
		ExpirationTimeSeconds: expiration,
		GasPrice:              gasPrice,
		SignerAddress:         common.HexToAddress(etx.SignerAddress),
		Data:                  etx.Data,
	}, nil
}
