package core

import (
	"context"
	"math/big"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/zrx"
)

// Transfer implements grpc.HotWalletServer
func (hw *HotWallet) Transfer(ctx context.Context, req *grpc.TransferRequest) (*grpc.TransferResponse, error) {
	token := common.HexToAddress(req.TokenAddress)
	to := common.HexToAddress(req.ToAddress)

	amount, ok := new(big.Int).SetString(req.Amount, 10)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "unable to parse transfer amount")
	}

	if token == zrx.NULL_ADDRESS {
		return hw.transferEther(ctx, to, amount)
	} else {
		return nil, status.Error(codes.Unimplemented, "token transfers not supported yet")
	}
}

func (hw *HotWallet) transferEther(ctx context.Context, to common.Address, amount *big.Int) (*grpc.TransferResponse, error) {
	nonce, err := hw.provider.Nonce(ctx, hw.makerAddress)
	if err != nil {
		return nil, err
	}

	gasPrice, err := hw.provider.Client().SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	account := accounts.Account{Address: hw.makerAddress}
	tx, err := hw.provider.SignTx(ctx, account, types.NewTransaction(nonce, to, amount, 21000, gasPrice, nil))
	if err != nil {
		return nil, err
	}

	if err := hw.provider.Client().SendTransaction(ctx, tx); err != nil {
		return nil, err
	}

	return &grpc.TransferResponse{
		TransactionHash: tx.Hash().Hex(),
	}, nil
}
