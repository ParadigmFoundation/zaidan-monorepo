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

	gasPrice := new(big.Int)
	if req.GasPrice == "" {
		gp, err := hw.provider.Client().SuggestGasPrice(ctx)
		if err != nil {
			return nil, err
		}
		gasPrice = gp
	}

	value := new(big.Int)
	if req.Value == "" {
		value = big.NewInt(0)
	} else {
		_, ok := value.SetString(req.Value, 10)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "unable to parse 'value'")
		}
	}

	nonce, err := hw.provider.Nonce(ctx, hw.makerAddress)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		From:     hw.makerAddress,
		To:       &to,
		Gas:      req.GasLimit,
		GasPrice: gasPrice,
		Value:    value,
		Data:     req.Data,
	}

	gasLimit := msg.Gas
	if gasLimit == 0 {
		gasLimit, err = hw.provider.Client().EstimateGas(ctx, msg)
		if err != nil {
			return nil, err
		}
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
