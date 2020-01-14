package core

import (
	"context"
	"math/big"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

// SendTransaction implements grpc.HotWalletServer
func (hw *HotWallet) SendTransaction(ctx context.Context, req *grpc.SendTransactionRequest) (*grpc.SendTransactionResponse, error) {
	to := common.HexToAddress(req.ToAddress)
	fromAcct, err := hw.provider.GetAccount(hw.makerAddress)
	if err != nil {
		return nil, err
	}

	gasPrice, err := hw.parseGasPrice(ctx, req.GasPrice)
	if err != nil {
		return nil, err
	}

	value, err := hw.parseValue(req.Value)
	if err != nil {
		return nil, err
	}

	nonce, err := hw.provider.Nonce(ctx, hw.makerAddress)
	if err != nil {
		return nil, err
	}

	gasLimit, err := hw.getGasLimit(ctx, to, req.GasLimit, value, gasPrice, req.Data)
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, req.Data)
	stx, err := hw.provider.SignTx(ctx, fromAcct, tx)
	if err != nil {
		return nil, err
	}

	if err := hw.provider.Client().SendTransaction(ctx, stx); err != nil {
		return nil, err
	}

	return &grpc.SendTransactionResponse{
		TransactionHash: stx.Hash().Hex(),
		SentAt:          time.Now().Unix(),
	}, nil
}

// parses gas price from request string, if empty string, uses suggested gas price
func (hw *HotWallet) parseGasPrice(ctx context.Context, raw string) (*big.Int, error) {
	if raw == "" {
		gasPrice, err := hw.provider.Client().SuggestGasPrice(ctx)
		if err != nil {
			return nil, err
		}
		return gasPrice, nil
	} else {
		gasPrice, ok := new(big.Int).SetString(raw, 10)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "unable to parse 'gasPrice'")
		}
		return gasPrice, nil
	}
}

// parses ether value from request string (if raw is empty string, returns 0)
func (hw *HotWallet) parseValue(raw string) (*big.Int, error) {
	if raw == "" {
		return big.NewInt(0), nil
	} else {
		value, ok := new(big.Int).SetString(raw, 10)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "unable to parse 'value'")
		}
		return value, nil
	}
}

// estimates gas required for transaction if gas is 0, otherwise returns gas
func (hw *HotWallet) getGasLimit(ctx context.Context, to common.Address, gas uint64, val *big.Int, gasPrice *big.Int, data []byte) (uint64, error) {
	msg := ethereum.CallMsg{
		From:     hw.makerAddress,
		To:       &to,
		Gas:      gas,
		GasPrice: gasPrice,
		Value:    val,
		Data:     data,
	}

	var err error
	gasLimit := msg.Gas
	if gasLimit == 0 {
		gasLimit, err = hw.provider.Client().EstimateGas(ctx, msg)
		if err != nil {
			return 0, err
		}
	}
	return gasLimit, nil
}
