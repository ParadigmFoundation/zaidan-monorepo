package grpc

import (
	fmt "fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

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

func (ztx *ZeroExTransaction) ToZeroExTransaction() (*zrx.Transaction, error) {
	salt, ok := new(big.Int).SetString(ztx.Salt, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse 'salt' (%v)", ztx.Salt)
	}

	gasPrice, ok := new(big.Int).SetString(ztx.GasPrice, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse 'gas_price' (%v)", ztx.GasPrice)
	}

	callData, err := hexutil.Decode(ztx.Data)
	if err != nil {
		return nil, err
	}

	return &zrx.Transaction{
		Salt:                  salt,
		ExpirationTimeSeconds: big.NewInt(ztx.ExpirationTimeSeconds),
		GasPrice:              gasPrice,
		SignerAddress:         common.HexToAddress(ztx.SignerAddress),
		Data:                  callData,
	}, nil
}

func ZeroExTransactionToProto(ztx *zrx.Transaction) *ZeroExTransaction {
	return &ZeroExTransaction{
		Salt:                  ztx.Salt.String(),
		ExpirationTimeSeconds: ztx.ExpirationTimeSeconds.Int64(),
		GasPrice:              ztx.GasPrice.String(),
		SignerAddress:         strings.ToLower(ztx.SignerAddress.Hex()),
		Data:                  hexutil.Encode(ztx.Data),
	}
}
